package community

const FetchMembersMetaEndPoint = "/api/community/fetch_members_meta"
const AllMembersV1EndPoint = "/api/v1/all_members"
const EditQuestionsEndPoint = "/api/community/edit_questions"
const CommunityMemberEndPoint = "/api/community/member"

const ParamPage = "page"
const ParamCommunityId = "community_id"

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
