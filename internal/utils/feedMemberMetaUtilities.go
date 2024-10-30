package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/internal/cache"
	"github.com/nateshr/likeminds-authentication/internal/constants"
	"github.com/nateshr/likeminds-authentication/internal/logging"
)

type SdkClientInfo struct {
	CommunityId  int    `json:"community"`
	UserId       int    `json:"user"`
	UserUniqueId string `json:"user_unique_id"`
	UUID         string `json:"uuid"`
	WidgetId     string `json:"widget_id"`
}

type MemberMeta struct {
	Id              int            `json:"id"`
	Name            string         `json:"name"`
	ImageUrl        string         `json:"image_url"`
	UserUniqueId    string         `json:"user_unique_id"`
	SdkClientInfo   *SdkClientInfo `json:"sdk_client_info"`
	UUID            string         `json:"uuid"`
	IsGuest         bool           `json:"is_guest"`
	IsDeleted       bool           `json:"is_deleted"`
	CustomTitle     string         `json:"custom_title"`
	State           int            `json:"state"`
	QuestionAnswers []interface{}  `json:"question_answers"`
}

type MemberMetaResponse struct {
	Success      bool         `json:"success"`
	ErrorMessage string       `json:"error_message"`
	Members      []MemberMeta `json:"members"`
}

func fetchMembersMetaFromCache(redisClient *redis.Client, communityId int, userUniqueIds []string) ([]MemberMeta, []string, error) {

	membersMeta := []MemberMeta{}
	remainingMemberIds := []string{}

	// fetch member meta from cache
	cachKeys, userIdsMap := []string{}, map[string]bool{}
	for _, userUniqueId := range userUniqueIds {
		if _, ok := userIdsMap[userUniqueId]; !ok {
			cachKeys = append(cachKeys, fmt.Sprintf(cache.UserMetaCacheKey, communityId, userUniqueId))
			userIdsMap[userUniqueId] = true
		}
	}

	// Fetch keys from cache
	values, err := cache.GetFromMultipleKeys(redisClient, cachKeys...)
	if err != nil {
		return nil, nil, err
	}

	// unmarshall values from cache
	for i, value := range values {
		if value != nil {
			var memberMeta MemberMeta
			if err := json.Unmarshal([]byte(value.(string)), &memberMeta); err != nil {
				return nil, nil, err
			}
			membersMeta = append(membersMeta, memberMeta)
		} else {
			userUniqueId := strings.Split(cachKeys[i], "_")[1]
			remainingMemberIds = append(remainingMemberIds, userUniqueId)
		}
	}

	return membersMeta, remainingMemberIds, nil
}

func fetchMembersMetaFromAPI(redisClient *redis.Client, communityId int, headers map[string]interface{}, userUniqueIds []string) ([]MemberMeta, error) {

	//Params to be sent in the api/community_member/fetch_access request
	params := map[string]string{
		ParamMemberIds: ParseStringArrayToString(userUniqueIds),
	}

	//Send Request
	respBytes, _, err := GetRequestResponseWithoutContext(CoreService, FetchMembersMetaEndPoint, GETRequest, headers, params, nil)
	if err != nil {
		return nil, err
	}

	//Parse response
	var membersMetaResponse MemberMetaResponse
	if err := json.Unmarshal(respBytes, &membersMetaResponse); err != nil {
		return nil, err
	}

	if !membersMetaResponse.Success {
		return nil, fmt.Errorf("error fetching members meta: %s", membersMetaResponse.ErrorMessage)
	}

	membersMeta := membersMetaResponse.Members

	// Save fetched members meta to cache in background
	if redisClient != nil {
		go saveMembersMetaInCache(redisClient, communityId, membersMeta)
	}

	return membersMeta, nil
}

// save response to cache
func saveMembersMetaInCache(redisClient *redis.Client, communityId int, membersMeta []MemberMeta) {

	for _, memberMeta := range membersMeta {

		parsedData, err := json.Marshal(memberMeta)
		if err != nil {
			logging.Error(fmt.Sprint("error marshalling member data", err))
			continue
		}

		cacheKey := fmt.Sprintf(cache.UserMetaCacheKey, communityId, memberMeta.UserUniqueId)
		if err := cache.Set(redisClient, cacheKey, parsedData, cache.UserMetaCacheTTL*time.Hour); err != nil {
			logging.Error(fmt.Sprint("error saving member data to cache", err))
		}
		logging.Info(fmt.Sprintf("Saved member data to cache with key: %s", cacheKey))
	}
}

// Utility method to fetch Anonymous user meta for feed
func getAnonymousUserMetaForFeed(redisClient *redis.Client, headers map[string]interface{}, communityId int) MemberMeta {

	// Fetch feed_metadata configurations
	feedMetadata, err := getFeedMetaConfig(redisClient, headers, communityId)
	if err != nil {
		logging.Error(fmt.Sprintf("error fetching feed metadata: %s", err))
	}

	name, imageUrl := constants.AnonymousUserName, ""

	anonUserMetaConfig, ok := feedMetadata.Value[constants.AnonymousUserMetaConfigKey].(map[string]interface{})
	if ok {
		_, ok := anonUserMetaConfig["name"]
		if ok {
			name = anonUserMetaConfig["name"].(string)
		}

		_, ok = anonUserMetaConfig["image_url"]
		if ok {
			imageUrl = anonUserMetaConfig["image_url"].(string)
		}
	}

	anonymousUserMeta := MemberMeta{
		UserUniqueId: constants.AnonymousUserUUID,
		UUID:         constants.AnonymousUserUUID,
		Name:         name,
		ImageUrl:     imageUrl,
		SdkClientInfo: &SdkClientInfo{
			UserUniqueId: constants.AnonymousUserUUID,
			UUID:         constants.AnonymousUserUUID,
		},
	}

	return anonymousUserMeta
}

// Exposed method to fetch members meta map for user_unique_ids from cache if present else from api
func FetchMemberMetaMapForUserUniqueIds(redisClient *redis.Client, headers map[string]interface{}, userUniqueIds []string,
) (map[string]MemberMeta, error) {

	// Fetch communityId from ApiKey
	communityId, err := FetchCommunityIdFromApiKey(redisClient, headers[HeadersApiKey].(string))
	if err != nil {
		return nil, err
	}

	memberMetaMap := map[string]MemberMeta{}

	if len(userUniqueIds) == 0 {
		return memberMetaMap, nil
	}

	// fetch member meta from cache
	if redisClient != nil {
		membersMeta, remainingMemberIds, err := fetchMembersMetaFromCache(redisClient, communityId, userUniqueIds)
		if err != nil {
			return nil, err
		}

		// Add fetched members meta to memberMetaMap
		for _, memberData := range membersMeta {
			memberMetaMap[memberData.UserUniqueId] = memberData
		}

		// update userUniqueIds with remaining userUniqueIds
		userUniqueIds = remainingMemberIds
	}

	// fetch remaining members meta from api
	if len(userUniqueIds) > 0 {

		// Generate anonymous user data if userUniqueIds contains anonymous-user
		for _, userUniqueId := range userUniqueIds {
			if userUniqueId == constants.AnonymousUserUUID {
				memberMetaMap[userUniqueId] = getAnonymousUserMetaForFeed(redisClient, headers, communityId)

				// save in cache
				go saveMembersMetaInCache(redisClient, communityId, []MemberMeta{memberMetaMap[userUniqueId]})
			}
		}

		membersMeta, err := fetchMembersMetaFromAPI(redisClient, communityId, headers, userUniqueIds)
		if err != nil {
			return nil, err
		}

		// Add fetched members meta to memberMetaMap
		for _, memberData := range membersMeta {
			memberMetaMap[memberData.UserUniqueId] = memberData
		}
	}

	// Generate user data for remaining userUniqueIds
	for _, memberId := range userUniqueIds {
		if _, ok := memberMetaMap[memberId]; !ok {
			memberMetaMap[memberId] = MemberMeta{
				IsDeleted:     true,
				SdkClientInfo: &SdkClientInfo{}, // to avoid nil pointer exception
			}
		}
	}

	return memberMetaMap, nil
}

// This function is used to fetch members meta from user_ids of feed entity data
func GetUsersMetaFromFeedData(redisClient *redis.Client, headers map[string]interface{}, feedDataArray []interface{},
	dataResponse map[string]interface{},
) (map[string]MemberMeta, []string, error) {

	userUniqueIds := []string{}

	// Fetch user ids from array
	for _, data := range feedDataArray {
		if user_unique_id, ok := data.(map[string]interface{})["uuid"]; ok {
			userUniqueIds = append(userUniqueIds, user_unique_id.(string))
		}
	}

	userUniqueIds = AppendRepostPostUsersFromFeedDataResponse(dataResponse, userUniqueIds)
	userUniqueIds = AppendPollOptionCreatorsFromFeedDataResponse(dataResponse, userUniqueIds)

	// Fetch user data for given userUniqueIds
	user_data, err := FetchMemberMetaMapForUserUniqueIds(redisClient, headers, userUniqueIds)
	if err != nil {
		return nil, nil, err
	}

	return user_data, userUniqueIds, nil
}
