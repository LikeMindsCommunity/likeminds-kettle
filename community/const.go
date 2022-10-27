package community

const FetchMembersMetaEndPoint = "/api/community/fetch_members_meta"
const AllMembersV1EndPoint = "/api/v1/all_members"
const EditQuestionsEndPoint = "/api/community/edit_questions"
const CommunityMemberEndPoint = "/api/community/member"
const MemberStateEndPoint = "/api/members_state"
const RemoveMemberEndPoint = "/api/remove_from_member"
const RemoveCMEndPoint = "/api/remove_community_manager"
const FetchManagementToolsEndPoint = "/api/fetch_management_tools"
const FetchReportsEndPoint = "/api/fetch_reports"
const PushReportEndPoint = "/api/push_report"
const CloseReportsEndPoint = "/api/close_report"
const FetchReportTagsEndPoint = "/api/fetch_report_tags"
const FetchCommunitySettingsEndPoint = "/api/community/fetch_community_settings"
const EditCommunitySettingsEndPoint = "/api/community/update_community_settings"
const EditCommunityRightsEndPoint = "/api/update_community_rights"
const FetchCommunityDMSettingsEndPoint = "/api/community/fetch_community_dm_settings"
const EditCommunityDMSettingsEndPoint = "/api/community/update_community_dm_settings"
const UserDMFeedEndpoint = "/api/user/fetch_dm_feed"
const UserCanDMEndpoint = "/api/community_member/can_dm"
const UserSearchEndpoint = "/api/search/member_directory"
const FetchMemberProfileEndPoint = "/api/community_member/fetch_profile"
const EditMemberProfileEndPoint = "/api/community_member/edit_profile"
const FetchMemberChatroomEndPoint = "/api/fetch_user_chatrooms"

const ParamPage = "page"
const ParamMemberId = "member_id"
const ParamUserId = "user_id"
const MemberIds = "member_ids"
const ChatroomIDParam = "chatroom_id"
const RequestFromParam = "req_from"
const SearchParam = "search"
const SearchTypeParam = "search_type"
const ParamState = "state"
const ParamType = "type"

type CommunityObject struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Purpose        string `json:"purpose"`
	ImageUrl       string `json:"image_url"`
	ImageUrlRound  string `json:"image_url_round"`
	CreatedBy      string `json:"created_by"`
	PromotersCount int32  `json:"promoters_count"`
	MembersCount   int32  `json:"members_count"`
	MemberState    int32  `json:"member_state"`
}
