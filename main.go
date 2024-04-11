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
	var AppVersion string = "2.27.0"

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
	router.POST("/sdk/initiate", middleware.VTMValidationMiddleware(false), middleware.RateLimitingMiddleware(redisClient), sdk.InitiateSDK)
	router.GET("/sdk/initiate", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.FetchSdkUserInfo)
	router.POST("/sdk/project", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.CreateProject)
	router.GET("/sdk/project", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.GetProject)
	router.PUT("/sdk/project", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.EditProject)
	router.DELETE("/sdk/project", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.DeleteProject)
	router.GET("/sdk/onboarding", middleware.OTMValidationMiddleware(), middleware.RateLimitingMiddleware(redisClient), sdk.GetScreen)
	router.POST("/sdk/onboarding", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.CreateScreen)
	router.PUT("/sdk/onboarding", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.EditScreen)
	router.DELETE("/sdk/onboarding", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.DeleteScreen)

	// Chatroom Apis
	router.GET("/chatroom", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetChatroom)
	router.POST("/chatroom", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.CreateChatroom)
	router.PUT("/chatroom", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.EditChatroom)
	router.DELETE("/chatroom", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.DeleteChatroom)
	router.GET("/chatroom/type", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetChatroomTypeStatus)
	router.PUT("/chatroom/type", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ChangeChatroomType)
	router.POST("/chatroom/schedule_follow", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ScheduleFollow)
	router.PUT("/chatroom/pin", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.PinChatroom)
	router.GET("/chatroom/tag", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetTaggingList)
	router.GET("/chatroom/:chatroom_id/tag", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetTaggingList)
	router.GET("/chatroom/participants", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetParticipants)
	router.POST("/chatroom/participants", middleware.LTMValidationMiddleware(redisClient, true), chatroom.AddParticipants)
	router.DELETE("/chatroom/participants", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.RemoveParticipants)
	router.GET("/chatroom/settings", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetChatroomSettings)
	router.PUT("/chatroom/settings", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.EditChatroomSettings)
	router.PUT("/chatroom/enable_member_message", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.EnableMemberMessage)
	router.PUT("/chatroom/auto_follow_members", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.AutoFollowMembers)
	router.PUT("/chatroom/files", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.UpdateFiles)
	router.GET("/chatroom/mine", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.MyChatrooms)
	router.PUT("/chatroom/seen", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.CollabcardSeen)
	router.PUT("/chatroom/follow", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ChatroomFollow)
	router.PUT("/chatroom/mute", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.MuteChatroom)
	router.PUT("/chatroom/rename", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.RenameChatroom)
	router.GET("/chatroom/share", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.FetchShareUrl)
	router.GET("/chatroom/pending", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.FetchPendingChatroom)
	router.PUT("/chatroom/pending", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ActionPendingChatroom)
	router.GET("/chatroom/sync", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.SyncChatrooms)
	router.POST("/chatroom/dm/block", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ChatroomBlock)
	router.POST("/chatroom/dm/request", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.InitiatingDMRequest)
	router.POST("/chatroom/dm/create", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.CreateDM)
	router.GET("/chatroom/dm", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ListDMChatrooms)
	router.GET("/chatroom/dm/limit", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.DMLimit)
	router.GET("/chatroom/search", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ChatroomSearch)
	router.POST("/chatroom/cohort", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.AddCohortToChatroom)
	router.DELETE("/chatroom/cohort", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.RemoveCohortFromChatroom)
	router.GET("/chatroom/cohort/access", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetCohortAccess)
	router.PUT("/chatroom/cohort/access", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.EditCohortAccess)
	router.GET("/chatroom/home", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetChatroomHome)
	router.POST("/chatroom/mark_read", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ChatroomMarkRead)
	router.GET("/chatroom/event", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.FetchEvents)
	router.POST("/chatroom/event", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.CreateEvent)
	router.PUT("/chatroom/event", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.EditEvent)
	router.GET("/chatroom/event/meta", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.FetchEventMeta)
	router.GET("/chatroom/event/link", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.FetchEventLinks)
	router.GET("/chatroom/event/unseen_count", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.FetchEventUnseenCount)
	router.POST("/chatroom/event/recordings", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.UploadEventRecordings)
	router.DELETE("/chatroom/event/recordings", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.DeleteEventRecordings)
	router.POST("/chatroom/event/recordings/meta", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.UploadEventRecordingsMeta)
	router.DELETE("/chatroom/event/recordings/meta", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.DeleteEventRecordingsMeta)
	router.POST("/chatroom/event/instructors", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.AddEventInstructors)
	router.POST("/chatroom/event/highlights", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.AddEventHighlights)
	router.POST("/chatroom/event/testimonials", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.AddEventTestimonials)
	router.POST("/chatroom/event/faq", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.AddEventFAQ)

	// Community Apis
	router.GET("/community", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.Community)
	router.GET("/community/branding", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.Branding)
	router.POST("/community/questions", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditQuestions)
	router.GET("/community/questions", middleware.LTMorVTMValidationMiddleware(), middleware.RateLimitingMiddleware(redisClient), community.GetQuestions)
	router.GET("/community/question/filters", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetCommunityQuestionFilters)
	router.GET("/community/member", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetMember)
	router.POST("/community/member", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.AddMember)
	router.DELETE("/community/member", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.RemoveMembers)
	router.DELETE("/community/member/leave", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.LeaveCommunity)
	router.PUT("/community/member", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditMember)
	router.GET("/community/member/state", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.FetchMemberState)
	router.GET("/community/member/role", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.FetchMemberRole)
	router.DELETE("/community/manager/remove", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.RemoveCommunityManager)
	router.DELETE("/community/admin", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.RemoveCommunityManager)
	router.DELETE("/community/member/remove", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.RemoveMember)
	router.GET("/community/management/tool", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetManagementTools)
	router.GET("/community/report", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetReport)
	router.POST("/community/report", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.PushReport)
	router.DELETE("/community/report", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.CloseReport)
	router.PATCH("/community/report", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.UpdateReports)
	router.GET("/community/report/tag", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetReportTags)
	router.GET("/community/settings", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetCommunitySettings)
	router.PUT("/community/settings", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.UpdateCommunitySettings)
	router.GET("/community/rights", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetCommunityRights)
	router.PUT("/community/rights", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditCommunityRights)
	router.PATCH("/community/rights", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.UpdateCommunityRights)
	router.GET("/community/settings/dm", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetCommunityDMSettings)
	router.PUT("/community/settings/dm", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditCommunityDMSettings)
	router.GET("/community/feed/dm", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.DMFeed)
	router.GET("/community/dm/status", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.DMStatus)
	router.GET("/community/member/search", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.MemberSearch)
	router.GET("/community/member/profile", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetMemberProfile)
	router.PUT("/community/member/profile", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditMemberProfile)
	router.GET("/community/member/chatroom", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.MemberChatroom)
	router.GET("/community/member/:user_id/channel", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.CommunityMemberChannels)
	router.GET("/community/member/channel/status", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetMemberChannels)
	router.POST("/community/cohort", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.CreateCohort)
	router.GET("/community/cohort", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetCohort)
	router.GET("/community/cohort/:cohort_id", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.FetchCohort)
	router.DELETE("/community/cohort", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.DeleteCohort)
	router.PUT("/community/cohort", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditCohort)
	router.DELETE("/community/cohort/member", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.RemoveCohortMember)
	router.GET("/community/feed", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetCommunityFeed)
	router.GET("/community/settings/notification/conversation", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetConversationNotificationSettings)
	router.PUT("/community/settings/notification/conversation", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditConversationNotificationSettings)
	router.GET("/community/settings/notification/feed", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetFeedNotificationSettings)
	router.PUT("/community/settings/notification/feed", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditFeedNotificationSettings)
	router.GET("/community/settings/notification", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetNotificationSettings)
	router.PUT("/community/settings/notification", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditNotificationSettings)
	router.GET("/community/tag", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetTaggingList)
	router.GET("/community/settings/content_download", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetContentDownloadSettings)
	router.PUT("/community/settings/content_download", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditContentDownloadSettings)
	router.GET("/community/member/home/meta", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.MemberHomeMeta)
	router.PUT("/community/member/join", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.AcceptRejectJoinCommunity)
	router.GET("/community/intro_examples", middleware.LTMorVTMValidationMiddleware(), middleware.RateLimitingMiddleware(redisClient), community.GetIntroExamples)
	router.POST("/community/invite", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.SendCommunityInvite)
	router.GET("/community/configurations", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetCommunityConfigurations)
	router.PATCH("/community/configurations", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.UpdateCommunityConfigurations)
	router.GET("/community/member/pending", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetPendingCommunityMembers)
	router.GET("/community/removal_reports", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetRemovalReports)
	router.POST("/community/member/:user_id/connection", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.CreateMemberConnection)
	router.PATCH("/community/member/:user_id/connection", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.AcceptRejectMemberConnection)
	router.GET("/community/member/:user_id/connection", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetMemberConnection)

	// Moderation Apis
	router.GET("/moderation/rights", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), moderation.GetRights)
	router.PUT("/moderation/rights", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), moderation.EditRights)
	router.PATCH("/moderation/rights", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), moderation.UpdateRights)

	// Conversation Apis
	router.GET("/conversation", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.GetConversation)
	router.POST("/conversation", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.CreateConversation)
	router.PUT("/conversation", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.EditConversation)
	router.DELETE("/conversation", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.DeleteConversation)
	router.PUT("/conversation/reaction", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.AddReaction)
	router.DELETE("/conversation/reaction", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.RemoveReaction)
	router.POST("/conversation/poll", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.AddPoll)
	router.POST("/conversation/poll/submit", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.SubmitPoll)
	router.GET("/conversation/poll/users", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.PollUsers)
	router.PUT("/conversation/topic", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.SetTopic)
	router.PUT("/conversation/event/attend", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.EventAttend)
	router.PUT("/conversation/event/attended", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.EventAttended)
	router.GET("/conversation/event/unseen", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.FetchEventUnseenCount)
	router.PUT("/conversation/event/last_seen", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.UpdateLastSeenEvent)
	router.GET("/conversation/event/link", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.FetchEventLink)
	router.GET("/conversation/event", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.FetchAllEvents)
	router.GET("/conversation/preview/unread", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.FetchUnreadPreviews)
	router.GET("/conversation/preview/unread_count", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.FetchPreviewUnreadMessagesCount)
	router.GET("/conversation/notification/unread", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.UnreadConversationNotification)
	router.GET("/conversation/sync", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.SyncConversation)
	router.GET("/conversation/search", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.ConversationSearch)

	// Feed Apis
	router.POST("/feed/post", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreatePost)
	router.GET("/feed/post/:post_id", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feed.GetPost)
	router.PUT("/feed/post/:post_id", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.EditPost)
	router.DELETE("/feed/post/:post_id", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.DeletePost)
	router.PUT("/feed/post/:post_id/like", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreatePostLike)
	router.GET("/feed/post/:post_id/like", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetPostLikes)
	router.PUT("/feed/post/:post_id/pin", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.PinPost)
	router.PUT("/feed/post/:post_id/save", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreateSavePost)
	router.POST("/feed/post/:post_id/comment", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CommentPost)
	router.PUT("/feed/post/:post_id/comment/:comment_id", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.EditCommentPost)
	router.POST("/feed/post/:post_id/comment/:comment_id/comment", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreateCommentReply)
	router.GET("/feed/post/:post_id/comment/:comment_id", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feed.GetComment)
	router.DELETE("/feed/post/:post_id/comment/:comment_id", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.DeleteComment)
	router.PUT("/feed/post/:post_id/comment/:comment_id/like", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreateCommentLike)
	router.GET("/feed/post/:post_id/comment/:comment_id/like", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetCommentLikes)
	router.GET("/feed/user/:user_id/save", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetSavedPosts)
	router.GET("/feed/user/:user_id/post", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.FetchUserCreatedPosts)
	router.GET("/feed/user/:user_id/comment", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetUserComments)
	router.GET("/feed/user/activity", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetUserActivity)
	router.GET("/feed/user/:user_id/activity", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.FetchUserProfileActivity)
	router.POST("/feed/user/:user_id/activity", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreateUserActivity)
	router.GET("/feed/user/activity/unread_count", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetUserActivityUnreadCount)
	router.POST("/feed/user/activity/:activity_id/mark_read", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.UserActivityMarkRead)
	router.GET("/feed/user/:user_id/meta", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetUserFeedMeta)
	router.GET("/feed/universal", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feed.FetchUniversalFeed)
	router.GET("/feed/group", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.FetchGroupFeed)
	router.POST("/feed/topic", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreateTopics)
	router.GET("/feed/topic", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feed.GetTopic)
	router.DELETE("/feed/topic", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.DeleteTopics)
	router.PUT("/feed/topic/:topic_id", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.EditTopic)
	router.GET("/feed/connection", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetConnectionFeed)
	router.POST("feed/post/pending", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreatePendingPost)
	router.GET("/feed/user/topics", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.FetchUsersTopics)
	router.PATCH("/feed/user/:uuid/topics", middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.UpdateUserTopics)

	// Utility Apis
	router.GET("/helper/url", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), utility.DecodeUrl)
	router.POST("/helper/media/upload", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), utility.UploadFiles)
	router.PUT("/helper/s3/upload", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), utility.UploadFilesToS3)

	// Feedroom Apis
	router.POST("/feedroom", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.CreateFeedroom)
	router.PUT("/feedroom", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.EditFeedroom)
	router.DELETE("/feedroom", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.DeleteFeedroom)
	router.GET("/feedroom", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetFeedroom)
	router.GET("/feedroom/action", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetFeedroomMenu)
	router.GET("/feedroom/settings", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetFeedroomSettings)
	router.PUT("/feedroom/type", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.ChangeFeedroomType)
	router.GET("/feedroom/type", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetFeedroomTypeStatus)
	router.PUT("/feedroom/enable_member_post", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.EnableMemberPost)
	router.PUT("/feedroom/pin", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.PinFeedroom)
	router.PUT("/feedroom/auto_join_members", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.AutoJoinMembers)
	router.POST("/feedroom/participants", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.AddParticipants)
	router.GET("/feedroom/participants", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetParticipants)
	router.DELETE("/feedroom/participants", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.RemoveParticipants)
	router.GET("/feedroom/cohort/access", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetCohortAccess)
	router.PUT("/feedroom/cohort/access", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.EditCohortAccess)
	router.GET("/feedroom/mine", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.MyFeedrooms)
	router.PUT("/feedroom/follow", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.FeedroomFollow)
	router.GET("/feedroom/tag", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetTaggingList)
	router.GET("/feedroom/:feedroom_id/tag", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetTaggingList)

	// Channel Apis
	router.GET("/channel", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), channel.FetchChannel)
	router.GET("/channel/invites", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), channel.GetChannelInvites)
	router.PUT("/channel/invite", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), channel.UpdateChannelInvite)
	router.GET("/channel/:channel_id/settings/member/:participant_uuid", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), channel.GetUserChannelSettings)
	router.PUT("/channel/:channel_id/settings/member/:participant_uuid", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), channel.UpdateUserChannelSettings)

	// Search Apis
	router.GET("/search/channel", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), search.ChannelSearch)
	router.GET("/search/message", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), search.MessageSearch)
	router.GET("/search/post", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), search.PostSearch)
	router.GET("/search/post/user/:user_id", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), search.UserCreatedPostSearch)
	router.GET("/search", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), search.GeneralSearch)

	// Widget Apis
	router.POST("/widget", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), widget.CreateWidget)
	router.GET("/widget", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), widget.GetWidget)
	router.PUT("/widget/:widget_id", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), widget.EditWidget)

	// Poll Apis
	router.PUT("/poll/:poll_id", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), poll.AddPollOption)
	router.PUT("/poll/:poll_id/vote", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), poll.CreatePollVote)
	router.GET("/poll/:poll_id/vote", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), poll.GetPollVotes)

	// Webhook Apis
	router.POST("/webhook", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), webhook.CreateWebhook)
	router.GET("/webhook", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), webhook.GetWebhooks)
	router.GET("/webhook/:webhook_id", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), webhook.GetWebhook)
	router.PATCH("/webhook/:webhook_id", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), webhook.EditWebhook)
	router.DELETE("/webhook/:webhook_id", middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), webhook.DeleteWebhook)

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
