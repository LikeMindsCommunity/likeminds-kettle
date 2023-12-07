package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/nateshr/likeminds-authentication/poll"
	"github.com/nateshr/likeminds-authentication/utility/logger"
	"github.com/nateshr/likeminds-authentication/utility/monitoring"
	"github.com/nateshr/likeminds-authentication/webhook"

	log "github.com/nateshr/likeminds-authentication/logging"
	"github.com/nateshr/likeminds-authentication/widget"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/api_client"
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
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utility"
	"github.com/nateshr/likeminds-authentication/utils"
	"github.com/nateshr/likeminds-authentication/web"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	redisClient *redis.Client
	router      *gin.Engine
)

func main() {
	var AppVersion string = "2.13.0"

	initGin()
	redisClient = cache.InitRedis()
	router.Use(cors.New(enableCors()))
	router.Use(ApiMiddleware(redisClient))
	router.Use(LoggingMiddleware())
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
	router.POST("/user/login", OTMValidationMiddleware(), user.Login)
	router.POST("/user/refresh", RTMValidationMiddleware(), user.Refresh)
	router.POST("/user/logout", LogoutValidationMiddleware(), user.Logout)
	router.POST("/user/merge_account", LTMValidationMiddleware(), user.MergeAccount)
	router.GET("/user/config", LTMValidationMiddleware(), APIKeyValidationMiddleware(), user.Config)
	router.GET("/user/bot", LTMValidationMiddleware(), APIKeyValidationMiddleware(), user.GetBot)
	router.POST("/user/device/push", LTMValidationMiddleware(), APIKeyValidationMiddleware(), user.PushUserToken)
	router.POST("/user/subscription/whatsapp", user.WASubscription)
	router.GET("/user/meta", LTMValidationMiddleware(), APIKeyValidationMiddleware(), user.UserMeta)
	router.POST("/user/otp", OTMValidationMiddleware(), APIKeyValidationMiddleware(), user.GenerateUserOTP)
	router.GET("/user/otp/verify", OTMValidationMiddleware(), APIKeyValidationMiddleware(), user.VerifyUserOTP)
	router.GET("/user/social/login", OTMValidationMiddleware(), APIKeyValidationMiddleware(), user.UserSocialLogin)

	// Home Apis
	router.POST("/home/fetch_communities", LTMValidationMiddleware(), home.FetchCommunities)
	router.GET("/home/dm/meta", LTMValidationMiddleware(), APIKeyValidationMiddleware(), home.DMHome)

	// SDK Apis
	router.POST("/sdk/initiate", VTMValidationMiddleware(false), APIKeyValidationMiddleware(), sdk.InitiateSDK)
	router.POST("/sdk/project", LTMValidationMiddleware(), sdk.CreateProject)
	router.GET("/sdk/project", LTMValidationMiddleware(), sdk.GetProject)
	router.PUT("/sdk/project", LTMValidationMiddleware(), APIKeyValidationMiddleware(), sdk.EditProject)
	router.DELETE("/sdk/project", LTMValidationMiddleware(), APIKeyValidationMiddleware(), sdk.DeleteProject)
	router.GET("/sdk/onboarding", OTMValidationMiddleware(), APIKeyValidationMiddleware(), sdk.GetScreen)
	router.POST("/sdk/onboarding", LTMValidationMiddleware(), APIKeyValidationMiddleware(), sdk.CreateScreen)
	router.PUT("/sdk/onboarding", LTMValidationMiddleware(), APIKeyValidationMiddleware(), sdk.EditScreen)
	router.DELETE("/sdk/onboarding", LTMValidationMiddleware(), APIKeyValidationMiddleware(), sdk.DeleteScreen)

	// Chatroom Apis
	router.GET("/chatroom", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.GetChatroom)
	router.POST("/chatroom", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.CreateChatroom)
	router.PUT("/chatroom", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.EditChatroom)
	router.DELETE("/chatroom", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.DeleteChatroom)
	router.GET("/chatroom/type", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.GetChatroomTypeStatus)
	router.PUT("/chatroom/type", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.ChangeChatroomType)
	router.POST("/chatroom/schedule_follow", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.ScheduleFollow)
	router.PUT("/chatroom/pin", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.PinChatroom)
	router.GET("/chatroom/tag", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.GetTaggingList)
	router.GET("/chatroom/:chatroom_id/tag", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.GetTaggingList)
	router.GET("/chatroom/participants", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.GetParticipants)
	router.POST("/chatroom/participants", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.AddParticipants)
	router.DELETE("/chatroom/participants", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.RemoveParticipants)
	router.GET("/chatroom/settings", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.GetChatroomSettings)
	router.PUT("/chatroom/settings", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.EditChatroomSettings)
	router.PUT("/chatroom/enable_member_message", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.EnableMemberMessage)
	router.PUT("/chatroom/auto_follow_members", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.AutoFollowMembers)
	router.PUT("/chatroom/files", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.UpdateFiles)
	router.GET("/chatroom/mine", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.MyChatrooms)
	router.PUT("/chatroom/seen", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.CollabcardSeen)
	router.PUT("/chatroom/follow", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.ChatroomFollow)
	router.PUT("/chatroom/mute", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.MuteChatroom)
	router.PUT("/chatroom/rename", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.RenameChatroom)
	router.GET("/chatroom/share", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.FetchShareUrl)
	router.GET("/chatroom/pending", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.FetchPendingChatroom)
	router.PUT("/chatroom/pending", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.ActionPendingChatroom)
	router.GET("/chatroom/sync", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.SyncChatrooms)
	router.POST("/chatroom/dm/block", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.ChatroomBlock)
	router.POST("/chatroom/dm/request", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.InitiatingDMRequest)
	router.POST("/chatroom/dm/create", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.CreateDM)
	router.GET("/chatroom/dm", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.ListDMChatrooms)
	router.GET("/chatroom/dm/limit", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.DMLimit)
	router.GET("/chatroom/search", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.ChatroomSearch)
	router.POST("/chatroom/cohort", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.AddCohortToChatroom)
	router.DELETE("/chatroom/cohort", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.RemoveCohortFromChatroom)
	router.GET("/chatroom/cohort/access", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.GetCohortAccess)
	router.PUT("/chatroom/cohort/access", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.EditCohortAccess)
	router.GET("/chatroom/home", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.GetChatroomHome)
	router.POST("/chatroom/mark_read", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.ChatroomMarkRead)
	router.GET("/chatroom/event", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.FetchEvents)
	router.POST("/chatroom/event", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.CreateEvent)
	router.PUT("/chatroom/event", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.EditEvent)
	router.GET("/chatroom/event/meta", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.FetchEventMeta)
	router.GET("/chatroom/event/link", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.FetchEventLinks)
	router.GET("/chatroom/event/unseen_count", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.FetchEventUnseenCount)
	router.POST("/chatroom/event/recordings", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.UploadEventRecordings)
	router.DELETE("/chatroom/event/recordings", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.DeleteEventRecordings)
	router.POST("/chatroom/event/recordings/meta", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.UploadEventRecordingsMeta)
	router.DELETE("/chatroom/event/recordings/meta", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.DeleteEventRecordingsMeta)
	router.POST("/chatroom/event/instructors", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.AddEventInstructors)
	router.POST("/chatroom/event/highlights", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.AddEventHighlights)
	router.POST("/chatroom/event/testimonials", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.AddEventTestimonials)
	router.POST("/chatroom/event/faq", LTMValidationMiddleware(), APIKeyValidationMiddleware(), chatroom.AddEventFAQ)

	// Community Apis
	router.GET("/community", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.Community)
	router.GET("/community/branding", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.Branding)
	router.POST("/community/questions", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.EditQuestions)
	router.GET("/community/questions", LTMorVTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetQuestions)
	router.GET("/community/question/filters", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetCommunityQuestionFilters)
	router.GET("/community/member", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetMember)
	router.POST("/community/member", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.AddMember)
	router.DELETE("/community/member", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.RemoveMembers)
	router.DELETE("/community/member/leave", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.LeaveCommunity)
	router.PUT("/community/member", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.EditMember)
	router.GET("/community/member/state", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.FetchMemberState)
	router.GET("/community/member/role", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.FetchMemberRole)
	router.DELETE("/community/manager/remove", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.RemoveCommunityManager)
	router.DELETE("/community/admin", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.RemoveCommunityManager)
	router.DELETE("/community/member/remove", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.RemoveMember)
	router.GET("/community/management/tool", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetManagementTools)
	router.GET("/community/report", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetReport)
	router.POST("/community/report", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.PushReport)
	router.DELETE("/community/report", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.CloseReport)
	router.GET("/community/report/tag", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetReportTags)
	router.GET("/community/settings", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetCommunitySettings)
	router.PUT("/community/settings", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.UpdateCommunitySettings)
	router.GET("/community/rights", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetCommunityRights)
	router.PUT("/community/rights", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.EditCommunityRights)
	router.PATCH("/community/rights", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.UpdateCommunityRights)
	router.GET("/community/settings/dm", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetCommunityDMSettings)
	router.PUT("/community/settings/dm", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.EditCommunityDMSettings)
	router.GET("/community/feed/dm", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.DMFeed)
	router.GET("/community/dm/status", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.DMStatus)
	router.GET("/community/member/search", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.MemberSearch)
	router.GET("/community/member/profile", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetMemberProfile)
	router.PUT("/community/member/profile", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.EditMemberProfile)
	router.GET("/community/member/chatroom", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.MemberChatroom)
	router.GET("/community/member/:user_id/channel", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.CommunityMemberChannels)
	router.GET("/community/member/channel/status", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetMemberChannels)
	router.POST("/community/cohort", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.CreateCohort)
	router.GET("/community/cohort", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetCohort)
	router.GET("/community/cohort/:cohort_id", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.FetchCohort)
	router.DELETE("/community/cohort", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.DeleteCohort)
	router.PUT("/community/cohort", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.EditCohort)
	router.DELETE("/community/cohort/member", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.RemoveCohortMember)
	router.GET("/community/feed", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetCommunityFeed)
	router.GET("/community/settings/notification/conversation", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetConversationNotificationSettings)
	router.PUT("/community/settings/notification/conversation", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.EditConversationNotificationSettings)
	router.GET("/community/settings/notification/feed", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetFeedNotificationSettings)
	router.PUT("/community/settings/notification/feed", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.EditFeedNotificationSettings)
	router.GET("/community/settings/notification", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetNotificationSettings)
	router.PUT("/community/settings/notification", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.EditNotificationSettings)
	router.GET("/community/tag", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetTaggingList)
	router.GET("/community/settings/content_download", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetContentDownloadSettings)
	router.PUT("/community/settings/content_download", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.EditContentDownloadSettings)
	router.GET("/community/member/home/meta", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.MemberHomeMeta)
	router.PUT("/community/member/join", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.AcceptRejectJoinCommunity)
	router.GET("/community/intro_examples", LTMorVTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetIntroExamples)
	router.POST("/community/invite", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.SendCommunityInvite)
	router.GET("/community/configurations", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetCommunityConfigurations)
	router.GET("/community/member/pending", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetPendingCommunityMembers)
	router.GET("/community/removal_reports", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetRemovalReports)
	router.POST("/community/member/:user_id/connection", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.CreateMemberConnection)
	router.PATCH("/community/member/:user_id/connection", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.AcceptRejectMemberConnection)
	router.GET("/community/member/:user_id/connection", LTMValidationMiddleware(), APIKeyValidationMiddleware(), community.GetMemberConnection)

	// Moderation Apis
	router.GET("/moderation/rights", LTMValidationMiddleware(), APIKeyValidationMiddleware(), moderation.GetRights)
	router.PUT("/moderation/rights", LTMValidationMiddleware(), APIKeyValidationMiddleware(), moderation.EditRights)
	router.PATCH("/moderation/rights", LTMValidationMiddleware(), APIKeyValidationMiddleware(), moderation.UpdateRights)

	// Conversation Apis
	router.GET("/conversation", LTMValidationMiddleware(), APIKeyValidationMiddleware(), conversation.GetConversation)
	router.POST("/conversation", LTMValidationMiddleware(), APIKeyValidationMiddleware(), conversation.CreateConversation)
	router.PUT("/conversation", LTMValidationMiddleware(), APIKeyValidationMiddleware(), conversation.EditConversation)
	router.DELETE("/conversation", LTMValidationMiddleware(), APIKeyValidationMiddleware(), conversation.DeleteConversation)
	router.PUT("/conversation/reaction", LTMValidationMiddleware(), APIKeyValidationMiddleware(), conversation.AddReaction)
	router.DELETE("/conversation/reaction", LTMValidationMiddleware(), APIKeyValidationMiddleware(), conversation.RemoveReaction)
	router.POST("/conversation/poll", LTMValidationMiddleware(), APIKeyValidationMiddleware(), conversation.AddPoll)
	router.POST("/conversation/poll/submit", LTMValidationMiddleware(), APIKeyValidationMiddleware(), conversation.SubmitPoll)
	router.GET("/conversation/poll/users", LTMValidationMiddleware(), APIKeyValidationMiddleware(), conversation.PollUsers)
	router.PUT("/conversation/topic", LTMValidationMiddleware(), APIKeyValidationMiddleware(), conversation.SetTopic)
	router.PUT("/conversation/event/attend", LTMValidationMiddleware(), APIKeyValidationMiddleware(), conversation.EventAttend)
	router.PUT("/conversation/event/attended", LTMValidationMiddleware(), APIKeyValidationMiddleware(), conversation.EventAttended)
	router.GET("/conversation/event/unseen", LTMValidationMiddleware(), APIKeyValidationMiddleware(), conversation.FetchEventUnseenCount)
	router.PUT("/conversation/event/last_seen", LTMValidationMiddleware(), APIKeyValidationMiddleware(), conversation.UpdateLastSeenEvent)
	router.GET("/conversation/event/link", LTMValidationMiddleware(), APIKeyValidationMiddleware(), conversation.FetchEventLink)
	router.GET("/conversation/event", LTMValidationMiddleware(), APIKeyValidationMiddleware(), conversation.FetchAllEvents)
	router.GET("/conversation/preview/unread", LTMValidationMiddleware(), APIKeyValidationMiddleware(), conversation.FetchUnreadPreviews)
	router.GET("/conversation/preview/unread_count", LTMValidationMiddleware(), APIKeyValidationMiddleware(), conversation.FetchPreviewUnreadMessagesCount)
	router.GET("/conversation/notification/unread", LTMValidationMiddleware(), APIKeyValidationMiddleware(), conversation.UnreadConversationNotification)
	router.GET("/conversation/sync", LTMValidationMiddleware(), APIKeyValidationMiddleware(), conversation.SyncConversation)
	router.GET("/conversation/search", LTMValidationMiddleware(), APIKeyValidationMiddleware(), conversation.ConversationSearch)

	// Feed Apis
	router.POST("/feed/post", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.CreatePost)
	router.GET("/feed/post/:post_id", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.GetPost)
	router.PUT("/feed/post/:post_id", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.EditPost)
	router.DELETE("/feed/post/:post_id", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.DeletePost)
	router.PUT("/feed/post/:post_id/like", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.CreatePostLike)
	router.GET("/feed/post/:post_id/like", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.GetPostLikes)
	router.PUT("/feed/post/:post_id/pin", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.PinPost)
	router.PUT("/feed/post/:post_id/save", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.CreateSavePost)
	router.POST("/feed/post/:post_id/comment", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.CommentPost)
	router.PUT("/feed/post/:post_id/comment/:comment_id", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.EditCommentPost)
	router.POST("/feed/post/:post_id/comment/:comment_id/comment", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.CreateCommentReply)
	router.GET("/feed/post/:post_id/comment/:comment_id", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.GetComment)
	router.DELETE("/feed/post/:post_id/comment/:comment_id", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.DeleteComment)
	router.PUT("/feed/post/:post_id/comment/:comment_id/like", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.CreateCommentLike)
	router.GET("/feed/post/:post_id/comment/:comment_id/like", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.GetCommentLikes)
	router.GET("/feed/user/:user_id/save", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.GetSavedPosts)
	router.GET("/feed/user/:user_id/post", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.FetchUserCreatedPosts)
	router.GET("/feed/user/activity", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.GetUserActivity)
	router.POST("/feed/user/:user_id/activity", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.CreateUserActivity)
	router.GET("/feed/user/activity/unread_count", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.GetUserActivityUnreadCount)
	router.POST("/feed/user/activity/:activity_id/mark_read", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.UserActivityMarkRead)
	router.GET("/feed/universal", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.FetchUniversalFeed)
	router.GET("/feed/group", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.FetchGroupFeed)
	router.POST("/feed/topic", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.CreateTopic)
	router.GET("/feed/topic", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.GetTopic)
	router.PUT("/feed/topic/:topic_id", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.EditTopic)
	router.GET("feed/:user_id/follow", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feed.GetFollowingFeed)

	// Utility Apis
	router.GET("/helper/url", LTMValidationMiddleware(), APIKeyValidationMiddleware(), utility.DecodeUrl)
	router.POST("/helper/media/upload", LTMValidationMiddleware(), APIKeyValidationMiddleware(), utility.UploadFiles)

	// Feedroom Apis
	router.POST("/feedroom", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feedroom.CreateFeedroom)
	router.PUT("/feedroom", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feedroom.EditFeedroom)
	router.DELETE("/feedroom", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feedroom.DeleteFeedroom)
	router.GET("/feedroom", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feedroom.GetFeedroom)
	router.GET("/feedroom/action", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feedroom.GetFeedroomMenu)
	router.GET("/feedroom/settings", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feedroom.GetFeedroomSettings)
	router.PUT("/feedroom/type", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feedroom.ChangeFeedroomType)
	router.GET("/feedroom/type", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feedroom.GetFeedroomTypeStatus)
	router.PUT("/feedroom/enable_member_post", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feedroom.EnableMemberPost)
	router.PUT("/feedroom/pin", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feedroom.PinFeedroom)
	router.PUT("/feedroom/auto_join_members", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feedroom.AutoJoinMembers)
	router.POST("/feedroom/participants", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feedroom.AddParticipants)
	router.GET("/feedroom/participants", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feedroom.GetParticipants)
	router.DELETE("/feedroom/participants", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feedroom.RemoveParticipants)
	router.GET("/feedroom/cohort/access", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feedroom.GetCohortAccess)
	router.PUT("/feedroom/cohort/access", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feedroom.EditCohortAccess)
	router.GET("/feedroom/mine", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feedroom.MyFeedrooms)
	router.PUT("/feedroom/follow", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feedroom.FeedroomFollow)
	router.GET("/feedroom/tag", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feedroom.GetTaggingList)
	router.GET("/feedroom/:feedroom_id/tag", LTMValidationMiddleware(), APIKeyValidationMiddleware(), feedroom.GetTaggingList)

	// Channel Apis
	router.GET("/channel", LTMValidationMiddleware(), APIKeyValidationMiddleware(), channel.FetchChannel)
	router.GET("/channel/invites", LTMValidationMiddleware(), APIKeyValidationMiddleware(), channel.GetChannelInvites)
	router.PUT("/channel/invite", LTMValidationMiddleware(), APIKeyValidationMiddleware(), channel.UpdateChannelInvite)
	router.GET("/channel/:channel_id/settings/member/:participant_uuid", LTMValidationMiddleware(), APIKeyValidationMiddleware(), channel.GetUserChannelSettings)
	router.PUT("/channel/:channel_id/settings/member/:participant_uuid", LTMValidationMiddleware(), APIKeyValidationMiddleware(), channel.UpdateUserChannelSettings)

	// Search Apis
	router.GET("/search/channel", LTMValidationMiddleware(), APIKeyValidationMiddleware(), search.ChannelSearch)
	router.GET("/search/message", LTMValidationMiddleware(), APIKeyValidationMiddleware(), search.MessageSearch)
	router.GET("/search/post", LTMValidationMiddleware(), APIKeyValidationMiddleware(), search.PostSearch)
	router.GET("/search/post/user/:user_id", LTMValidationMiddleware(), APIKeyValidationMiddleware(), search.UserCreatedPostSearch)
	router.GET("/search", LTMValidationMiddleware(), APIKeyValidationMiddleware(), search.GeneralSearch)

	// Widget Apis
	router.POST("/widget", LTMValidationMiddleware(), APIKeyValidationMiddleware(), widget.CreateWidget)
	router.GET("/widget", LTMValidationMiddleware(), APIKeyValidationMiddleware(), widget.GetWidget)
	router.PUT("/widget/:widget_id", LTMValidationMiddleware(), APIKeyValidationMiddleware(), widget.EditWidget)

	// Poll Apis
	router.PUT("/poll/:poll_id", LTMValidationMiddleware(), APIKeyValidationMiddleware(), poll.AddPollOption)
	router.PUT("/poll/:poll_id/vote", LTMValidationMiddleware(), APIKeyValidationMiddleware(), poll.CreatePollVote)
	router.GET("/poll/:poll_id/vote", LTMValidationMiddleware(), APIKeyValidationMiddleware(), poll.GetPollVotes)

	// Webhook Apis
	router.POST("/webhook", LTMValidationMiddleware(), APIKeyValidationMiddleware(), webhook.CreateWebhook)
	router.GET("/webhook", LTMValidationMiddleware(), APIKeyValidationMiddleware(), webhook.GetWebhooks)
	router.GET("/webhook/:webhook_id", LTMValidationMiddleware(), APIKeyValidationMiddleware(), webhook.GetWebhook)
	router.PATCH("/webhook/:webhook_id", LTMValidationMiddleware(), APIKeyValidationMiddleware(), webhook.EditWebhook)
	router.DELETE("/webhook/:webhook_id", LTMValidationMiddleware(), APIKeyValidationMiddleware(), webhook.DeleteWebhook)

	// Logging Apis
	router.POST("/logs", LTMValidationMiddleware(), APIKeyValidationMiddleware(), logger.PushLogs)

	log.Info(fmt.Sprintf("application version: %s", AppVersion))
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	log.Fatal(router.Run(":8080"))
}

func initGin() {
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
}

// Method to process API request to log
func processRequest(c *gin.Context) interface{} {
	requestBodyData := gin.H{}

	// Reading request body
	requestBody, err := ioutil.ReadAll(c.Request.Body)

	// Updating request body after read
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	// Unmarshalling request body
	if err == nil {
		_ = json.Unmarshal(requestBody, &requestBodyData)
	}

	return gin.H{
		"host":         c.Request.Host,
		"absolute_uri": c.Request.RequestURI,
		"method":       c.Request.Method,
		"headers":      c.Request.Header,
		"body":         requestBodyData,
	}
}

// responseBodyWriter | Custom Response Writer
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write | Custom Write Method for responseBodyWriter
func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// LoggingMiddleware will log the request and response of API
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.RequestURI == "/" {

			c.Next()

		} else {

			data := gin.H{}

			// Starting time
			startTime := time.Now()

			// Implementing custom response body writer in the context
			w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
			c.Writer = w

			// Updating Request Data
			data["request"] = processRequest(c)

			// Processing request
			c.Next()

			// End Time
			endTime := time.Now()

			response := gin.H{}
			statusCode := c.Writer.Status()

			// Unmarshalling Request Response
			_ = json.Unmarshal(w.body.Bytes(), &response)

			// Updating Request Response
			data["response"] = gin.H{
				"http_response_code": statusCode,
				"content":            response,
			}

			if statusCode < http.StatusBadRequest {
				data["response"].(gin.H)["content"] = gin.H{}
			}

			// Updating Request Meta Data
			data["meta"] = gin.H{
				"latency":   endTime.Sub(startTime),
				"client_ip": c.ClientIP(),
			}

			// Marshalling the final Data
			marshelledData, _ := json.Marshal(data)

			if statusCode >= http.StatusOK && statusCode < http.StatusBadRequest {
				// Logging the generated request data as Info
				log.Info(string(marshelledData))
			} else {
				// Logging the generated request data as Error
				log.Error(string(marshelledData))
			}

			c.Next()
		}
	}
}

// ApiMiddleware will add the db connection to the context
func ApiMiddleware(client *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(cache.ParamRedisClient, client)
		c.Next()
	}
}

func OTMValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract OTM from token, internally it checks if token is valid or not
		otm, err := token.ExtractOTM(c.Request.Header.Get(token.HeaderAuthorization))

		if otm == nil {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidOTM,
			})
			return

		} else {
			// If valid, set "otm" in context, to be used in later APIs
			c.Set(token.ParamOTM, otm)
		}
		c.Next()
	}
}

func VTMValidationMiddleware(isMandatory bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract VTM from token, internally it checks if token is valid or not
		vtm, err := token.ExtractVTM(c.Request.Header.Get(token.HeaderAuthorization))

		if vtm == nil && isMandatory {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidVTM,
			})
			return

		} else if vtm == nil {
			log.Error(err)
			c.Next()

		} else {
			// If valid, set "vtm" in context, to be used in later APIs
			c.Set(token.ParamVTM, vtm)
		}
		c.Next()
	}
}

func LTMValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get(token.HeaderAuthorization)
		//Extract LTM from token, internally it checks if token is valid or not
		ltm, err := token.ExtractLTM(bearerToken)
		if ltm == nil {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidLTM,
			})
			return
		} else {
			//Check if LTM is black listed or not
			if cache.IsLTMBlacklisted(redisClient, ltm) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
					Success:      false,
					ErrorMessage: utils.ErrorDeviceLoggedOut,
				})
				return
			}
			//If valid and not blacklisted, set "ltm" in context, to be used in later APIs
			c.Set(token.ParamLTM, ltm)
		}
		c.Next()
	}
}

func RTMValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Extract RTM from token, internally it checks if token is valid or not
		rtm, err := token.ExtractRTM(c.Request.Header.Get(token.HeaderAuthorization))
		if rtm == nil {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidRTM,
			})
			return
		} else {
			//Check if RTM is black listed or not
			if cache.IsRTMBlacklisted(redisClient, rtm) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
					Success:      false,
					ErrorMessage: utils.ErrorDeviceLoggedOut,
				})
				return
			}
			//If valid and not blacklisted, set "rtm" in context, to be used in later APIs
			c.Set(token.ParamRTM, rtm)
		}
		c.Next()
	}
}

func LTMorVTMValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from headers
		bearerToken := c.Request.Header.Get(token.HeaderAuthorization)

		if bearerToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidLTMorVTM,
			})
			return
		}

		// Extract LTM info from token, internally it checks if token is valid or not
		ltm, ltmErr := token.ExtractLTM(bearerToken)

		if ltmErr == nil {
			c.Set(token.ParamLTM, ltm)
			c.Next()
		}

		// Extract VTM info from token, internally it checks if token is valid or not
		vtm, vtmErr := token.ExtractVTM(bearerToken)

		if vtmErr == nil {
			c.Set(token.ParamVTM, vtm)
			c.Next()

		} else {
			log.Error(ltmErr)
			log.Error(vtmErr)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidLTMorVTM,
			})
			return
		}
	}
}

func LogoutValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Extract LTM from token, internally it checks if token is valid or not
		ltm, err := token.ExtractLTM(c.Request.Header.Get(token.HeaderAuthorization))
		if ltm == nil {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidLTM,
			})
			return
		} else {
			//Check if LTM is black listed or not
			if cache.IsLTMBlacklisted(redisClient, ltm) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
					Success:      false,
					ErrorMessage: utils.ErrorDeviceLoggedOut,
				})
				return
			}
			//Get RTM token from body
			var logoutRequest user.LogoutRequest
			if err := c.ShouldBindJSON(&logoutRequest); err != nil {
				c.JSON(http.StatusUnprocessableEntity, utils.Response{
					Success:      false,
					ErrorMessage: utils.ErrorInvalidRequest,
				})
				return
			}
			//Extract RTM from token, internally it checks if token is valid or not
			rtm, err := token.ExtractRTM(logoutRequest.RefreshToken)
			if rtm == nil {
				log.Error(err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
					Success:      false,
					ErrorMessage: token.ErrorInvalidRTM,
				})
				return
			} else {
				//Check if RTM is black listed or not
				if cache.IsRTMBlacklisted(redisClient, rtm) {
					c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
						Success:      false,
						ErrorMessage: utils.ErrorDeviceLoggedOut,
					})
					return
				}
				//If valid and not blacklisted, set "ltm" and "rtm" in context, to be used in later APIs
				c.Set(token.ParamLTM, ltm)
				c.Set(token.ParamRTM, rtm)
			}
		}
		c.Next()
	}
}

const SDKAuthenticateEndPoint = "/api/sdk/authenticate"
const ResponseCommunityId = "community_id"

func APIKeyValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request has LTM token or not
		ltm, ok := c.Get(token.ParamLTM)
		if ok && ltm.(*token.LoginTokenMeta).ApiKey != "" {
			c.Request.Header["X-Api-Key"] = []string{ltm.(*token.LoginTokenMeta).ApiKey}
			c.Next()
		}

		// Check if request has OTM token or not
		otm, ok := c.Get(token.ParamOTM)
		if ok && otm.(*token.OnboardingTokenMeta).ApiKey != "" {
			c.Request.Header["X-Api-Key"] = []string{otm.(*token.OnboardingTokenMeta).ApiKey}
			c.Next()
		}

		// Check if request has VTM token or not
		vtm, ok := c.Get(token.ParamVTM)
		if ok && vtm.(*token.VerifyTokenMeta).ApiKey != "" {
			c.Request.Header["X-Api-Key"] = []string{vtm.(*token.VerifyTokenMeta).ApiKey}
			c.Next()
		}

		//Create internal API client
		client := api_client.NewAPIClient()
		options := api_client.GetRequestOptions{
			Url:           client.CoreServiceBaseURL + SDKAuthenticateEndPoint,
			CustomHeaders: utils.CreateHeaders(c, ""),
		}
		//Send request
		respBytes, _, err := client.GetRequest(&options)
		if err != nil {
			//If API fails or any other error
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: err.Error(),
			})
			return
		}
		//Parse response
		var apiCR api_client.APIClientResponse
		err = api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)
		if err != nil {
			//Internal unmarshal error
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: err.Error(),
			})
			return
		}

		if !apiCR.Success {
			//If api/sdk/authenticate returns success as false
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: utils.ErrorInvalidAPIKey,
			})
			return
		}
		c.Set(ResponseCommunityId, apiCR.Response[ResponseCommunityId])
		c.Next()
	}
}

// GuestAccessCheckMiddleware | restrict guest access on endpoints
func GuestAccessCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Create internal API client
		client := api_client.NewAPIClient()
		options := api_client.GetRequestOptions{
			Url:           client.CoreServiceBaseURL + user.UserFetchEndPoint,
			CustomHeaders: utils.CreateHeaders(c, c.GetHeader(utils.HeadersMemberId)),
		}
		//Send request
		respBytes, statusCode, err := client.GetRequest(&options)
		if err != nil {
			//If API fails or any other error
			utils.GeneralAPIError(c, err.Error())
			return
		}
		//Parse response
		var apiCR api_client.APIClientResponse
		err = api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)
		if err != nil {
			//Internal unmarshal error
			utils.GeneralAPIError(c, err.Error())
			return
		}

		if !apiCR.Success {
			//If api/user/fetch returns success as false
			c.JSON(statusCode, apiCR)
			return
		}

		isGuest := apiCR.Response[user.ResponseUser].(map[string]interface{})[user.ResponseUserIsGuest].(bool)
		if isGuest {
			type GuestAccessDeniedResponseData struct {
				Route string `json:"route"`
			}
			response := utils.Response{
				Success:      false,
				ErrorMessage: utils.ErrorGuestAccessDenied,
				Data:         GuestAccessDeniedResponseData{Route: user.GuestLoginRoute},
			}

			// If user is guest returns success as false
			utils.APIError(c, http.StatusUnauthorized, response)
			return
		}
		c.Next()
	}
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
		log.Fatal(err.Error())
		return nil
	}
	return prometheusService
}
