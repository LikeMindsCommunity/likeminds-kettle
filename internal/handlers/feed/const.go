package feed

const CreatePostEndPoint = "/post"
const SinglePostEndPoint = "/post/%s"
const SinglePostLikeEndPoint = "/post/%s/like"
const SinglePostCommentEndPoint = "/post/%s/comment"
const SinglePostPinEndPoint = "/post/%s/pin"
const SinglePostSaveEndPoint = "/post/%s/save"
const SingleCommentEndPoint = "/post/%s/comment/%s"
const SingleCommentByIdEndPoint = "/comment/%s"
const SingleCommentLikeEndPoint = "/post/%s/comment/%s/like"
const SingleCommentReplyEndPoint = "/post/%s/comment/%s/comment"
const FetchUserSavedPostsEndPoint = "/user/%s/save"
const FetchUserCreatedPostsEndPoint = "/user/%s/post"
const SingleUserActivityEndPoint = "/user/%s/activity"
const FetchUserFeedMetaEndPoint = "/user/%s/meta"
const UserCommentsEndPoint = "/user/%s/comment"
const CreatePendingPostEndPoint = "/post/pending"
const PendingPostEndPoint = "/post/pending/%s"
const FetchUserCreatedPendingPostsEndPoint = "/user/%s/post/pending"
const RecomputePersonalisedFeedEndPoint = "/personalised/recompute"
const ReorderPersonalisedFeedEndPoint = "/personalised/reorder"
const FetchPersonalisedFeedEndPoint = "/personalised"

// GetUserActivityEndPoint | swarm API endpoint for user activity list
const GetUserActivityEndPoint = "/user/activity"

// GetUserActivityUnreadCountEndPoint | swarm API endpoint for user activity unread count
const GetUserActivityUnreadCountEndPoint = "/user/%s/activity/unread_count"

// UserActivityMarkReadEndPoint | swarm API endpoint for user activity mark read
const UserActivityMarkReadEndPoint = "/user/%s/activity/%s/mark_read"

const FetchUniversalFeedEndPoint = "/feed/universal"
const FeedExploreEndPoint = "/feed/explore"
const FetchGroupFeedEndPoint = "/feed/group"
const DeleteUserDataEndPoint = "/user"
const FetchPostsEndpoint = "/post"
const FetchCommentsEndpoint = "/comment"
const TopicEndPoint = "/topic"
const SingleTopicEndPoint = "/topic/%s"
const ConnectionFeedEndPoint = "/feed/connection"
const FetchUserTopicsEndPoint = "/user/topics"
const UpdateUserTopicsEndPoint = "/user/%s/topics"

const ParamPage = "page"
const ParamPageSize = "page_size"
const ParamUserIsCm = "user_is_cm"
const ParamFeedroomId = "feedroom_id"
const ParamPostId = "post_id"
const ParamCommentId = "comment_id"
const ParamPostIds = "post_ids"
const ParamPendingPostId = "pending_post_id"
const ParamPendingPostIds = "pending_post_ids"
const ParamCommentIds = "comment_ids"
const ParamTopicIds = "topic_ids"
const ParamWidgetIds = "widget_ids"
const ParamIsEnabled = "is_enabled"
const ParamSearchType = "search_type"
const ParamSearch = "search"
const ParamMinPosts = "min_posts"
const ParamParentIds = "parent_ids"
const ParamTopicsOrderBy = "order_by"
const ParamUserId = "user_id"
const ParamUUIDs = "uuids"

const FeedRepostCommunitySettingType = "feed_repost"

const (
	CREATE_POST_ACTION     = "create_post"
	EDIT_POST_ACTION       = "edit_post"
	VIEW_POST_ACTION       = "view_post"
	DELETE_POST_ACTION     = "delete_post"
	PIN_POST_ACTION        = "pin_post"
	LIKE_POST_ACTION       = "like_post"
	SAVE_POST_ACTION       = "save_post"
	CREATE_COMMENT_ACTION  = "create_comment"
	EDIT_COMMENT_ACTION    = "edit_comment"
	VIEW_COMMENT_ACTION    = "view_comment"
	DELETE_COMMENT_ACTION  = "delete_comment"
	LIKE_COMMENT_ACTION    = "like_comment"
	CREATE_ACTIVITY_ACTION = "create_activity"
	VIEW_USER_ACTIVITY     = "view_user_activity"
	VIEW_ACTIVITY_ACTION   = "view_activity"
	VIEW_REPORT_ENTITY     = "view_report_entity"
	CREATE_TOPIC_ACTION    = "create_topic"
	DELETE_TOPIC_ACTION    = "delete_topic"
	EDIT_TOPIC_ACTION      = "edit_topic"
	IS_MEMBER              = "is_member"
	CHANGE_AUTHOR_ACTION   = "change_author"
	CREATE_FEED_POLL       = "create_feed_poll"
)

const (
	POST_REPORT_TYPE         = 5
	COMMENT_REPORT_TYPE      = 6
	REPLY_REPORT_TYPE        = 7
	PENDING_POST_REPORT_TYPE = 8
)

// constants for attachment_type
const (
	ImageWidget int = iota + 1
	VideoWidget
	DocumentWidget
	LinkWidget
	CustomWidget
	PollWidget
	ArticleWidget
)
