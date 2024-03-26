package main

import (
	"fmt"

	"github.com/nateshr/likeminds-authentication/internalServices"
	"github.com/nateshr/likeminds-authentication/logging"
	"github.com/nateshr/likeminds-authentication/middleware"
	"github.com/nateshr/likeminds-authentication/poll"
	"github.com/nateshr/likeminds-authentication/utility/logger"
	"github.com/nateshr/likeminds-authentication/utility/monitoring"
	"github.com/nateshr/likeminds-authentication/webhook"

	"github.com/nateshr/likeminds-authentication/widget"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/cache"
	"github.com/nateshr/likeminds-authentication/channel"
	"github.com/nateshr/likeminds-authentication/chatroom"
	"github.com/nateshr/likeminds-authentication/community"
	"github.com/nateshr/likeminds-authentication/conversation"
	"github.com/nateshr/likeminds-authentication/feed"
	"github.com/nateshr/likeminds-authentication/feedroom"
	"github.com/nateshr/likeminds-authentication/home"
	"github.com/nateshr/likeminds-authentication/moderation"
	"github.com/nateshr/likeminds-authentication/otp"
	"github.com/nateshr/likeminds-authentication/sdk"
	"github.com/nateshr/likeminds-authentication/search"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utility"
	"github.com/nateshr/likeminds-authentication/web"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	redisClient *redis.Client
	router      *gin.Engine
)

func main() {
	var AppVersion string = "2.24.1"

	initGin()
	redisClient = cache.InitRedis()
	router.Use(cors.New(enableCors()))
	router.Use(middleware.ApiMiddleware(redisClient))
	router.Use(middleware.LoggingMiddleware())
	//Attach prometheus service as middleware
	prometheusService := getPrometheusMetricService()
	if prometheusService != nil {
		router.Use(monitoring.PrometheusMiddleware(prometheusService))
	}

	router.GET("", web.Home)

	// OTP Apis
	router.GET("/otp/generate", otp.GenerateOTP)
	router.GET("/otp/verify", otp.VerifyOTP)
	router.GET("/user/token", user.CreateToken)

	// User Apis
	router.POST("/user/login", middleware.OTMValidationMiddleware(), user.Login)
	router.POST("/user/refresh", middleware.RTMValidationMiddleware(redisClient), user.Refresh)
	router.POST("/user/logout", middleware.LogoutValidationMiddleware(redisClient), user.Logout)
	router.POST("/user/merge_account", middleware.LTMValidationMiddleware(redisClient, true), user.MergeAccount)
	router.GET("/user/config", middleware.LTMValidationMiddleware(redisClient, true), user.Config)
	router.GET("/user/bot", middleware.LTMValidationMiddleware(redisClient, true), user.GetBot)
	router.POST("/user/device/push", middleware.LTMValidationMiddleware(redisClient, true), user.PushUserToken)
	router.POST("/user/subscription/whatsapp", user.WASubscription)
	router.GET("/user/meta", middleware.LTMValidationMiddleware(redisClient, true), user.UserMeta)
	router.POST("/user/otp", middleware.OTMValidationMiddleware(), user.GenerateUserOTP)
	router.GET("/user/otp/verify", middleware.OTMValidationMiddleware(), user.VerifyUserOTP)
	router.GET("/user/social/login", middleware.OTMValidationMiddleware(), user.UserSocialLogin)

	// Home Apis
	router.POST("/home/fetch_communities", middleware.LTMValidationMiddleware(redisClient, true), home.FetchCommunities)
	router.GET("/home/dm/meta", middleware.LTMValidationMiddleware(redisClient, true), home.DMHome)

	// SDK Apis
	router.POST("/sdk/initiate", middleware.VTMValidationMiddleware(false), sdk.InitiateSDK)
	router.GET("/sdk/initiate", middleware.LTMValidationMiddleware(redisClient, true), sdk.FetchSdkUserInfo)
	router.POST("/sdk/project", middleware.LTMValidationMiddleware(redisClient, true), sdk.CreateProject)
	router.GET("/sdk/project", middleware.LTMValidationMiddleware(redisClient, true), sdk.GetProject)
	router.PUT("/sdk/project", middleware.LTMValidationMiddleware(redisClient, true), sdk.EditProject)
	router.DELETE("/sdk/project", middleware.LTMValidationMiddleware(redisClient, true), sdk.DeleteProject)
	router.GET("/sdk/onboarding", middleware.OTMValidationMiddleware(), sdk.GetScreen)
	router.POST("/sdk/onboarding", middleware.LTMValidationMiddleware(redisClient, true), sdk.CreateScreen)
	router.PUT("/sdk/onboarding", middleware.LTMValidationMiddleware(redisClient, true), sdk.EditScreen)
	router.DELETE("/sdk/onboarding", middleware.LTMValidationMiddleware(redisClient, true), sdk.DeleteScreen)

	// Chatroom Apis
	router.GET("/chatroom", middleware.LTMValidationMiddleware(redisClient, true), chatroom.GetChatroom)
	router.POST("/chatroom", middleware.LTMValidationMiddleware(redisClient, true), chatroom.CreateChatroom)
	router.PUT("/chatroom", middleware.LTMValidationMiddleware(redisClient, true), chatroom.EditChatroom)
	router.DELETE("/chatroom", middleware.LTMValidationMiddleware(redisClient, true), chatroom.DeleteChatroom)
	router.GET("/chatroom/type", middleware.LTMValidationMiddleware(redisClient, true), chatroom.GetChatroomTypeStatus)
	router.PUT("/chatroom/type", middleware.LTMValidationMiddleware(redisClient, true), chatroom.ChangeChatroomType)
	router.POST("/chatroom/schedule_follow", middleware.LTMValidationMiddleware(redisClient, true), chatroom.ScheduleFollow)
	router.PUT("/chatroom/pin", middleware.LTMValidationMiddleware(redisClient, true), chatroom.PinChatroom)
	router.GET("/chatroom/tag", middleware.LTMValidationMiddleware(redisClient, true), chatroom.GetTaggingList)
	router.GET("/chatroom/:chatroom_id/tag", middleware.LTMValidationMiddleware(redisClient, true), chatroom.GetTaggingList)
	router.GET("/chatroom/participants", middleware.LTMValidationMiddleware(redisClient, true), chatroom.GetParticipants)
	router.POST("/chatroom/participants", middleware.LTMValidationMiddleware(redisClient, true), chatroom.AddParticipants)
	router.DELETE("/chatroom/participants", middleware.LTMValidationMiddleware(redisClient, true), chatroom.RemoveParticipants)
	router.GET("/chatroom/settings", middleware.LTMValidationMiddleware(redisClient, true), chatroom.GetChatroomSettings)
	router.PUT("/chatroom/settings", middleware.LTMValidationMiddleware(redisClient, true), chatroom.EditChatroomSettings)
	router.PUT("/chatroom/enable_member_message", middleware.LTMValidationMiddleware(redisClient, true), chatroom.EnableMemberMessage)
	router.PUT("/chatroom/auto_follow_members", middleware.LTMValidationMiddleware(redisClient, true), chatroom.AutoFollowMembers)
	router.PUT("/chatroom/files", middleware.LTMValidationMiddleware(redisClient, true), chatroom.UpdateFiles)
	router.GET("/chatroom/mine", middleware.LTMValidationMiddleware(redisClient, true), chatroom.MyChatrooms)
	router.PUT("/chatroom/seen", middleware.LTMValidationMiddleware(redisClient, true), chatroom.CollabcardSeen)
	router.PUT("/chatroom/follow", middleware.LTMValidationMiddleware(redisClient, true), chatroom.ChatroomFollow)
	router.PUT("/chatroom/mute", middleware.LTMValidationMiddleware(redisClient, true), chatroom.MuteChatroom)
	router.PUT("/chatroom/rename", middleware.LTMValidationMiddleware(redisClient, true), chatroom.RenameChatroom)
	router.GET("/chatroom/share", middleware.LTMValidationMiddleware(redisClient, true), chatroom.FetchShareUrl)
	router.GET("/chatroom/pending", middleware.LTMValidationMiddleware(redisClient, true), chatroom.FetchPendingChatroom)
	router.PUT("/chatroom/pending", middleware.LTMValidationMiddleware(redisClient, true), chatroom.ActionPendingChatroom)
	router.GET("/chatroom/sync", middleware.LTMValidationMiddleware(redisClient, true), chatroom.SyncChatrooms)
	router.POST("/chatroom/dm/block", middleware.LTMValidationMiddleware(redisClient, true), chatroom.ChatroomBlock)
	router.POST("/chatroom/dm/request", middleware.LTMValidationMiddleware(redisClient, true), chatroom.InitiatingDMRequest)
	router.POST("/chatroom/dm/create", middleware.LTMValidationMiddleware(redisClient, true), chatroom.CreateDM)
	router.GET("/chatroom/dm", middleware.LTMValidationMiddleware(redisClient, true), chatroom.ListDMChatrooms)
	router.GET("/chatroom/dm/limit", middleware.LTMValidationMiddleware(redisClient, true), chatroom.DMLimit)
	router.GET("/chatroom/search", middleware.LTMValidationMiddleware(redisClient, true), chatroom.ChatroomSearch)
	router.POST("/chatroom/cohort", middleware.LTMValidationMiddleware(redisClient, true), chatroom.AddCohortToChatroom)
	router.DELETE("/chatroom/cohort", middleware.LTMValidationMiddleware(redisClient, true), chatroom.RemoveCohortFromChatroom)
	router.GET("/chatroom/cohort/access", middleware.LTMValidationMiddleware(redisClient, true), chatroom.GetCohortAccess)
	router.PUT("/chatroom/cohort/access", middleware.LTMValidationMiddleware(redisClient, true), chatroom.EditCohortAccess)
	router.GET("/chatroom/home", middleware.LTMValidationMiddleware(redisClient, true), chatroom.GetChatroomHome)
	router.POST("/chatroom/mark_read", middleware.LTMValidationMiddleware(redisClient, true), chatroom.ChatroomMarkRead)
	router.GET("/chatroom/event", middleware.LTMValidationMiddleware(redisClient, true), chatroom.FetchEvents)
	router.POST("/chatroom/event", middleware.LTMValidationMiddleware(redisClient, true), chatroom.CreateEvent)
	router.PUT("/chatroom/event", middleware.LTMValidationMiddleware(redisClient, true), chatroom.EditEvent)
	router.GET("/chatroom/event/meta", middleware.LTMValidationMiddleware(redisClient, true), chatroom.FetchEventMeta)
	router.GET("/chatroom/event/link", middleware.LTMValidationMiddleware(redisClient, true), chatroom.FetchEventLinks)
	router.GET("/chatroom/event/unseen_count", middleware.LTMValidationMiddleware(redisClient, true), chatroom.FetchEventUnseenCount)
	router.POST("/chatroom/event/recordings", middleware.LTMValidationMiddleware(redisClient, true), chatroom.UploadEventRecordings)
	router.DELETE("/chatroom/event/recordings", middleware.LTMValidationMiddleware(redisClient, true), chatroom.DeleteEventRecordings)
	router.POST("/chatroom/event/recordings/meta", middleware.LTMValidationMiddleware(redisClient, true), chatroom.UploadEventRecordingsMeta)
	router.DELETE("/chatroom/event/recordings/meta", middleware.LTMValidationMiddleware(redisClient, true), chatroom.DeleteEventRecordingsMeta)
	router.POST("/chatroom/event/instructors", middleware.LTMValidationMiddleware(redisClient, true), chatroom.AddEventInstructors)
	router.POST("/chatroom/event/highlights", middleware.LTMValidationMiddleware(redisClient, true), chatroom.AddEventHighlights)
	router.POST("/chatroom/event/testimonials", middleware.LTMValidationMiddleware(redisClient, true), chatroom.AddEventTestimonials)
	router.POST("/chatroom/event/faq", middleware.LTMValidationMiddleware(redisClient, true), chatroom.AddEventFAQ)

	// Community Apis
	router.GET("/community", middleware.LTMValidationMiddleware(redisClient, true), community.Community)
	router.GET("/community/branding", middleware.LTMValidationMiddleware(redisClient, true), community.Branding)
	router.POST("/community/questions", middleware.LTMValidationMiddleware(redisClient, true), community.EditQuestions)
	router.GET("/community/questions", middleware.LTMorVTMValidationMiddleware(), community.GetQuestions)
	router.GET("/community/question/filters", middleware.LTMValidationMiddleware(redisClient, true), community.GetCommunityQuestionFilters)
	router.GET("/community/member", middleware.LTMValidationMiddleware(redisClient, true), community.GetMember)
	router.POST("/community/member", middleware.LTMValidationMiddleware(redisClient, true), community.AddMember)
	router.DELETE("/community/member", middleware.LTMValidationMiddleware(redisClient, true), community.RemoveMembers)
	router.DELETE("/community/member/leave", middleware.LTMValidationMiddleware(redisClient, true), community.LeaveCommunity)
	router.PUT("/community/member", middleware.LTMValidationMiddleware(redisClient, true), community.EditMember)
	router.GET("/community/member/state", middleware.LTMValidationMiddleware(redisClient, true), community.FetchMemberState)
	router.GET("/community/member/role", middleware.LTMValidationMiddleware(redisClient, true), community.FetchMemberRole)
	router.DELETE("/community/manager/remove", middleware.LTMValidationMiddleware(redisClient, true), community.RemoveCommunityManager)
	router.DELETE("/community/admin", middleware.LTMValidationMiddleware(redisClient, true), community.RemoveCommunityManager)
	router.DELETE("/community/member/remove", middleware.LTMValidationMiddleware(redisClient, true), community.RemoveMember)
	router.GET("/community/management/tool", middleware.LTMValidationMiddleware(redisClient, true), community.GetManagementTools)
	router.GET("/community/report", middleware.LTMValidationMiddleware(redisClient, true), community.GetReport)
	router.POST("/community/report", middleware.LTMValidationMiddleware(redisClient, true), community.PushReport)
	router.DELETE("/community/report", middleware.LTMValidationMiddleware(redisClient, true), community.CloseReport)
	router.PATCH("/community/report", middleware.LTMValidationMiddleware(redisClient, true), community.UpdateReports)
	router.GET("/community/report/tag", middleware.LTMValidationMiddleware(redisClient, true), community.GetReportTags)
	router.GET("/community/settings", middleware.LTMValidationMiddleware(redisClient, true), community.GetCommunitySettings)
	router.PUT("/community/settings", middleware.LTMValidationMiddleware(redisClient, true), community.UpdateCommunitySettings)
	router.GET("/community/rights", middleware.LTMValidationMiddleware(redisClient, true), community.GetCommunityRights)
	router.PUT("/community/rights", middleware.LTMValidationMiddleware(redisClient, true), community.EditCommunityRights)
	router.PATCH("/community/rights", middleware.LTMValidationMiddleware(redisClient, true), community.UpdateCommunityRights)
	router.GET("/community/settings/dm", middleware.LTMValidationMiddleware(redisClient, true), community.GetCommunityDMSettings)
	router.PUT("/community/settings/dm", middleware.LTMValidationMiddleware(redisClient, true), community.EditCommunityDMSettings)
	router.GET("/community/feed/dm", middleware.LTMValidationMiddleware(redisClient, true), community.DMFeed)
	router.GET("/community/dm/status", middleware.LTMValidationMiddleware(redisClient, true), community.DMStatus)
	router.GET("/community/member/search", middleware.LTMValidationMiddleware(redisClient, true), community.MemberSearch)
	router.GET("/community/member/profile", middleware.LTMValidationMiddleware(redisClient, true), community.GetMemberProfile)
	router.PUT("/community/member/profile", middleware.LTMValidationMiddleware(redisClient, true), community.EditMemberProfile)
	router.GET("/community/member/chatroom", middleware.LTMValidationMiddleware(redisClient, true), community.MemberChatroom)
	router.GET("/community/member/:user_id/channel", middleware.LTMValidationMiddleware(redisClient, true), community.CommunityMemberChannels)
	router.GET("/community/member/channel/status", middleware.LTMValidationMiddleware(redisClient, true), community.GetMemberChannels)
	router.POST("/community/cohort", middleware.LTMValidationMiddleware(redisClient, true), community.CreateCohort)
	router.GET("/community/cohort", middleware.LTMValidationMiddleware(redisClient, true), community.GetCohort)
	router.GET("/community/cohort/:cohort_id", middleware.LTMValidationMiddleware(redisClient, true), community.FetchCohort)
	router.DELETE("/community/cohort", middleware.LTMValidationMiddleware(redisClient, true), community.DeleteCohort)
	router.PUT("/community/cohort", middleware.LTMValidationMiddleware(redisClient, true), community.EditCohort)
	router.DELETE("/community/cohort/member", middleware.LTMValidationMiddleware(redisClient, true), community.RemoveCohortMember)
	router.GET("/community/feed", middleware.LTMValidationMiddleware(redisClient, true), community.GetCommunityFeed)
	router.GET("/community/settings/notification/conversation", middleware.LTMValidationMiddleware(redisClient, true), community.GetConversationNotificationSettings)
	router.PUT("/community/settings/notification/conversation", middleware.LTMValidationMiddleware(redisClient, true), community.EditConversationNotificationSettings)
	router.GET("/community/settings/notification/feed", middleware.LTMValidationMiddleware(redisClient, true), community.GetFeedNotificationSettings)
	router.PUT("/community/settings/notification/feed", middleware.LTMValidationMiddleware(redisClient, true), community.EditFeedNotificationSettings)
	router.GET("/community/settings/notification", middleware.LTMValidationMiddleware(redisClient, true), community.GetNotificationSettings)
	router.PUT("/community/settings/notification", middleware.LTMValidationMiddleware(redisClient, true), community.EditNotificationSettings)
	router.GET("/community/tag", middleware.LTMValidationMiddleware(redisClient, true), community.GetTaggingList)
	router.GET("/community/settings/content_download", middleware.LTMValidationMiddleware(redisClient, true), community.GetContentDownloadSettings)
	router.PUT("/community/settings/content_download", middleware.LTMValidationMiddleware(redisClient, true), community.EditContentDownloadSettings)
	router.GET("/community/member/home/meta", middleware.LTMValidationMiddleware(redisClient, true), community.MemberHomeMeta)
	router.PUT("/community/member/join", middleware.LTMValidationMiddleware(redisClient, true), community.AcceptRejectJoinCommunity)
	router.GET("/community/intro_examples", middleware.LTMorVTMValidationMiddleware(), community.GetIntroExamples)
	router.POST("/community/invite", middleware.LTMValidationMiddleware(redisClient, true), community.SendCommunityInvite)
	router.GET("/community/configurations", middleware.LTMValidationMiddleware(redisClient, true), community.GetCommunityConfigurations)
	router.PATCH("/community/configurations", middleware.LTMValidationMiddleware(redisClient, true), community.UpdateCommunityConfigurations)
	router.GET("/community/member/pending", middleware.LTMValidationMiddleware(redisClient, true), community.GetPendingCommunityMembers)
	router.GET("/community/removal_reports", middleware.LTMValidationMiddleware(redisClient, true), community.GetRemovalReports)
	router.POST("/community/member/:user_id/connection", middleware.LTMValidationMiddleware(redisClient, true), community.CreateMemberConnection)
	router.PATCH("/community/member/:user_id/connection", middleware.LTMValidationMiddleware(redisClient, true), community.AcceptRejectMemberConnection)
	router.GET("/community/member/:user_id/connection", middleware.LTMValidationMiddleware(redisClient, true), community.GetMemberConnection)

	// Moderation Apis
	router.GET("/moderation/rights", middleware.LTMValidationMiddleware(redisClient, true), moderation.GetRights)
	router.PUT("/moderation/rights", middleware.LTMValidationMiddleware(redisClient, true), moderation.EditRights)
	router.PATCH("/moderation/rights", middleware.LTMValidationMiddleware(redisClient, true), moderation.UpdateRights)

	// Conversation Apis
	router.GET("/conversation", middleware.LTMValidationMiddleware(redisClient, true), conversation.GetConversation)
	router.POST("/conversation", middleware.LTMValidationMiddleware(redisClient, true), conversation.CreateConversation)
	router.PUT("/conversation", middleware.LTMValidationMiddleware(redisClient, true), conversation.EditConversation)
	router.DELETE("/conversation", middleware.LTMValidationMiddleware(redisClient, true), conversation.DeleteConversation)
	router.PUT("/conversation/reaction", middleware.LTMValidationMiddleware(redisClient, true), conversation.AddReaction)
	router.DELETE("/conversation/reaction", middleware.LTMValidationMiddleware(redisClient, true), conversation.RemoveReaction)
	router.POST("/conversation/poll", middleware.LTMValidationMiddleware(redisClient, true), conversation.AddPoll)
	router.POST("/conversation/poll/submit", middleware.LTMValidationMiddleware(redisClient, true), conversation.SubmitPoll)
	router.GET("/conversation/poll/users", middleware.LTMValidationMiddleware(redisClient, true), conversation.PollUsers)
	router.PUT("/conversation/topic", middleware.LTMValidationMiddleware(redisClient, true), conversation.SetTopic)
	router.PUT("/conversation/event/attend", middleware.LTMValidationMiddleware(redisClient, true), conversation.EventAttend)
	router.PUT("/conversation/event/attended", middleware.LTMValidationMiddleware(redisClient, true), conversation.EventAttended)
	router.GET("/conversation/event/unseen", middleware.LTMValidationMiddleware(redisClient, true), conversation.FetchEventUnseenCount)
	router.PUT("/conversation/event/last_seen", middleware.LTMValidationMiddleware(redisClient, true), conversation.UpdateLastSeenEvent)
	router.GET("/conversation/event/link", middleware.LTMValidationMiddleware(redisClient, true), conversation.FetchEventLink)
	router.GET("/conversation/event", middleware.LTMValidationMiddleware(redisClient, true), conversation.FetchAllEvents)
	router.GET("/conversation/preview/unread", middleware.LTMValidationMiddleware(redisClient, true), conversation.FetchUnreadPreviews)
	router.GET("/conversation/preview/unread_count", middleware.LTMValidationMiddleware(redisClient, true), conversation.FetchPreviewUnreadMessagesCount)
	router.GET("/conversation/notification/unread", middleware.LTMValidationMiddleware(redisClient, true), conversation.UnreadConversationNotification)
	router.GET("/conversation/sync", middleware.LTMValidationMiddleware(redisClient, true), conversation.SyncConversation)
	router.GET("/conversation/search", middleware.LTMValidationMiddleware(redisClient, true), conversation.ConversationSearch)

	// Feed Apis
	router.POST("/feed/post", middleware.LTMValidationMiddleware(redisClient, false), feed.CreatePost)
	router.GET("/feed/post/:post_id", middleware.LTMValidationMiddleware(redisClient, true), feed.GetPost)
	router.PUT("/feed/post/:post_id", middleware.LTMValidationMiddleware(redisClient, false), feed.EditPost)
	router.DELETE("/feed/post/:post_id", middleware.LTMValidationMiddleware(redisClient, false), feed.DeletePost)
	router.PUT("/feed/post/:post_id/like", middleware.LTMValidationMiddleware(redisClient, false), feed.CreatePostLike)
	router.GET("/feed/post/:post_id/like", middleware.LTMValidationMiddleware(redisClient, false), feed.GetPostLikes)
	router.PUT("/feed/post/:post_id/pin", middleware.LTMValidationMiddleware(redisClient, false), feed.PinPost)
	router.PUT("/feed/post/:post_id/save", middleware.LTMValidationMiddleware(redisClient, false), feed.CreateSavePost)
	router.POST("/feed/post/:post_id/comment", middleware.LTMValidationMiddleware(redisClient, false), feed.CommentPost)
	router.PUT("/feed/post/:post_id/comment/:comment_id", middleware.LTMValidationMiddleware(redisClient, false), feed.EditCommentPost)
	router.POST("/feed/post/:post_id/comment/:comment_id/comment", middleware.LTMValidationMiddleware(redisClient, false), feed.CreateCommentReply)
	router.GET("/feed/post/:post_id/comment/:comment_id", middleware.LTMValidationMiddleware(redisClient, true), feed.GetComment)
	router.DELETE("/feed/post/:post_id/comment/:comment_id", middleware.LTMValidationMiddleware(redisClient, false), feed.DeleteComment)
	router.PUT("/feed/post/:post_id/comment/:comment_id/like", middleware.LTMValidationMiddleware(redisClient, false), feed.CreateCommentLike)
	router.GET("/feed/post/:post_id/comment/:comment_id/like", middleware.LTMValidationMiddleware(redisClient, false), feed.GetCommentLikes)
	router.GET("/feed/user/:user_id/save", middleware.LTMValidationMiddleware(redisClient, false), feed.GetSavedPosts)
	router.GET("/feed/user/:user_id/post", middleware.LTMValidationMiddleware(redisClient, false), feed.FetchUserCreatedPosts)
	router.GET("/feed/user/:user_id/comment", middleware.LTMValidationMiddleware(redisClient, false), feed.GetUserComments)
	router.GET("/feed/user/activity", middleware.LTMValidationMiddleware(redisClient, false), feed.GetUserActivity)
	router.GET("/feed/user/:user_id/activity", middleware.LTMValidationMiddleware(redisClient, false), feed.FetchUserProfileActivity)
	router.POST("/feed/user/:user_id/activity", middleware.LTMValidationMiddleware(redisClient, false), feed.CreateUserActivity)
	router.GET("/feed/user/activity/unread_count", middleware.LTMValidationMiddleware(redisClient, false), feed.GetUserActivityUnreadCount)
	router.POST("/feed/user/activity/:activity_id/mark_read", middleware.LTMValidationMiddleware(redisClient, false), feed.UserActivityMarkRead)
	router.GET("/feed/user/:user_id/meta", middleware.LTMValidationMiddleware(redisClient, false), feed.GetUserFeedMeta)
	router.GET("/feed/universal", middleware.LTMValidationMiddleware(redisClient, true), feed.FetchUniversalFeed)
	router.GET("/feed/group", middleware.LTMValidationMiddleware(redisClient, false), feed.FetchGroupFeed)
	router.POST("/feed/topic", middleware.LTMValidationMiddleware(redisClient, false), feed.CreateTopics)
	router.GET("/feed/topic", middleware.LTMValidationMiddleware(redisClient, true), feed.GetTopic)
	router.DELETE("/feed/topic", middleware.LTMValidationMiddleware(redisClient, false), feed.DeleteTopics)
	router.PUT("/feed/topic/:topic_id", middleware.LTMValidationMiddleware(redisClient, false), feed.EditTopic)
	router.GET("/feed/connection", middleware.LTMValidationMiddleware(redisClient, false), feed.GetConnectionFeed)
	router.POST("feed/post/pending", middleware.LTMValidationMiddleware(redisClient, false), feed.CreatePendingPost)
	router.GET("/feed/user/topics", middleware.LTMValidationMiddleware(redisClient, false), feed.FetchUsersTopics)
	router.PATCH("/feed/user/:uuid/topics", middleware.LTMValidationMiddleware(redisClient, false), feed.UpdateUserTopics)

	// Utility Apis
	router.GET("/helper/url", middleware.LTMValidationMiddleware(redisClient, true), utility.DecodeUrl)
	router.POST("/helper/media/upload", middleware.LTMValidationMiddleware(redisClient, true), utility.UploadFiles)

	// Feedroom Apis
	router.POST("/feedroom", middleware.LTMValidationMiddleware(redisClient, true), feedroom.CreateFeedroom)
	router.PUT("/feedroom", middleware.LTMValidationMiddleware(redisClient, true), feedroom.EditFeedroom)
	router.DELETE("/feedroom", middleware.LTMValidationMiddleware(redisClient, true), feedroom.DeleteFeedroom)
	router.GET("/feedroom", middleware.LTMValidationMiddleware(redisClient, true), feedroom.GetFeedroom)
	router.GET("/feedroom/action", middleware.LTMValidationMiddleware(redisClient, true), feedroom.GetFeedroomMenu)
	router.GET("/feedroom/settings", middleware.LTMValidationMiddleware(redisClient, true), feedroom.GetFeedroomSettings)
	router.PUT("/feedroom/type", middleware.LTMValidationMiddleware(redisClient, true), feedroom.ChangeFeedroomType)
	router.GET("/feedroom/type", middleware.LTMValidationMiddleware(redisClient, true), feedroom.GetFeedroomTypeStatus)
	router.PUT("/feedroom/enable_member_post", middleware.LTMValidationMiddleware(redisClient, true), feedroom.EnableMemberPost)
	router.PUT("/feedroom/pin", middleware.LTMValidationMiddleware(redisClient, true), feedroom.PinFeedroom)
	router.PUT("/feedroom/auto_join_members", middleware.LTMValidationMiddleware(redisClient, true), feedroom.AutoJoinMembers)
	router.POST("/feedroom/participants", middleware.LTMValidationMiddleware(redisClient, true), feedroom.AddParticipants)
	router.GET("/feedroom/participants", middleware.LTMValidationMiddleware(redisClient, true), feedroom.GetParticipants)
	router.DELETE("/feedroom/participants", middleware.LTMValidationMiddleware(redisClient, true), feedroom.RemoveParticipants)
	router.GET("/feedroom/cohort/access", middleware.LTMValidationMiddleware(redisClient, true), feedroom.GetCohortAccess)
	router.PUT("/feedroom/cohort/access", middleware.LTMValidationMiddleware(redisClient, true), feedroom.EditCohortAccess)
	router.GET("/feedroom/mine", middleware.LTMValidationMiddleware(redisClient, true), feedroom.MyFeedrooms)
	router.PUT("/feedroom/follow", middleware.LTMValidationMiddleware(redisClient, true), feedroom.FeedroomFollow)
	router.GET("/feedroom/tag", middleware.LTMValidationMiddleware(redisClient, true), feedroom.GetTaggingList)
	router.GET("/feedroom/:feedroom_id/tag", middleware.LTMValidationMiddleware(redisClient, true), feedroom.GetTaggingList)

	// Channel Apis
	router.GET("/channel", middleware.LTMValidationMiddleware(redisClient, true), channel.FetchChannel)
	router.GET("/channel/invites", middleware.LTMValidationMiddleware(redisClient, true), channel.GetChannelInvites)
	router.PUT("/channel/invite", middleware.LTMValidationMiddleware(redisClient, true), channel.UpdateChannelInvite)
	router.GET("/channel/:channel_id/settings/member/:participant_uuid", middleware.LTMValidationMiddleware(redisClient, true), channel.GetUserChannelSettings)
	router.PUT("/channel/:channel_id/settings/member/:participant_uuid", middleware.LTMValidationMiddleware(redisClient, true), channel.UpdateUserChannelSettings)

	// Search Apis
	router.GET("/search/channel", middleware.LTMValidationMiddleware(redisClient, true), search.ChannelSearch)
	router.GET("/search/message", middleware.LTMValidationMiddleware(redisClient, true), search.MessageSearch)
	router.GET("/search/post", middleware.LTMValidationMiddleware(redisClient, true), search.PostSearch)
	router.GET("/search/post/user/:user_id", middleware.LTMValidationMiddleware(redisClient, true), search.UserCreatedPostSearch)
	router.GET("/search", middleware.LTMValidationMiddleware(redisClient, true), search.GeneralSearch)

	// Widget Apis
	router.POST("/widget", middleware.LTMValidationMiddleware(redisClient, true), widget.CreateWidget)
	router.GET("/widget", middleware.LTMValidationMiddleware(redisClient, true), widget.GetWidget)
	router.PUT("/widget/:widget_id", middleware.LTMValidationMiddleware(redisClient, true), widget.EditWidget)

	// Poll Apis
	router.PUT("/poll/:poll_id", middleware.LTMValidationMiddleware(redisClient, true), poll.AddPollOption)
	router.PUT("/poll/:poll_id/vote", middleware.LTMValidationMiddleware(redisClient, true), poll.CreatePollVote)
	router.GET("/poll/:poll_id/vote", middleware.LTMValidationMiddleware(redisClient, true), poll.GetPollVotes)

	// Webhook Apis
	router.POST("/webhook", middleware.LTMValidationMiddleware(redisClient, true), webhook.CreateWebhook)
	router.GET("/webhook", middleware.LTMValidationMiddleware(redisClient, true), webhook.GetWebhooks)
	router.GET("/webhook/:webhook_id", middleware.LTMValidationMiddleware(redisClient, true), webhook.GetWebhook)
	router.PATCH("/webhook/:webhook_id", middleware.LTMValidationMiddleware(redisClient, true), webhook.EditWebhook)
	router.DELETE("/webhook/:webhook_id", middleware.LTMValidationMiddleware(redisClient, true), webhook.DeleteWebhook)

	// Logging Apis
	router.POST("/logs", middleware.LTMValidationMiddleware(redisClient, true), logger.PushLogs)

	// Internal Apis
	router.DELETE("/cache", middleware.InternalServiceValidationMiddleware(), internalServices.DeleteCache)

	logging.Info(fmt.Sprintf("application version: %s", AppVersion))
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	logging.Fatal(router.Run(":8080"))
}

func initGin() {
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
}

func enableCors() cors.Config {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders(
		"x-member-id",
		"x-platform-code",
		"x-platform-type",
		"x-version-code",
		"x-sdk-source",
		"x-accept-version",
		"x-username",
		"x-password",
		"x-device-id",
		"x-api-key",
		"Authorization",
	)
	return config
}

// getPrometheusMetricService returns prometheus metrics service
func getPrometheusMetricService() *monitoring.PrometheusService {
	prometheusService, err := monitoring.NewPrometheusService()
	if err != nil {
		logging.Fatal(err.Error())
		return nil
	}
	return prometheusService
}
