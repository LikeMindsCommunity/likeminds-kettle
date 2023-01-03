package feed

const CreatePostEndPoint = "/post"
const SinglePostEndPoint = "/post/%s"
const SinglePostLikeEndPoint = "/post/%s/like"
const SinglePostCommentEndPoint = "/post/%s/comment"
const SinglePostPinEndPoint = "/post/%s/pin"
const SinglePostSaveEndPoint = "/post/%s/save"
const SingleCommentEndPoint = "/post/%s/comment/%s"
const SingleCommentLikeEndPoint = "/post/%s/comment/%s/like"
const SingleCommentReplyEndPoint = "/post/%s/comment/%s/comment"
const FetchUserSavedPostsEndPoint = "/user/%s/save"
const FetchUserCreatedPostsEndPoint = "/user/%s/post"
const SingleUserActivityEndPoint = "/user/%s/activity"
const FetchUniversalFeedEndPoint = "/feed/universal"

const ParamPage = "page"
const ParamPageSize = "page_size"
const ParamUserIsCm = "user_is_cm"

const (
	CREATE_POST_ACTION     = "create_post"
	VIEW_POST_ACTION       = "view_post"
	DELETE_POST_ACTION     = "delete_post"
	PIN_POST_ACTION        = "pin_post"
	LIKE_POST_ACTION       = "like_post"
	SAVE_POST_ACTION       = "save_post"
	CREATE_COMMENT_ACTION  = "create_comment"
	VIEW_COMMENT_ACTION    = "view_comment"
	DELETE_COMMENT_ACTION  = "delete_comment"
	LIKE_COMMENT_ACTION    = "like_comment"
	CREATE_ACTIVITY_ACTION = "create_activity"
	VIEW_ACTIVITY_ACTION   = "view_activity"
)
