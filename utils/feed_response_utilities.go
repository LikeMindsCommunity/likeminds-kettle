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

// AppendPollOptionCreatorsFromFeedDataResponse | adds uuids from poll options to API response user uuid list
func AppendPollOptionCreatorsFromFeedDataResponse(dataResponse map[string]interface{}, userIDs []string) []string {
	// fetch options creator uuid
	if widgets, ok := dataResponse["widgets"]; ok {

		widgetsData := widgets.(map[string]interface{})
		for _, widgetData := range widgetsData {

			// extract _lm_meta from each widget object
			if lmMeta, ok := widgetData.(map[string]interface{})["_lm_meta"]; ok {

				// extract options array from _lm_meta object
				if options, ok := lmMeta.(map[string]interface{})["options"].([]interface{}); ok {
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

	return userIDs
}
