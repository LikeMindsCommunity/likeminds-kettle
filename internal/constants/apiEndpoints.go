package constants

// OTP Endpoints
const (
	OTPGenerateRoute = "/otp/generate"
	OTPVerifyRoute   = "/otp/verify"
	UserTokenRoute   = "/user/token"
)

// User Endpoints
const (
	UserLoginRoute                = "/user/login"
	UserRefreshRoute              = "/user/refresh"
	UserLogoutRoute               = "/user/logout"
	UserMergeAccountRoute         = "/user/merge_account"
	UserConfigRoute               = "/user/config"
	UserBotRoute                  = "/user/bot"
	UserDevicePushRoute           = "/user/device/push"
	UserSubscriptionWhatsappRoute = "/user/subscription/whatsapp"
	UserMetaRoute                 = "/user/meta"
	UserOTPRoute                  = "/user/otp"
	UserOTPVerifyRoute            = "/user/otp/verify"
	UserSocialLoginRoute          = "/user/social/login"
)

// Home Endpoints
const (
	HomeFetchCommunitiesRoute = "/home/fetch_communities"
	HomeDmMetaRoute           = "/home/dm/meta"
)

// SDK Endpoints
const (
	SDKInitiateRoute   = "/sdk/initiate"
	SDKProjectRoute    = "/sdk/project"
	SDKOnboardingRoute = "/sdk/onboarding"
)

// Chatroom Endpoints
const (
	ChatroomRoute                    = "/chatroom"
	ChatroomTypeRoute                = "/chatroom/type"
	ChatroomScheduleFollowRoute      = "/chatroom/schedule_follow"
	ChatroomPinRoute                 = "/chatroom/pin"
	ChatroomTagRoute                 = "/chatroom/tag"
	CHatroomIdTagRoute               = "/chatroom/:chatroom_id/tag"
	ChatroomParticipantsRoute        = "/chatroom/participants"
	ChatroomSettingsRoute            = "/chatroom/settings"
	ChatroomEnableMemberMessageRoute = "/chatroom/enable_member_message"
	ChatroomAutoFollowMembersRoute   = "/chatroom/auto_follow_members"
	ChatroomFilesRoute               = "/chatroom/files"
	ChatroomMineRoute                = "/chatroom/mine"
	ChatroomSeenRoute                = "/chatroom/seen"
	ChatroomFollowRoute              = "/chatroom/follow"
	ChatroomMuteRoute                = "/chatroom/mute"
	ChatroomRenameRoute              = "/chatroom/rename"
	ChatroomShareRoute               = "/chatroom/share"
	ChatroomPendingRoute             = "/chatroom/pending"
	ChatroomSyncRoute                = "/chatroom/sync"
	ChatroomDMBlockRoute             = "/chatroom/dm/block"
	ChatroomDMRequestRoute           = "/chatroom/dm/request"
	ChatroomDMCreateRoute            = "/chatroom/dm/create"
	ChatroomDMRoute                  = "/chatroom/dm"
	ChatroomDMLimitRoute             = "/chatroom/dm/limit"
	ChatroomSearchRoute              = "/chatroom/search"
	ChatroomCohortRoute              = "/chatroom/cohort"
	ChatroomCohortAccessRoute        = "/chatroom/cohort/access"
	ChatroomHomeRoute                = "/chatroom/home"
	ChatroomMarkReadRoute            = "/chatroom/mark_read"
	ChatroomEventRoute               = "/chatroom/event"
	ChatroomEventMetaRoute           = "/chatroom/event/meta"
	ChatroomEventLinkRoute           = "/chatroom/event/link"
	ChatroomEventUnseenCountRoute    = "/chatroom/event/unseen_count"
	ChatroomEventRecordingsRoute     = "/chatroom/event/recordings"
	ChatroomEventRecordingsMetaRoute = "/chatroom/event/recordings/meta"
	ChatroomEventInstructorsRoute    = "/chatroom/event/instructors"
	ChatroomEventHighlightsRoute     = "/chatroom/event/highlights"
	ChatroomEventTestimonialsRoute   = "/chatroom/event/testimonials"
	ChatroomEventFAQRoute            = "/chatroom/event/faq"
)

// Community Endpoints
const (
	CommunityRoute                                 = "/community"
	CommunityBrandingRoute                         = "/community/branding"
	CommunityQuestionsRoute                        = "/community/questions"
	CommunityQuestionFiltersRoute                  = "/community/question/filters"
	CommunityMemberRoute                           = "/community/member"
	CommunityMemberLeaveRoute                      = "/community/member/leave"
	CommunityMemberStateRoute                      = "/community/member/state"
	CommunityMemberRoleRoute                       = "/community/member/role"
	CommunityManagerRemoveRoute                    = "/community/manager/remove"
	CommunityAdminRoute                            = "/community/admin"
	CommunityMemberRemoveRoute                     = "/community/member/remove"
	CommunityManagementToolRoute                   = "/community/management/tool"
	CommunityReportRoute                           = "/community/report"
	CommunityReportTagRoute                        = "/community/report/tag"
	CommunitySettingsRoute                         = "/community/settings"
	CommunityRightsRoute                           = "/community/rights"
	CommunityDMSettingsRoute                       = "/community/settings/dm"
	CommunityFeedDMRoute                           = "/community/feed/dm"
	CommunityDMStatusRoute                         = "/community/dm/status"
	CommunityMemberSearchRoute                     = "/community/member/search"
	CommunityMemberProfileRoute                    = "/community/member/profile"
	CommunityMemberChatroomRoute                   = "/community/member/chatroom"
	CommunityMemberUserIdChannelsRoute             = "/community/member/:user_id/channel"
	CommunityMemberChannelStatusRoute              = "/community/member/channel/status"
	CommunityCohortRoute                           = "/community/cohort"
	CommunityCohortIdRoute                         = "/community/cohort/:cohort_id"
	CommunityCohortMemberRoute                     = "/community/cohort/member"
	CommunityFeedRoute                             = "/community/feed"
	CommunitySettingsNotificationConversationRoute = "/community/settings/notification/conversation"
	CommunitySettingsNotificationFeedRoute         = "/community/settings/notification/feed"
	CommunitySettingsNotificationRoute             = "/community/settings/notification"
	CommunityTagRoute                              = "/community/tag"
	CommunityContentDownloadSettingsRoute          = "/community/settings/content_download"
	CommunityMemberHomeMetaRoute                   = "/community/member/home/meta"
	CommunityMemberJoin                            = "/community/member/join"
	CommunityIntroExamplesRoute                    = "/community/intro_examples"
	CommunityInviteRoute                           = "/community/invite"
	CommunityConfigurationsRoute                   = "/community/configurations"
	CommunityMemberPendingRoute                    = "/community/member/pending"
	CommunityRemovalReportsRoute                   = "/community/removal_reports"
	CommunityMemberConnectionRoute                 = "/community/member/:user_id/connection"
	CommunityMemberConnectionMetaRoute             = "/community/member/:user_id/connection_meta"
)

// Moderation Endpoints
const (
	ModerationRightsRoute = "/moderation/rights"
)

// Conversation Endpoints
const (
	ConversationRoute                   = "/conversation"
	ConversationReactionRoute           = "/conversation/reaction"
	ConversationPollRoute               = "/conversation/poll"
	ConversationPollSubmitRoute         = "/conversation/poll/submit"
	ConversationPollUsersRoute          = "/conversation/poll/users"
	ConversationTopicRoute              = "/conversation/topic"
	ConversationEventAttendRoute        = "/conversation/event/attend"
	ConversationEventAttendedRoute      = "/conversation/event/attended"
	ConversationEventUnseenRoute        = "/conversation/event/unseen"
	ConversationEventLastSeenRoute      = "/conversation/event/last_seen"
	ConversationEventLinkRoute          = "/conversation/event/link"
	ConversationEventRoute              = "/conversation/event"
	ConversationPreviewUnreadRoute      = "/conversation/preview/unread"
	ConversationPreviewUnreadCountRoute = "/conversation/preview/unread_count"
	ConversationNotificationUnreadRoute = "/conversation/notification/unread"
	ConversationSyncRoute               = "/conversation/sync"
	ConversationSearchRoute             = "/conversation/search"
)

// FEED Endpoints
const (
	FeedPostRoute                   = "/feed/post"
	FeedPostIDRoute                 = "/feed/post/:post_id"
	FeedPostIDLikeRoute             = "/feed/post/:post_id/like"
	FeedPostIDPinRoute              = "/feed/post/:post_id/pin"
	FeedPostIDSaveRoute             = "/feed/post/:post_id/save"
	FeedPostIDCommentRoute          = "/feed/post/:post_id/comment"
	FeedPostIDCommentIDRoute        = "/feed/post/:post_id/comment/:comment_id"
	FeedPostIDCommentIDLikeRoute    = "/feed/post/:post_id/comment/:comment_id/like"
	FeedPostIDCommentIDCommentRoute = "/feed/post/:post_id/comment/:comment_id/comment"
	FeedUserIDSaveRoute             = "/feed/user/:user_id/save"
	FeedUserIDPostRoute             = "/feed/user/:user_id/post"
	FeedUserIDCommentRoute          = "/feed/user/:user_id/comment"
	FeedUserActivityRoute           = "/feed/user/activity"
	FeedUserIdActivityRoute         = "/feed/user/:user_id/activity"
	FeedUserActivityUnreadCount     = "/feed/user/activity/unread_count"
	FeedUserActivityIDMarkReadRoute = "/feed/user/activity/:activity_id/mark_read"
	FeedUserIDMetaRoute             = "/feed/user/:user_id/meta"
	FeedUserIDPostPendingRoute      = "/feed/user/:user_id/post/pending"
	FeedUniversalRoute              = "/feed/universal"
	FeedGroupRoute                  = "/feed/group"
	FeedTopicRoute                  = "/feed/topic"
	FeedTopicIDRoute                = "/feed/topic/:topic_id"
	FeedConnectionRoute             = "/feed/connection"
	FeedPostPendingRoute            = "feed/post/pending"
	FeedPostPendingPostIDRoute      = "feed/post/pending/:pending_post_id"
	FeedUserTopicsRoute             = "/feed/user/topics"
	FeedUserUUIDTopicsRoute         = "/feed/user/:uuid/topics"
	FeedPersonalisedRoute           = "/feed/personalised"
	FeedPersonalisedRecomputeRoute  = "/feed/personalised/recompute"
	FeedPersonalisedReorderRoute    = "/feed/personalised/reorder"
	FeedPostSeenRoute               = "/feed/post/seen"
)

// Utility Endpoints
const (
	HelperUrlRoute         = "/helper/url"
	HelperMediaUploadRoute = "/helper/media/upload"
	HelperS3UploadRoute    = "/helper/s3/upload"
)

// Feedroom Endpoints
const (
	FeedroomRoute                 = "/feedroom"
	FeedroomActionRoute           = "/feedroom/action"
	FeedroomSettingsRoute         = "/feedroom/settings"
	FeedroomTypeRoute             = "/feedroom/type"
	FeedroomEnableMemberPostRoute = "/feedroom/enable_member_post"
	FeedroomPinRoute              = "/feedroom/pin"
	FeedroomAutoJoinMembersRoute  = "/feedroom/auto_join_members"
	FeedroomParticipantsRoute     = "/feedroom/participants"
	FeedroomCohortAccessRoute     = "/feedroom/cohort/access"
	FeedroomMineRoute             = "/feedroom/mine"
	FeedroomFollowRoute           = "/feedroom/follow"
	FeedroomTagRoute              = "/feedroom/tag"
	FeedroomIDTagRoute            = "/feedroom/:feedroom_id/tag"
)

// Channel Endpoints
const (
	ChannelRoute                                       = "/channel"
	ChannelInvitesRoute                                = "/channel/invites"
	ChannelInviteRoute                                 = "/channel/invite"
	ChannelChannelIdSettingsMemberParticipantUUIDRoute = "/channel/:channel_id/settings/member/:participant_uuid"
)

// Search Endpoints
const (
	SearchChannelRoute        = "/search/channel"
	SearchMessageRoute        = "/search/message"
	SearchPostRoute           = "/search/post"
	SearchPostUserUserIdRoute = "/search/post/user/:user_id"
	SearchRoute               = "/search"
)

// Widget Endpoints
const (
	WidgetRoute   = "/widget"
	WidgetIdRoute = "/widget/:widget_id"
)

// Feed Poll Endpoints
const (
	PollIdRoute     = "/poll/:poll_id"
	PollIdVoteRoute = "/poll/:poll_id/vote"
)

// Webhook Endpoints
const (
	WebhookRoute   = "/webhook"
	WebhookIdRoute = "/webhook/:webhook_id"
)

// Logging Endpoints
const (
	LogsEndpoint = "/logs"
)

// Internal Endpoints
const (
	CacheRoute = "/cache"
)
