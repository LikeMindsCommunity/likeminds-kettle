package utils

import "github.com/go-redis/redis/v7"

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

			// extract _lm_meta from each widget object
			if lmMeta, ok := widgetData.(map[string]interface{})["_lm_meta"]; ok && lmMeta != nil {

				// extract options array from _lm_meta object
				if lmMetaOptions, ok := lmMeta.(map[string]interface{})["options"]; ok {
					if options, ok := lmMetaOptions.([]interface{}); ok {
						for _, option := range options {

							// extract option creator uuid from each option
							if uuid, ok := option.(map[string]interface{})["uuid"]; ok {
								userIDs = append(userIDs, uuid.(string))
							}
						}
					}
				}
			}
		}
	}

	return userIDs
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

		// Update user topics data in dataResponse
		dataResponse = FetchAndUpdateUserTopicsDataForResponse(redisClient, headers, dataResponse, userUniqueIds)
	}

	return dataResponse, nil
}
