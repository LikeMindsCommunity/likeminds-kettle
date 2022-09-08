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
const CloseReportsEndPoint = "/api/close_report"

const ParamPage = "page"
const ParamMemberId = "member_id"
const MemberIds = "member_ids"

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
