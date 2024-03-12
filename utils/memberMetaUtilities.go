package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/cache"
	"github.com/nateshr/likeminds-authentication/logging"
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
	QuestionAnswers []interface{}  `json:"question_answers"`
}

type MemberMetaResponse struct {
	Success      bool         `json:"success"`
	ErrorMessage string       `json:"error_message"`
	Members      []MemberMeta `json:"members"`
}

func fetchmembersMetaFromCache(redisClient *redis.Client, member_ids []string) ([]MemberMeta, []string, error) {

	membersMeta := []MemberMeta{}
	remainingMemberIds := []string{}

	// fetch member meta from cache
	cachKeys := []string{}
	for _, member_id := range member_ids {
		cachKeys = append(cachKeys, fmt.Sprintf("%s_user_meta", member_id)) // TODO: Move to constants
	}

	// Fetch keys from cache
	values, err := cache.MGet(redisClient, cachKeys...)
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
			remainingMemberIds = append(remainingMemberIds, member_ids[i])
		}
	}

	return membersMeta, remainingMemberIds, nil
}

func fetchMembersMetaFromAPI(headers map[string]interface{}, member_ids []string) ([]MemberMeta, error) {

	//Params to be sent in the api/community_member/fetch_access request
	params := map[string]string{
		ParamMemberIds: ParseStringArrayToString(member_ids),
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

	return membersMetaResponse.Members, nil
}

// save response to cache
func saveMembersMetaInCache(redisClient *redis.Client, membersMeta []MemberMeta) {

	for _, memberMeta := range membersMeta {

		parsedData, err := json.Marshal(memberMeta)
		if err != nil {
			logging.Error(fmt.Sprint("error marshalling member data", err))
			continue
		}

		cacheKey := fmt.Sprintf("%s_user_meta", memberMeta.UserUniqueId) //TODO: Move to constants
		if err := cache.Set(redisClient, cacheKey, parsedData, 7*24*time.Hour); err != nil {
			logging.Error(fmt.Sprint("error saving member data to cache", err))
		}
		logging.Info(fmt.Sprintf("Saved member data to cache for user_unique_id: %s", memberMeta.UserUniqueId))
	}
}

// Exposed method to fetch members meta map for user_unique_ids from cache if present else from api
func FetchMemberMetaMapForUserUniqueIds(redisClient *redis.Client, headers map[string]interface{}, user_unique_ids []string,
) (map[string]MemberMeta, error) {

	memberMetaMap := map[string]MemberMeta{}

	if len(user_unique_ids) == 0 {
		return memberMetaMap, nil
	}

	// fetch member meta from cache
	if redisClient != nil {
		membersMeta, remainingMemberIds, err := fetchmembersMetaFromCache(redisClient, user_unique_ids)
		if err != nil {
			return nil, err
		}

		// Add fetched members meta to memberMetaMap
		for _, memberData := range membersMeta {
			memberMetaMap[memberData.UserUniqueId] = memberData
		}

		// update user_unique_ids with remaining user_unique_ids
		user_unique_ids = remainingMemberIds
	}

	// fetch remaining members meta from api
	if len(user_unique_ids) > 0 {

		membersMeta, err := fetchMembersMetaFromAPI(headers, user_unique_ids)
		if err != nil {
			return nil, err
		}

		// Add fetched members meta to memberMetaMap
		for _, memberData := range membersMeta {
			memberMetaMap[memberData.UserUniqueId] = memberData
		}

		// Save fetched members meta to cache in background
		if redisClient != nil {
			go saveMembersMetaInCache(redisClient, membersMeta) // TODO: Test this and move to background
		}
	}

	//Generate user data for remaining user_unique_ids
	for _, memberId := range user_unique_ids {
		if _, ok := memberMetaMap[memberId]; !ok {
			memberMetaMap[memberId] = MemberMeta{
				IsDeleted: true,
			}
		}
	}

	return memberMetaMap, nil
}

// This function is used to fetch members meta from user_ids of feed entity data
func GetUsersMetaFromFeedData(redisClient *redis.Client, headers map[string]interface{}, feedDataArray []interface{},
	dataResponse map[string]interface{},
) (map[string]MemberMeta, error) {

	user_unique_ids := []string{}

	// Fetch user ids from array
	for _, data := range feedDataArray {
		if user_unique_id, ok := data.(map[string]interface{})["uuid"]; ok {
			user_unique_ids = append(user_unique_ids, user_unique_id.(string))
		}
	}

	user_unique_ids = AppendRepostPostUsersFromFeedDataResponse(dataResponse, user_unique_ids)

	// Fetch user data for given user_unique_ids
	user_data, err := FetchMemberMetaMapForUserUniqueIds(redisClient, headers, user_unique_ids)
	if err != nil {
		return nil, err
	}

	return user_data, nil
}
