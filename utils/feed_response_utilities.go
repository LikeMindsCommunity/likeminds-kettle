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
