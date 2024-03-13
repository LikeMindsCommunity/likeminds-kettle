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
	var AppVersion string = "2.22.0"

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
	router.POST("/user/merge_account", middleware.LTMValidationMiddleware(redisClient), user.MergeAccount)
	router.GET("/user/config", middleware.LTMValidationMiddleware(redisClient), user.Config)
	router.GET("/user/bot", middleware.LTMValidationMiddleware(redisClient), user.GetBot)
	router.POST("/user/device/push", middleware.LTMValidationMiddleware(redisClient), user.PushUserToken)
	router.POST("/user/subscription/whatsapp", user.WASubscription)
	router.GET("/user/meta", middleware.LTMValidationMiddleware(redisClient), user.UserMeta)
	router.POST("/user/otp", middleware.OTMValidationMiddleware(), user.GenerateUserOTP)
	router.GET("/user/otp/verify", middleware.OTMValidationMiddleware(), user.VerifyUserOTP)
	router.GET("/user/social/login", middleware.OTMValidationMiddleware(), user.UserSocialLogin)

	// Home Apis
	router.POST("/home/fetch_communities", middleware.LTMValidationMiddleware(redisClient), home.FetchCommunities)
	router.GET("/home/dm/meta", middleware.LTMValidationMiddleware(redisClient), home.DMHome)

	// SDK Apis
	router.POST("/sdk/initiate", middleware.VTMValidationMiddleware(false), sdk.InitiateSDK)
	router.POST("/sdk/project", middleware.LTMValidationMiddleware(redisClient), sdk.CreateProject)
	router.GET("/sdk/project", middleware.LTMValidationMiddleware(redisClient), sdk.GetProject)
	router.PUT("/sdk/project", middleware.LTMValidationMiddleware(redisClient), sdk.EditProject)
	router.DELETE("/sdk/project", middleware.LTMValidationMiddleware(redisClient), sdk.DeleteProject)
	router.GET("/sdk/onboarding", middleware.OTMValidationMiddleware(), sdk.GetScreen)
	router.POST("/sdk/onboarding", middleware.LTMValidationMiddleware(redisClient), sdk.CreateScreen)
	router.PUT("/sdk/onboarding", middleware.LTMValidationMiddleware(redisClient), sdk.EditScreen)
	router.DELETE("/sdk/onboarding", middleware.LTMValidationMiddleware(redisClient), sdk.DeleteScreen)

	// Chatroom Apis
	router.GET("/chatroom", middleware.LTMValidationMiddleware(redisClient), chatroom.GetChatroom)
	router.POST("/chatroom", middleware.LTMValidationMiddleware(redisClient), chatroom.CreateChatroom)
	router.PUT("/chatroom", middleware.LTMValidationMiddleware(redisClient), chatroom.EditChatroom)
	router.DELETE("/chatroom", middleware.LTMValidationMiddleware(redisClient), chatroom.DeleteChatroom)
	router.GET("/chatroom/type", middleware.LTMValidationMiddleware(redisClient), chatroom.GetChatroomTypeStatus)
	router.PUT("/chatroom/type", middleware.LTMValidationMiddleware(redisClient), chatroom.ChangeChatroomType)
	router.POST("/chatroom/schedule_follow", middleware.LTMValidationMiddleware(redisClient), chatroom.ScheduleFollow)
	router.PUT("/chatroom/pin", middleware.LTMValidationMiddleware(redisClient), chatroom.PinChatroom)
	router.GET("/chatroom/tag", middleware.LTMValidationMiddleware(redisClient), chatroom.GetTaggingList)
	router.GET("/chatroom/:chatroom_id/tag", middleware.LTMValidationMiddleware(redisClient), chatroom.GetTaggingList)
	router.GET("/chatroom/participants", middleware.LTMValidationMiddleware(redisClient), chatroom.GetParticipants)
	router.POST("/chatroom/participants", middleware.LTMValidationMiddleware(redisClient), chatroom.AddParticipants)
	router.DELETE("/chatroom/participants", middleware.LTMValidationMiddleware(redisClient), chatroom.RemoveParticipants)
	router.GET("/chatroom/settings", middleware.LTMValidationMiddleware(redisClient), chatroom.GetChatroomSettings)
	router.PUT("/chatroom/settings", middleware.LTMValidationMiddleware(redisClient), chatroom.EditChatroomSettings)
	router.PUT("/chatroom/enable_member_message", middleware.LTMValidationMiddleware(redisClient), chatroom.EnableMemberMessage)
	router.PUT("/chatroom/auto_follow_members", middleware.LTMValidationMiddleware(redisClient), chatroom.AutoFollowMembers)
	router.PUT("/chatroom/files", middleware.LTMValidationMiddleware(redisClient), chatroom.UpdateFiles)
	router.GET("/chatroom/mine", middleware.LTMValidationMiddleware(redisClient), chatroom.MyChatrooms)
	router.PUT("/chatroom/seen", middleware.LTMValidationMiddleware(redisClient), chatroom.CollabcardSeen)
	router.PUT("/chatroom/follow", middleware.LTMValidationMiddleware(redisClient), chatroom.ChatroomFollow)
	router.PUT("/chatroom/mute", middleware.LTMValidationMiddleware(redisClient), chatroom.MuteChatroom)
	router.PUT("/chatroom/rename", middleware.LTMValidationMiddleware(redisClient), chatroom.RenameChatroom)
	router.GET("/chatroom/share", middleware.LTMValidationMiddleware(redisClient), chatroom.FetchShareUrl)
	router.GET("/chatroom/pending", middleware.LTMValidationMiddleware(redisClient), chatroom.FetchPendingChatroom)
	router.PUT("/chatroom/pending", middleware.LTMValidationMiddleware(redisClient), chatroom.ActionPendingChatroom)
	router.GET("/chatroom/sync", middleware.LTMValidationMiddleware(redisClient), chatroom.SyncChatrooms)
	router.POST("/chatroom/dm/block", middleware.LTMValidationMiddleware(redisClient), chatroom.ChatroomBlock)
	router.POST("/chatroom/dm/request", middleware.LTMValidationMiddleware(redisClient), chatroom.InitiatingDMRequest)
	router.POST("/chatroom/dm/create", middleware.LTMValidationMiddleware(redisClient), chatroom.CreateDM)
	router.GET("/chatroom/dm", middleware.LTMValidationMiddleware(redisClient), chatroom.ListDMChatrooms)
	router.GET("/chatroom/dm/limit", middleware.LTMValidationMiddleware(redisClient), chatroom.DMLimit)
	router.GET("/chatroom/search", middleware.LTMValidationMiddleware(redisClient), chatroom.ChatroomSearch)
	router.POST("/chatroom/cohort", middleware.LTMValidationMiddleware(redisClient), chatroom.AddCohortToChatroom)
	router.DELETE("/chatroom/cohort", middleware.LTMValidationMiddleware(redisClient), chatroom.RemoveCohortFromChatroom)
	router.GET("/chatroom/cohort/access", middleware.LTMValidationMiddleware(redisClient), chatroom.GetCohortAccess)
	router.PUT("/chatroom/cohort/access", middleware.LTMValidationMiddleware(redisClient), chatroom.EditCohortAccess)
	router.GET("/chatroom/home", middleware.LTMValidationMiddleware(redisClient), chatroom.GetChatroomHome)
	router.POST("/chatroom/mark_read", middleware.LTMValidationMiddleware(redisClient), chatroom.ChatroomMarkRead)
	router.GET("/chatroom/event", middleware.LTMValidationMiddleware(redisClient), chatroom.FetchEvents)
	router.POST("/chatroom/event", middleware.LTMValidationMiddleware(redisClient), chatroom.CreateEvent)
	router.PUT("/chatroom/event", middleware.LTMValidationMiddleware(redisClient), chatroom.EditEvent)
	router.GET("/chatroom/event/meta", middleware.LTMValidationMiddleware(redisClient), chatroom.FetchEventMeta)
	router.GET("/chatroom/event/link", middleware.LTMValidationMiddleware(redisClient), chatroom.FetchEventLinks)
	router.GET("/chatroom/event/unseen_count", middleware.LTMValidationMiddleware(redisClient), chatroom.FetchEventUnseenCount)
	router.POST("/chatroom/event/recordings", middleware.LTMValidationMiddleware(redisClient), chatroom.UploadEventRecordings)
	router.DELETE("/chatroom/event/recordings", middleware.LTMValidationMiddleware(redisClient), chatroom.DeleteEventRecordings)
	router.POST("/chatroom/event/recordings/meta", middleware.LTMValidationMiddleware(redisClient), chatroom.UploadEventRecordingsMeta)
	router.DELETE("/chatroom/event/recordings/meta", middleware.LTMValidationMiddleware(redisClient), chatroom.DeleteEventRecordingsMeta)
	router.POST("/chatroom/event/instructors", middleware.LTMValidationMiddleware(redisClient), chatroom.AddEventInstructors)
	router.POST("/chatroom/event/highlights", middleware.LTMValidationMiddleware(redisClient), chatroom.AddEventHighlights)
	router.POST("/chatroom/event/testimonials", middleware.LTMValidationMiddleware(redisClient), chatroom.AddEventTestimonials)
	router.POST("/chatroom/event/faq", middleware.LTMValidationMiddleware(redisClient), chatroom.AddEventFAQ)

	// Community Apis
	router.GET("/community", middleware.LTMValidationMiddleware(redisClient), community.Community)
	router.GET("/community/branding", middleware.LTMValidationMiddleware(redisClient), community.Branding)
	router.POST("/community/questions", middleware.LTMValidationMiddleware(redisClient), community.EditQuestions)
	router.GET("/community/questions", middleware.LTMorVTMValidationMiddleware(), community.GetQuestions)
	router.GET("/community/question/filters", middleware.LTMValidationMiddleware(redisClient), community.GetCommunityQuestionFilters)
	router.GET("/community/member", middleware.LTMValidationMiddleware(redisClient), community.GetMember)
	router.POST("/community/member", middleware.LTMValidationMiddleware(redisClient), community.AddMember)
	router.DELETE("/community/member", middleware.LTMValidationMiddleware(redisClient), community.RemoveMembers)
	router.DELETE("/community/member/leave", middleware.LTMValidationMiddleware(redisClient), community.LeaveCommunity)
	router.PUT("/community/member", middleware.LTMValidationMiddleware(redisClient), community.EditMember)
	router.GET("/community/member/state", middleware.LTMValidationMiddleware(redisClient), community.FetchMemberState)
	router.GET("/community/member/role", middleware.LTMValidationMiddleware(redisClient), community.FetchMemberRole)
	router.DELETE("/community/manager/remove", middleware.LTMValidationMiddleware(redisClient), community.RemoveCommunityManager)
	router.DELETE("/community/admin", middleware.LTMValidationMiddleware(redisClient), community.RemoveCommunityManager)
	router.DELETE("/community/member/remove", middleware.LTMValidationMiddleware(redisClient), community.RemoveMember)
	router.GET("/community/management/tool", middleware.LTMValidationMiddleware(redisClient), community.GetManagementTools)
	router.GET("/community/report", middleware.LTMValidationMiddleware(redisClient), community.GetReport)
	router.POST("/community/report", middleware.LTMValidationMiddleware(redisClient), community.PushReport)
	router.DELETE("/community/report", middleware.LTMValidationMiddleware(redisClient), community.CloseReport)
	router.PATCH("/community/report", middleware.LTMValidationMiddleware(redisClient), community.UpdateReports)
	router.GET("/community/report/tag", middleware.LTMValidationMiddleware(redisClient), community.GetReportTags)
	router.GET("/community/settings", middleware.LTMValidationMiddleware(redisClient), community.GetCommunitySettings)
	router.PUT("/community/settings", middleware.LTMValidationMiddleware(redisClient), community.UpdateCommunitySettings)
	router.GET("/community/rights", middleware.LTMValidationMiddleware(redisClient), community.GetCommunityRights)
	router.PUT("/community/rights", middleware.LTMValidationMiddleware(redisClient), community.EditCommunityRights)
	router.PATCH("/community/rights", middleware.LTMValidationMiddleware(redisClient), community.UpdateCommunityRights)
	router.GET("/community/settings/dm", middleware.LTMValidationMiddleware(redisClient), community.GetCommunityDMSettings)
	router.PUT("/community/settings/dm", middleware.LTMValidationMiddleware(redisClient), community.EditCommunityDMSettings)
	router.GET("/community/feed/dm", middleware.LTMValidationMiddleware(redisClient), community.DMFeed)
	router.GET("/community/dm/status", middleware.LTMValidationMiddleware(redisClient), community.DMStatus)
	router.GET("/community/member/search", middleware.LTMValidationMiddleware(redisClient), community.MemberSearch)
	router.GET("/community/member/profile", middleware.LTMValidationMiddleware(redisClient), community.GetMemberProfile)
	router.PUT("/community/member/profile", middleware.LTMValidationMiddleware(redisClient), community.EditMemberProfile)
	router.GET("/community/member/chatroom", middleware.LTMValidationMiddleware(redisClient), community.MemberChatroom)
	router.GET("/community/member/:user_id/channel", middleware.LTMValidationMiddleware(redisClient), community.CommunityMemberChannels)
	router.GET("/community/member/channel/status", middleware.LTMValidationMiddleware(redisClient), community.GetMemberChannels)
	router.POST("/community/cohort", middleware.LTMValidationMiddleware(redisClient), community.CreateCohort)
	router.GET("/community/cohort", middleware.LTMValidationMiddleware(redisClient), community.GetCohort)
	router.GET("/community/cohort/:cohort_id", middleware.LTMValidationMiddleware(redisClient), community.FetchCohort)
	router.DELETE("/community/cohort", middleware.LTMValidationMiddleware(redisClient), community.DeleteCohort)
	router.PUT("/community/cohort", middleware.LTMValidationMiddleware(redisClient), community.EditCohort)
	router.DELETE("/community/cohort/member", middleware.LTMValidationMiddleware(redisClient), community.RemoveCohortMember)
	router.GET("/community/feed", middleware.LTMValidationMiddleware(redisClient), community.GetCommunityFeed)
	router.GET("/community/settings/notification/conversation", middleware.LTMValidationMiddleware(redisClient), community.GetConversationNotificationSettings)
	router.PUT("/community/settings/notification/conversation", middleware.LTMValidationMiddleware(redisClient), community.EditConversationNotificationSettings)
	router.GET("/community/settings/notification/feed", middleware.LTMValidationMiddleware(redisClient), community.GetFeedNotificationSettings)
	router.PUT("/community/settings/notification/feed", middleware.LTMValidationMiddleware(redisClient), community.EditFeedNotificationSettings)
	router.GET("/community/settings/notification", middleware.LTMValidationMiddleware(redisClient), community.GetNotificationSettings)
	router.PUT("/community/settings/notification", middleware.LTMValidationMiddleware(redisClient), community.EditNotificationSettings)
	router.GET("/community/tag", middleware.LTMValidationMiddleware(redisClient), community.GetTaggingList)
	router.GET("/community/settings/content_download", middleware.LTMValidationMiddleware(redisClient), community.GetContentDownloadSettings)
	router.PUT("/community/settings/content_download", middleware.LTMValidationMiddleware(redisClient), community.EditContentDownloadSettings)
	router.GET("/community/member/home/meta", middleware.LTMValidationMiddleware(redisClient), community.MemberHomeMeta)
	router.PUT("/community/member/join", middleware.LTMValidationMiddleware(redisClient), community.AcceptRejectJoinCommunity)
	router.GET("/community/intro_examples", middleware.LTMorVTMValidationMiddleware(), community.GetIntroExamples)
	router.POST("/community/invite", middleware.LTMValidationMiddleware(redisClient), community.SendCommunityInvite)
	router.GET("/community/configurations", middleware.LTMValidationMiddleware(redisClient), community.GetCommunityConfigurations)
	router.PATCH("/community/configurations", middleware.LTMValidationMiddleware(redisClient), community.UpdateCommunityConfigurations)
	router.GET("/community/member/pending", middleware.LTMValidationMiddleware(redisClient), community.GetPendingCommunityMembers)
	router.GET("/community/removal_reports", middleware.LTMValidationMiddleware(redisClient), community.GetRemovalReports)
	router.POST("/community/member/:user_id/connection", middleware.LTMValidationMiddleware(redisClient), community.CreateMemberConnection)
	router.PATCH("/community/member/:user_id/connection", middleware.LTMValidationMiddleware(redisClient), community.AcceptRejectMemberConnection)
	router.GET("/community/member/:user_id/connection", middleware.LTMValidationMiddleware(redisClient), community.GetMemberConnection)

	// Moderation Apis
	router.GET("/moderation/rights", middleware.LTMValidationMiddleware(redisClient), moderation.GetRights)
	router.PUT("/moderation/rights", middleware.LTMValidationMiddleware(redisClient), moderation.EditRights)
	router.PATCH("/moderation/rights", middleware.LTMValidationMiddleware(redisClient), moderation.UpdateRights)

	// Conversation Apis
	router.GET("/conversation", middleware.LTMValidationMiddleware(redisClient), conversation.GetConversation)
	router.POST("/conversation", middleware.LTMValidationMiddleware(redisClient), conversation.CreateConversation)
	router.PUT("/conversation", middleware.LTMValidationMiddleware(redisClient), conversation.EditConversation)
	router.DELETE("/conversation", middleware.LTMValidationMiddleware(redisClient), conversation.DeleteConversation)
	router.PUT("/conversation/reaction", middleware.LTMValidationMiddleware(redisClient), conversation.AddReaction)
	router.DELETE("/conversation/reaction", middleware.LTMValidationMiddleware(redisClient), conversation.RemoveReaction)
	router.POST("/conversation/poll", middleware.LTMValidationMiddleware(redisClient), conversation.AddPoll)
	router.POST("/conversation/poll/submit", middleware.LTMValidationMiddleware(redisClient), conversation.SubmitPoll)
	router.GET("/conversation/poll/users", middleware.LTMValidationMiddleware(redisClient), conversation.PollUsers)
	router.PUT("/conversation/topic", middleware.LTMValidationMiddleware(redisClient), conversation.SetTopic)
	router.PUT("/conversation/event/attend", middleware.LTMValidationMiddleware(redisClient), conversation.EventAttend)
	router.PUT("/conversation/event/attended", middleware.LTMValidationMiddleware(redisClient), conversation.EventAttended)
	router.GET("/conversation/event/unseen", middleware.LTMValidationMiddleware(redisClient), conversation.FetchEventUnseenCount)
	router.PUT("/conversation/event/last_seen", middleware.LTMValidationMiddleware(redisClient), conversation.UpdateLastSeenEvent)
	router.GET("/conversation/event/link", middleware.LTMValidationMiddleware(redisClient), conversation.FetchEventLink)
	router.GET("/conversation/event", middleware.LTMValidationMiddleware(redisClient), conversation.FetchAllEvents)
	router.GET("/conversation/preview/unread", middleware.LTMValidationMiddleware(redisClient), conversation.FetchUnreadPreviews)
	router.GET("/conversation/preview/unread_count", middleware.LTMValidationMiddleware(redisClient), conversation.FetchPreviewUnreadMessagesCount)
	router.GET("/conversation/notification/unread", middleware.LTMValidationMiddleware(redisClient), conversation.UnreadConversationNotification)
	router.GET("/conversation/sync", middleware.LTMValidationMiddleware(redisClient), conversation.SyncConversation)
	router.GET("/conversation/search", middleware.LTMValidationMiddleware(redisClient), conversation.ConversationSearch)

	// Feed Apis
	router.POST("/feed/post", middleware.LTMValidationMiddleware(redisClient), feed.CreatePost)
	router.GET("/feed/post/:post_id", middleware.LTMValidationMiddleware(redisClient), feed.GetPost)
	router.PUT("/feed/post/:post_id", middleware.LTMValidationMiddleware(redisClient), feed.EditPost)
	router.DELETE("/feed/post/:post_id", middleware.LTMValidationMiddleware(redisClient), feed.DeletePost)
	router.PUT("/feed/post/:post_id/like", middleware.LTMValidationMiddleware(redisClient), feed.CreatePostLike)
	router.GET("/feed/post/:post_id/like", middleware.LTMValidationMiddleware(redisClient), feed.GetPostLikes)
	router.PUT("/feed/post/:post_id/pin", middleware.LTMValidationMiddleware(redisClient), feed.PinPost)
	router.PUT("/feed/post/:post_id/save", middleware.LTMValidationMiddleware(redisClient), feed.CreateSavePost)
	router.POST("/feed/post/:post_id/comment", middleware.LTMValidationMiddleware(redisClient), feed.CommentPost)
	router.PUT("/feed/post/:post_id/comment/:comment_id", middleware.LTMValidationMiddleware(redisClient), feed.EditCommentPost)
	router.POST("/feed/post/:post_id/comment/:comment_id/comment", middleware.LTMValidationMiddleware(redisClient), feed.CreateCommentReply)
	router.GET("/feed/post/:post_id/comment/:comment_id", middleware.LTMValidationMiddleware(redisClient), feed.GetComment)
	router.DELETE("/feed/post/:post_id/comment/:comment_id", middleware.LTMValidationMiddleware(redisClient), feed.DeleteComment)
	router.PUT("/feed/post/:post_id/comment/:comment_id/like", middleware.LTMValidationMiddleware(redisClient), feed.CreateCommentLike)
	router.GET("/feed/post/:post_id/comment/:comment_id/like", middleware.LTMValidationMiddleware(redisClient), feed.GetCommentLikes)
	router.GET("/feed/user/:user_id/save", middleware.LTMValidationMiddleware(redisClient), feed.GetSavedPosts)
	router.GET("/feed/user/:user_id/post", middleware.LTMValidationMiddleware(redisClient), feed.FetchUserCreatedPosts)
	router.GET("/feed/user/:user_id/comment", middleware.LTMValidationMiddleware(redisClient), feed.GetUserComments)
	router.GET("/feed/user/activity", middleware.LTMValidationMiddleware(redisClient), feed.GetUserActivity)
	router.GET("/feed/user/:user_id/activity", middleware.LTMValidationMiddleware(redisClient), feed.FetchUserProfileActivity)
	router.POST("/feed/user/:user_id/activity", middleware.LTMValidationMiddleware(redisClient), feed.CreateUserActivity)
	router.GET("/feed/user/activity/unread_count", middleware.LTMValidationMiddleware(redisClient), feed.GetUserActivityUnreadCount)
	router.POST("/feed/user/activity/:activity_id/mark_read", middleware.LTMValidationMiddleware(redisClient), feed.UserActivityMarkRead)
	router.GET("/feed/user/:user_id/meta", middleware.LTMValidationMiddleware(redisClient), feed.GetUserFeedMeta)
	router.GET("/feed/universal", middleware.LTMValidationMiddleware(redisClient), feed.FetchUniversalFeed)
	router.GET("/feed/group", middleware.LTMValidationMiddleware(redisClient), feed.FetchGroupFeed)
	router.POST("/feed/topic", middleware.LTMValidationMiddleware(redisClient), feed.CreateTopics)
	router.GET("/feed/topic", middleware.LTMValidationMiddleware(redisClient), feed.GetTopic)
	router.DELETE("/feed/topic", middleware.LTMValidationMiddleware(redisClient), feed.DeleteTopics)
	router.PUT("/feed/topic/:topic_id", middleware.LTMValidationMiddleware(redisClient), feed.EditTopic)
	router.GET("/feed/connection", middleware.LTMValidationMiddleware(redisClient), feed.GetConnectionFeed)
	router.POST("feed/post/pending", middleware.LTMValidationMiddleware(redisClient), feed.CreatePendingPost)
	router.GET("/feed/user/topics", middleware.LTMValidationMiddleware(redisClient), feed.FetchUsersTopics)
	router.PATCH("/feed/user/:uuid/topics", middleware.LTMValidationMiddleware(redisClient), feed.UpdateUserTopics)

	// Utility Apis
	router.GET("/helper/url", middleware.LTMValidationMiddleware(redisClient), utility.DecodeUrl)
	router.POST("/helper/media/upload", middleware.LTMValidationMiddleware(redisClient), utility.UploadFiles)

	// Feedroom Apis
	router.POST("/feedroom", middleware.LTMValidationMiddleware(redisClient), feedroom.CreateFeedroom)
	router.PUT("/feedroom", middleware.LTMValidationMiddleware(redisClient), feedroom.EditFeedroom)
	router.DELETE("/feedroom", middleware.LTMValidationMiddleware(redisClient), feedroom.DeleteFeedroom)
	router.GET("/feedroom", middleware.LTMValidationMiddleware(redisClient), feedroom.GetFeedroom)
	router.GET("/feedroom/action", middleware.LTMValidationMiddleware(redisClient), feedroom.GetFeedroomMenu)
	router.GET("/feedroom/settings", middleware.LTMValidationMiddleware(redisClient), feedroom.GetFeedroomSettings)
	router.PUT("/feedroom/type", middleware.LTMValidationMiddleware(redisClient), feedroom.ChangeFeedroomType)
	router.GET("/feedroom/type", middleware.LTMValidationMiddleware(redisClient), feedroom.GetFeedroomTypeStatus)
	router.PUT("/feedroom/enable_member_post", middleware.LTMValidationMiddleware(redisClient), feedroom.EnableMemberPost)
	router.PUT("/feedroom/pin", middleware.LTMValidationMiddleware(redisClient), feedroom.PinFeedroom)
	router.PUT("/feedroom/auto_join_members", middleware.LTMValidationMiddleware(redisClient), feedroom.AutoJoinMembers)
	router.POST("/feedroom/participants", middleware.LTMValidationMiddleware(redisClient), feedroom.AddParticipants)
	router.GET("/feedroom/participants", middleware.LTMValidationMiddleware(redisClient), feedroom.GetParticipants)
	router.DELETE("/feedroom/participants", middleware.LTMValidationMiddleware(redisClient), feedroom.RemoveParticipants)
	router.GET("/feedroom/cohort/access", middleware.LTMValidationMiddleware(redisClient), feedroom.GetCohortAccess)
	router.PUT("/feedroom/cohort/access", middleware.LTMValidationMiddleware(redisClient), feedroom.EditCohortAccess)
	router.GET("/feedroom/mine", middleware.LTMValidationMiddleware(redisClient), feedroom.MyFeedrooms)
	router.PUT("/feedroom/follow", middleware.LTMValidationMiddleware(redisClient), feedroom.FeedroomFollow)
	router.GET("/feedroom/tag", middleware.LTMValidationMiddleware(redisClient), feedroom.GetTaggingList)
	router.GET("/feedroom/:feedroom_id/tag", middleware.LTMValidationMiddleware(redisClient), feedroom.GetTaggingList)

	// Channel Apis
	router.GET("/channel", middleware.LTMValidationMiddleware(redisClient), channel.FetchChannel)
	router.GET("/channel/invites", middleware.LTMValidationMiddleware(redisClient), channel.GetChannelInvites)
	router.PUT("/channel/invite", middleware.LTMValidationMiddleware(redisClient), channel.UpdateChannelInvite)
	router.GET("/channel/:channel_id/settings/member/:participant_uuid", middleware.LTMValidationMiddleware(redisClient), channel.GetUserChannelSettings)
	router.PUT("/channel/:channel_id/settings/member/:participant_uuid", middleware.LTMValidationMiddleware(redisClient), channel.UpdateUserChannelSettings)

	// Search Apis
	router.GET("/search/channel", middleware.LTMValidationMiddleware(redisClient), search.ChannelSearch)
	router.GET("/search/message", middleware.LTMValidationMiddleware(redisClient), search.MessageSearch)
	router.GET("/search/post", middleware.LTMValidationMiddleware(redisClient), search.PostSearch)
	router.GET("/search/post/user/:user_id", middleware.LTMValidationMiddleware(redisClient), search.UserCreatedPostSearch)
	router.GET("/search", middleware.LTMValidationMiddleware(redisClient), search.GeneralSearch)

	// Widget Apis
	router.POST("/widget", middleware.LTMValidationMiddleware(redisClient), widget.CreateWidget)
	router.GET("/widget", middleware.LTMValidationMiddleware(redisClient), widget.GetWidget)
	router.PUT("/widget/:widget_id", middleware.LTMValidationMiddleware(redisClient), widget.EditWidget)

	// Poll Apis
	router.PUT("/poll/:poll_id", middleware.LTMValidationMiddleware(redisClient), poll.AddPollOption)
	router.PUT("/poll/:poll_id/vote", middleware.LTMValidationMiddleware(redisClient), poll.CreatePollVote)
	router.GET("/poll/:poll_id/vote", middleware.LTMValidationMiddleware(redisClient), poll.GetPollVotes)

	// Webhook Apis
	router.POST("/webhook", middleware.LTMValidationMiddleware(redisClient), webhook.CreateWebhook)
	router.GET("/webhook", middleware.LTMValidationMiddleware(redisClient), webhook.GetWebhooks)
	router.GET("/webhook/:webhook_id", middleware.LTMValidationMiddleware(redisClient), webhook.GetWebhook)
	router.PATCH("/webhook/:webhook_id", middleware.LTMValidationMiddleware(redisClient), webhook.EditWebhook)
	router.DELETE("/webhook/:webhook_id", middleware.LTMValidationMiddleware(redisClient), webhook.DeleteWebhook)

	// Logging Apis
	router.POST("/logs", middleware.LTMValidationMiddleware(redisClient), logger.PushLogs)

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
