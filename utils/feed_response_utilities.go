package utils

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

// AppendPollOptionAddedByUsersFromFeedDataResponse | adds uuids from poll options to API response user uuid list
func AppendPollOptionAddedByUsersFromFeedDataResponse(dataResponse map[string]interface{}, userIDs []string) []string {
	// fetch options uuid
	if widgets, ok := dataResponse["widgets"]; ok {
		widgetsData := widgets.(map[string]interface{})
		for _, widgetData := range widgetsData {
			if lmMeta, ok := widgetData.(map[string]interface{})["_lm_meta"]; ok {
				if options, ok := lmMeta.(map[string]interface{})["options"].([]interface{}); ok {
					for _, option := range options {
						if uuid, ok := option.(map[string]interface{})["uuid"]; ok {
							userIDs = append(userIDs, uuid.(string))
						}
					}
				}
			}
		}
	}
	return userIDs
}
