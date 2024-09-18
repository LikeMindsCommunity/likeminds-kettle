package utils

import (
	"fmt"

	"github.com/go-redis/redis/v7"
)

// AppendRepostPostUsersFromFeedDataResponse | adds uuids from reposted_posts container to API response user uuid list
func AppendRepostPostUsersFromFeedDataResponse(dataResponse map[string]interface{}, userIDs []string) []string {
	// fetch respost_post user ids
	if value, ok := dataResponse["reposted_posts"]; ok {
		repostedPost := value.(map[string]interface{})
		for _, postData := range repostedPost {
			if uuid, ok := postData.(map[string]interface{})["uuid"]; ok {
				userIDs = append(userIDs, uuid.(string))
			}
		}
	}

	return userIDs
}

// AppendPollOptionCreatorsFromFeedDataResponse | adds uuids from poll options to API response user uuid list
func AppendPollOptionCreatorsFromFeedDataResponse(dataResponse map[string]interface{}, userIDs []string) []string {

	// fetch options creator uuid
	if widgets, ok := dataResponse["widgets"]; ok {
		widgetsData := widgets.(map[string]interface{})
		for _, widgetData := range widgetsData {
			userIDs = appendOptionCreatorsFromWidget(widgetData, userIDs)
		}
	}

	return userIDs
}

func appendOptionCreatorsFromWidget(widgetData interface{}, userIDs []string) []string {
	if lmMeta, ok := widgetData.(map[string]interface{})["_lm_meta"]; ok && lmMeta != nil {
		userIDs = appendOptionCreatorsFromLMMeta(lmMeta, userIDs)
	}
	return userIDs
}

func appendOptionCreatorsFromLMMeta(lmMeta interface{}, userIDs []string) []string {
	if lmMetaOptions, ok := lmMeta.(map[string]interface{})["options"]; ok {
		if options, ok := lmMetaOptions.([]interface{}); ok {
			for _, option := range options {
				userIDs = appendOptionCreatorFromOption(option, userIDs)
			}
		}
	}
	return userIDs
}

func appendOptionCreatorFromOption(option interface{}, userIDs []string) []string {
	if uuid, ok := option.(map[string]interface{})["uuid"]; ok {
		userIDs = append(userIDs, uuid.(string))
	}
	return userIDs
}

// Add block user name in title of post menu
func AddUserNameInBlockMenuTitle(dataResponse map[string]interface{}) map[string]interface{} {
	if value, ok := dataResponse["posts"]; ok {
		if userData, ok := dataResponse["users"].(map[string]MemberMeta); ok {
			posts := value.([]interface{})
			updatedPostsData := []interface{}{}

			// Fetch menu items from array
			for _, data := range posts {
				var userMemberMetaData MemberMeta
				var dataMap map[string]interface{} = data.(map[string]interface{})

				if userUniqueId, ok := dataMap["uuid"]; ok {
					userMemberMetaData = userData[userUniqueId.(string)]
				}

				userFirstName := GetFirstNameFromName(userMemberMetaData.Name)

				if menuItems, ok := dataMap["menu_items"]; ok && userFirstName != "" {
					updatedMenuItems := []map[string]interface{}{}

					for _, menuItem := range menuItems.([]interface{}) {
						menuItemMap := menuItem.(map[string]interface{})
						menuItemId, ok := menuItemMap["id"].(float64)
						if ok && menuItemId == BlockUserMenuItemID {
							menuItemMap["title"] = fmt.Sprintf(BlockUserMenuItemTitle, userFirstName)
						}

						updatedMenuItems = append(updatedMenuItems, menuItemMap)
					}

					dataMap["menu_items"] = updatedMenuItems
				}

				updatedPostsData = append(updatedPostsData, dataMap)
			}

			dataResponse["posts"] = updatedPostsData
		}
	} else if value, ok := dataResponse["post"]; ok {
		if userData, ok := dataResponse["users"].(map[string]MemberMeta); ok {
			post := value.(interface{})

			// Fetch menu items from array
			var userMemberMetaData MemberMeta
			var dataMap map[string]interface{} = post.(map[string]interface{})

			if userUniqueId, ok := dataMap["uuid"]; ok {
				userMemberMetaData = userData[userUniqueId.(string)]
			}

			userFirstName := GetFirstNameFromName(userMemberMetaData.Name)

			if menuItems, ok := dataMap["menu_items"]; ok && userFirstName != "" {
				updatedMenuItems := []map[string]interface{}{}

				for _, menuItem := range menuItems.([]interface{}) {
					menuItemMap := menuItem.(map[string]interface{})
					menuItemId, ok := menuItemMap["id"].(float64)
					if ok && menuItemId == BlockUserMenuItemID {
						menuItemMap["title"] = fmt.Sprintf(BlockUserMenuItemTitle, userFirstName)
					}

					updatedMenuItems = append(updatedMenuItems, menuItemMap)
				}

				dataMap["menu_items"] = updatedMenuItems
			}

			dataResponse["post"] = post
		}
	}

	return dataResponse
}

func PopulateDataResponseForFeed(headers map[string]interface{}, redisClient *redis.Client, dataResponse map[string]interface{},
) (map[string]interface{}, error) {

	if value, ok := dataResponse["posts"]; ok {

		posts := value.([]interface{})

		if value, ok := dataResponse["filtered_comments"]; ok {
			if commentData, ok := value.(map[string]interface{}); ok {
				for _, val := range commentData {
					posts = append(posts, val)
				}
			}
		}

		userData, userUniqueIds, err := GetUsersMetaFromFeedData(redisClient, headers, posts, dataResponse)
		if err != nil {
			return dataResponse, err
		}

		//Update user data in dataResponse
		dataResponse["users"] = userData

		dataResponse = AddUserNameInBlockMenuTitle(dataResponse)

		// Update user topics data in dataResponse
		dataResponse = FetchAndUpdateUserTopicsDataForResponse(redisClient, headers, dataResponse, userUniqueIds)
	}

	return dataResponse, nil
}
