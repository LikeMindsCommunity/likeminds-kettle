package moderation

const FetchCMRights = "/api/fetch_community_manager_rights"
const FetchMemberRights = "/api/fetch_member_rights"
const UpdateCMRights = "/api/update_community_manager_rights"
const UpdateMemberRights = "/api/update_member_rights"

const ParamIsCm = "is_cm"
const ParamUserId = "user_id"
const ParamUUID = "uuid"

const (
	CREATE_POST_RIGHT_ID       = 10
	COMMENT_AND_REPLY_RIGHT_ID = 11
)

const (
	CREATE_POST_PERMISSION_ADDED_ACTION      = "create_post_permit_added"
	CREATE_POST_PERMISSION_REMOVED_ACTION    = "create_post_permit_removed"
	CREATE_COMMENT_PERMISSION_ADDED_ACTION   = "create_comment_permit_added"
	CREATE_COMMENT_PERMISSION_REMOVED_ACTION = "create_comment_permit_removed"
)
