package community

const FetchMembersMetaEndPoint = "/api/community/fetch_members_meta"
const AllMembersV1EndPoint = "/api/v1/all_members"
const EditQuestionsEndPoint = "/api/community/edit_questions"
const FetchQuestionsEndPoint = "/api/community/questions"
const CommunityMemberEndPoint = "/api/community/member"
const MemberStateEndPoint = "/api/members_state"
const RemoveMemberEndPoint = "/api/remove_from_member"
const RemoveCMEndPoint = "/api/remove_community_manager"
const FetchManagementToolsEndPoint = "/api/fetch_management_tools"
const FetchReportsEndPoint = "/api/fetch_reports"
const PushReportEndPoint = "/api/push_report"
const CommunityReportV1EndPoint = "/api/community/report"
const CloseReportsEndPoint = "/api/close_report"
const FetchReportTagsEndPoint = "/api/fetch_report_tags"
const FetchCommunityReportTagsEndPoint = "/api/community/report/tags"
const FetchCommunitySettingsEndPoint = "/api/community/fetch_community_settings"
const EditCommunitySettingsEndPoint = "/api/community/update_community_settings"
const FetchCommunityRightsEndPoint = "/api/fetch_community_setting_rights"
const EditCommunityRightsEndPoint = "/api/update_community_rights"
const FetchCommunityDMSettingsEndPoint = "/api/community/fetch_community_dm_settings"
const EditCommunityDMSettingsEndPoint = "/api/community/update_community_dm_settings"
const UserDMFeedEndpoint = "/api/user/fetch_dm_feed"
const UserCanDMEndpoint = "/api/community_member/can_dm"
const UserSearchEndpoint = "/api/search/member_directory"
const FetchMemberProfileEndPoint = "/api/community_member/fetch_profile"
const EditMemberProfileEndPoint = "/api/community_member/edit_profile"
const FetchMemberChatroomEndPoint = "/api/fetch_user_chatrooms"
const CreateCohortEndPoint = "/api/cohort/create"
const GetCohortEndPoint = "/api/cohort/fetch"
const GetCommunityCohortsEndPoint = "/api/cohort/fetch_community_cohorts"
const DeleteCohortEndPoint = "/api/cohort/delete"
const EditCohortEndPoint = "/api/cohort/update"
const RemoveCohortMemberEndPoint = "/api/cohort/remove_member"
const CommunityFetchFeedEndPoint = "/api/community_member/fetch_feed"
const CommunityFetchPostFeedEndPoint = "/api/community_member/post_feed"
const ConversationNotificationSettingsEndPoint = "/api/community/notification_settings"
const FeedNotificationSettingEndPoint = "/api/community/feed_notification_setting"
const CommunityExcludedChatroomsEndPoint = "/api/community_member/excluded_chatrooms"
const GetMemberChatroomsEndPoint = "/api/community_member/chatroom/status"
const FetchContentDownloadSettingsEndPoint = "/api/community/fetch_content_download_settings"
const EditContentDownloadSettingsEndPoint = "/api/community/update_content_download_settings"
const GetCommunityEndPoint = "/api/community/%s"
const GetCommunityBrandingEndPoint = "/api/community/%s/branding"
const GetCommunityV2Endpoint = "/api/community/fetch"
const UserHomeMetaEndpoint = "/api/community_member/home_meta"
const AcceptRejectCommunityJoinEndpoint = "/api/community_member/join"
const FetchCommunityQuestionFiltersEndpoint = "/api/fetch_filters"
const FetchIntroExamplesEndpoint = "/api/fetch_intro_examples"
const SendCommunityInviteEndpoint = "/api/community/invite"
const CommunityConfigurationsEndpoint = "/api/community/configurations"
const FetchPendingMembersEndpoint = "/api/community_member/pending_members"
const LeaveCommunityEndPoint = "/api/community_member/leave"
const FetchCommunityRemovalReports = "/api/community/removal_reports"
const MemberConnectionEndPoint = "/api/community_member/%s/connection"

const ParamCommunityID = "community_id"
const ParamPage = "page"
const ParamPageSize = "page_size"
const ParamMemberId = "member_id"
const ParamUUID = "uuid"
const ParamUserId = "user_id"
const MemberIds = "member_ids"
const ChatroomIDParam = "chatroom_id"
const FeedroomIDParam = "feedroom_id"
const ParamChatroomIds = "chatroom_ids"
const ParamExcludedChatroomIds = "excluded_chatroom_ids"
const RequestFromParam = "req_from"
const SearchParam = "search"
const SearchTypeParam = "search_type"
const SearchName = "search_name"

const ParamState = "state"
const ParamType = "type"
const ParamEntityType = "entity_type"
const ParamCohortID = "cohort_id"
const ParamPinned = "pinned"
const ParamOrderType = "order_type"
const ParamIsClosed = "is_closed"
const ParamFilterType = "filter_type"
const ParamMemberState = "member_state"
const ParamMemberStates = "member_states"
const ParamConfigurationTypes = "configuration_types"
const ParamQuestionAnswersVersion = "question_answers_version"
const ParamFilterMemberRoles = "filter_member_roles"
const ParamExcludeSelfMember = "exclude_self_user"
const ParamStatus = "status"

const UserChannelReqFrom = "user_channel"
const MemberProfileReqFrom = "member_profile"

const showDmResponse = "show_dm"

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

const (
	ChatFeedType = 0
	PostFeedType = 1
)

const (
	OrderTypeNewest           = 0
	OrderTypeRecentlyActive   = 1
	OrderTypeMostMessages     = 2
	OrderTypeMostParticipants = 3
)

// Notification Types
const (
	NotificationTypeChat = "chat"
	NotificationTypeFeed = "feed"
)
