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
	var AppVersion string = "2.20.0"

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
	router.GET("/user/config", LTMValidationMiddleware(), user.Config)
	router.GET("/user/bot", LTMValidationMiddleware(), user.GetBot)
	router.POST("/user/device/push", LTMValidationMiddleware(), user.PushUserToken)
	router.POST("/user/subscription/whatsapp", user.WASubscription)
	router.GET("/user/meta", LTMValidationMiddleware(), user.UserMeta)
	router.POST("/user/otp", OTMValidationMiddleware(), user.GenerateUserOTP)
	router.GET("/user/otp/verify", OTMValidationMiddleware(), user.VerifyUserOTP)
	router.GET("/user/social/login", OTMValidationMiddleware(), user.UserSocialLogin)

	// Home Apis
	router.POST("/home/fetch_communities", LTMValidationMiddleware(), home.FetchCommunities)
	router.GET("/home/dm/meta", LTMValidationMiddleware(), home.DMHome)

	// SDK Apis
	router.POST("/sdk/initiate", VTMValidationMiddleware(false), sdk.InitiateSDK)
	router.POST("/sdk/project", LTMValidationMiddleware(), sdk.CreateProject)
	router.GET("/sdk/project", LTMValidationMiddleware(), sdk.GetProject)
	router.PUT("/sdk/project", LTMValidationMiddleware(), sdk.EditProject)
	router.DELETE("/sdk/project", LTMValidationMiddleware(), sdk.DeleteProject)
	router.GET("/sdk/onboarding", OTMValidationMiddleware(), sdk.GetScreen)
	router.POST("/sdk/onboarding", LTMValidationMiddleware(), sdk.CreateScreen)
	router.PUT("/sdk/onboarding", LTMValidationMiddleware(), sdk.EditScreen)
	router.DELETE("/sdk/onboarding", LTMValidationMiddleware(), sdk.DeleteScreen)

	// Chatroom Apis
	router.GET("/chatroom", LTMValidationMiddleware(), chatroom.GetChatroom)
	router.POST("/chatroom", LTMValidationMiddleware(), chatroom.CreateChatroom)
	router.PUT("/chatroom", LTMValidationMiddleware(), chatroom.EditChatroom)
	router.DELETE("/chatroom", LTMValidationMiddleware(), chatroom.DeleteChatroom)
	router.GET("/chatroom/type", LTMValidationMiddleware(), chatroom.GetChatroomTypeStatus)
	router.PUT("/chatroom/type", LTMValidationMiddleware(), chatroom.ChangeChatroomType)
	router.POST("/chatroom/schedule_follow", LTMValidationMiddleware(), chatroom.ScheduleFollow)
	router.PUT("/chatroom/pin", LTMValidationMiddleware(), chatroom.PinChatroom)
	router.GET("/chatroom/tag", LTMValidationMiddleware(), chatroom.GetTaggingList)
	router.GET("/chatroom/:chatroom_id/tag", LTMValidationMiddleware(), chatroom.GetTaggingList)
	router.GET("/chatroom/participants", LTMValidationMiddleware(), chatroom.GetParticipants)
	router.POST("/chatroom/participants", LTMValidationMiddleware(), chatroom.AddParticipants)
	router.DELETE("/chatroom/participants", LTMValidationMiddleware(), chatroom.RemoveParticipants)
	router.GET("/chatroom/settings", LTMValidationMiddleware(), chatroom.GetChatroomSettings)
	router.PUT("/chatroom/settings", LTMValidationMiddleware(), chatroom.EditChatroomSettings)
	router.PUT("/chatroom/enable_member_message", LTMValidationMiddleware(), chatroom.EnableMemberMessage)
	router.PUT("/chatroom/auto_follow_members", LTMValidationMiddleware(), chatroom.AutoFollowMembers)
	router.PUT("/chatroom/files", LTMValidationMiddleware(), chatroom.UpdateFiles)
	router.GET("/chatroom/mine", LTMValidationMiddleware(), chatroom.MyChatrooms)
	router.PUT("/chatroom/seen", LTMValidationMiddleware(), chatroom.CollabcardSeen)
	router.PUT("/chatroom/follow", LTMValidationMiddleware(), chatroom.ChatroomFollow)
	router.PUT("/chatroom/mute", LTMValidationMiddleware(), chatroom.MuteChatroom)
	router.PUT("/chatroom/rename", LTMValidationMiddleware(), chatroom.RenameChatroom)
	router.GET("/chatroom/share", LTMValidationMiddleware(), chatroom.FetchShareUrl)
	router.GET("/chatroom/pending", LTMValidationMiddleware(), chatroom.FetchPendingChatroom)
	router.PUT("/chatroom/pending", LTMValidationMiddleware(), chatroom.ActionPendingChatroom)
	router.GET("/chatroom/sync", LTMValidationMiddleware(), chatroom.SyncChatrooms)
	router.POST("/chatroom/dm/block", LTMValidationMiddleware(), chatroom.ChatroomBlock)
	router.POST("/chatroom/dm/request", LTMValidationMiddleware(), chatroom.InitiatingDMRequest)
	router.POST("/chatroom/dm/create", LTMValidationMiddleware(), chatroom.CreateDM)
	router.GET("/chatroom/dm", LTMValidationMiddleware(), chatroom.ListDMChatrooms)
	router.GET("/chatroom/dm/limit", LTMValidationMiddleware(), chatroom.DMLimit)
	router.GET("/chatroom/search", LTMValidationMiddleware(), chatroom.ChatroomSearch)
	router.POST("/chatroom/cohort", LTMValidationMiddleware(), chatroom.AddCohortToChatroom)
	router.DELETE("/chatroom/cohort", LTMValidationMiddleware(), chatroom.RemoveCohortFromChatroom)
	router.GET("/chatroom/cohort/access", LTMValidationMiddleware(), chatroom.GetCohortAccess)
	router.PUT("/chatroom/cohort/access", LTMValidationMiddleware(), chatroom.EditCohortAccess)
	router.GET("/chatroom/home", LTMValidationMiddleware(), chatroom.GetChatroomHome)
	router.POST("/chatroom/mark_read", LTMValidationMiddleware(), chatroom.ChatroomMarkRead)
	router.GET("/chatroom/event", LTMValidationMiddleware(), chatroom.FetchEvents)
	router.POST("/chatroom/event", LTMValidationMiddleware(), chatroom.CreateEvent)
	router.PUT("/chatroom/event", LTMValidationMiddleware(), chatroom.EditEvent)
	router.GET("/chatroom/event/meta", LTMValidationMiddleware(), chatroom.FetchEventMeta)
	router.GET("/chatroom/event/link", LTMValidationMiddleware(), chatroom.FetchEventLinks)
	router.GET("/chatroom/event/unseen_count", LTMValidationMiddleware(), chatroom.FetchEventUnseenCount)
	router.POST("/chatroom/event/recordings", LTMValidationMiddleware(), chatroom.UploadEventRecordings)
	router.DELETE("/chatroom/event/recordings", LTMValidationMiddleware(), chatroom.DeleteEventRecordings)
	router.POST("/chatroom/event/recordings/meta", LTMValidationMiddleware(), chatroom.UploadEventRecordingsMeta)
	router.DELETE("/chatroom/event/recordings/meta", LTMValidationMiddleware(), chatroom.DeleteEventRecordingsMeta)
	router.POST("/chatroom/event/instructors", LTMValidationMiddleware(), chatroom.AddEventInstructors)
	router.POST("/chatroom/event/highlights", LTMValidationMiddleware(), chatroom.AddEventHighlights)
	router.POST("/chatroom/event/testimonials", LTMValidationMiddleware(), chatroom.AddEventTestimonials)
	router.POST("/chatroom/event/faq", LTMValidationMiddleware(), chatroom.AddEventFAQ)

	// Community Apis
	router.GET("/community", LTMValidationMiddleware(), community.Community)
	router.GET("/community/branding", LTMValidationMiddleware(), community.Branding)
	router.POST("/community/questions", LTMValidationMiddleware(), community.EditQuestions)
	router.GET("/community/questions", LTMorVTMValidationMiddleware(), community.GetQuestions)
	router.GET("/community/question/filters", LTMValidationMiddleware(), community.GetCommunityQuestionFilters)
	router.GET("/community/member", LTMValidationMiddleware(), community.GetMember)
	router.POST("/community/member", LTMValidationMiddleware(), community.AddMember)
	router.DELETE("/community/member", LTMValidationMiddleware(), community.RemoveMembers)
	router.DELETE("/community/member/leave", LTMValidationMiddleware(), community.LeaveCommunity)
	router.PUT("/community/member", LTMValidationMiddleware(), community.EditMember)
	router.GET("/community/member/state", LTMValidationMiddleware(), community.FetchMemberState)
	router.GET("/community/member/role", LTMValidationMiddleware(), community.FetchMemberRole)
	router.DELETE("/community/manager/remove", LTMValidationMiddleware(), community.RemoveCommunityManager)
	router.DELETE("/community/admin", LTMValidationMiddleware(), community.RemoveCommunityManager)
	router.DELETE("/community/member/remove", LTMValidationMiddleware(), community.RemoveMember)
	router.GET("/community/management/tool", LTMValidationMiddleware(), community.GetManagementTools)
	router.GET("/community/report", LTMValidationMiddleware(), community.GetReport)
	router.POST("/community/report", LTMValidationMiddleware(), community.PushReport)
	router.DELETE("/community/report", LTMValidationMiddleware(), community.CloseReport)
	router.PATCH("/community/report", LTMValidationMiddleware(), community.UpdateReports)
	router.GET("/community/report/tag", LTMValidationMiddleware(), community.GetReportTags)
	router.GET("/community/settings", LTMValidationMiddleware(), community.GetCommunitySettings)
	router.PUT("/community/settings", LTMValidationMiddleware(), community.UpdateCommunitySettings)
	router.GET("/community/rights", LTMValidationMiddleware(), community.GetCommunityRights)
	router.PUT("/community/rights", LTMValidationMiddleware(), community.EditCommunityRights)
	router.PATCH("/community/rights", LTMValidationMiddleware(), community.UpdateCommunityRights)
	router.GET("/community/settings/dm", LTMValidationMiddleware(), community.GetCommunityDMSettings)
	router.PUT("/community/settings/dm", LTMValidationMiddleware(), community.EditCommunityDMSettings)
	router.GET("/community/feed/dm", LTMValidationMiddleware(), community.DMFeed)
	router.GET("/community/dm/status", LTMValidationMiddleware(), community.DMStatus)
	router.GET("/community/member/search", LTMValidationMiddleware(), community.MemberSearch)
	router.GET("/community/member/profile", LTMValidationMiddleware(), community.GetMemberProfile)
	router.PUT("/community/member/profile", LTMValidationMiddleware(), community.EditMemberProfile)
	router.GET("/community/member/chatroom", LTMValidationMiddleware(), community.MemberChatroom)
	router.GET("/community/member/:user_id/channel", LTMValidationMiddleware(), community.CommunityMemberChannels)
	router.GET("/community/member/channel/status", LTMValidationMiddleware(), community.GetMemberChannels)
	router.POST("/community/cohort", LTMValidationMiddleware(), community.CreateCohort)
	router.GET("/community/cohort", LTMValidationMiddleware(), community.GetCohort)
	router.GET("/community/cohort/:cohort_id", LTMValidationMiddleware(), community.FetchCohort)
	router.DELETE("/community/cohort", LTMValidationMiddleware(), community.DeleteCohort)
	router.PUT("/community/cohort", LTMValidationMiddleware(), community.EditCohort)
	router.DELETE("/community/cohort/member", LTMValidationMiddleware(), community.RemoveCohortMember)
	router.GET("/community/feed", LTMValidationMiddleware(), community.GetCommunityFeed)
	router.GET("/community/settings/notification/conversation", LTMValidationMiddleware(), community.GetConversationNotificationSettings)
	router.PUT("/community/settings/notification/conversation", LTMValidationMiddleware(), community.EditConversationNotificationSettings)
	router.GET("/community/settings/notification/feed", LTMValidationMiddleware(), community.GetFeedNotificationSettings)
	router.PUT("/community/settings/notification/feed", LTMValidationMiddleware(), community.EditFeedNotificationSettings)
	router.GET("/community/settings/notification", LTMValidationMiddleware(), community.GetNotificationSettings)
	router.PUT("/community/settings/notification", LTMValidationMiddleware(), community.EditNotificationSettings)
	router.GET("/community/tag", LTMValidationMiddleware(), community.GetTaggingList)
	router.GET("/community/settings/content_download", LTMValidationMiddleware(), community.GetContentDownloadSettings)
	router.PUT("/community/settings/content_download", LTMValidationMiddleware(), community.EditContentDownloadSettings)
	router.GET("/community/member/home/meta", LTMValidationMiddleware(), community.MemberHomeMeta)
	router.PUT("/community/member/join", LTMValidationMiddleware(), community.AcceptRejectJoinCommunity)
	router.GET("/community/intro_examples", LTMorVTMValidationMiddleware(), community.GetIntroExamples)
	router.POST("/community/invite", LTMValidationMiddleware(), community.SendCommunityInvite)
	router.GET("/community/configurations", LTMValidationMiddleware(), community.GetCommunityConfigurations)
	router.PATCH("/community/configurations", LTMValidationMiddleware(), community.UpdateCommunityConfigurations)
	router.GET("/community/member/pending", LTMValidationMiddleware(), community.GetPendingCommunityMembers)
	router.GET("/community/removal_reports", LTMValidationMiddleware(), community.GetRemovalReports)
	router.POST("/community/member/:user_id/connection", LTMValidationMiddleware(), community.CreateMemberConnection)
	router.PATCH("/community/member/:user_id/connection", LTMValidationMiddleware(), community.AcceptRejectMemberConnection)
	router.GET("/community/member/:user_id/connection", LTMValidationMiddleware(), community.GetMemberConnection)

	// Moderation Apis
	router.GET("/moderation/rights", LTMValidationMiddleware(), moderation.GetRights)
	router.PUT("/moderation/rights", LTMValidationMiddleware(), moderation.EditRights)
	router.PATCH("/moderation/rights", LTMValidationMiddleware(), moderation.UpdateRights)

	// Conversation Apis
	router.GET("/conversation", LTMValidationMiddleware(), conversation.GetConversation)
	router.POST("/conversation", LTMValidationMiddleware(), conversation.CreateConversation)
	router.PUT("/conversation", LTMValidationMiddleware(), conversation.EditConversation)
	router.DELETE("/conversation", LTMValidationMiddleware(), conversation.DeleteConversation)
	router.PUT("/conversation/reaction", LTMValidationMiddleware(), conversation.AddReaction)
	router.DELETE("/conversation/reaction", LTMValidationMiddleware(), conversation.RemoveReaction)
	router.POST("/conversation/poll", LTMValidationMiddleware(), conversation.AddPoll)
	router.POST("/conversation/poll/submit", LTMValidationMiddleware(), conversation.SubmitPoll)
	router.GET("/conversation/poll/users", LTMValidationMiddleware(), conversation.PollUsers)
	router.PUT("/conversation/topic", LTMValidationMiddleware(), conversation.SetTopic)
	router.PUT("/conversation/event/attend", LTMValidationMiddleware(), conversation.EventAttend)
	router.PUT("/conversation/event/attended", LTMValidationMiddleware(), conversation.EventAttended)
	router.GET("/conversation/event/unseen", LTMValidationMiddleware(), conversation.FetchEventUnseenCount)
	router.PUT("/conversation/event/last_seen", LTMValidationMiddleware(), conversation.UpdateLastSeenEvent)
	router.GET("/conversation/event/link", LTMValidationMiddleware(), conversation.FetchEventLink)
	router.GET("/conversation/event", LTMValidationMiddleware(), conversation.FetchAllEvents)
	router.GET("/conversation/preview/unread", LTMValidationMiddleware(), conversation.FetchUnreadPreviews)
	router.GET("/conversation/preview/unread_count", LTMValidationMiddleware(), conversation.FetchPreviewUnreadMessagesCount)
	router.GET("/conversation/notification/unread", LTMValidationMiddleware(), conversation.UnreadConversationNotification)
	router.GET("/conversation/sync", LTMValidationMiddleware(), conversation.SyncConversation)
	router.GET("/conversation/search", LTMValidationMiddleware(), conversation.ConversationSearch)

	// Feed Apis
	router.POST("/feed/post", LTMValidationMiddleware(), feed.CreatePost)
	router.GET("/feed/post/:post_id", LTMValidationMiddleware(), feed.GetPost)
	router.PUT("/feed/post/:post_id", LTMValidationMiddleware(), feed.EditPost)
	router.DELETE("/feed/post/:post_id", LTMValidationMiddleware(), feed.DeletePost)
	router.PUT("/feed/post/:post_id/like", LTMValidationMiddleware(), feed.CreatePostLike)
	router.GET("/feed/post/:post_id/like", LTMValidationMiddleware(), feed.GetPostLikes)
	router.PUT("/feed/post/:post_id/pin", LTMValidationMiddleware(), feed.PinPost)
	router.PUT("/feed/post/:post_id/save", LTMValidationMiddleware(), feed.CreateSavePost)
	router.POST("/feed/post/:post_id/comment", LTMValidationMiddleware(), feed.CommentPost)
	router.PUT("/feed/post/:post_id/comment/:comment_id", LTMValidationMiddleware(), feed.EditCommentPost)
	router.POST("/feed/post/:post_id/comment/:comment_id/comment", LTMValidationMiddleware(), feed.CreateCommentReply)
	router.GET("/feed/post/:post_id/comment/:comment_id", LTMValidationMiddleware(), feed.GetComment)
	router.DELETE("/feed/post/:post_id/comment/:comment_id", LTMValidationMiddleware(), feed.DeleteComment)
	router.PUT("/feed/post/:post_id/comment/:comment_id/like", LTMValidationMiddleware(), feed.CreateCommentLike)
	router.GET("/feed/post/:post_id/comment/:comment_id/like", LTMValidationMiddleware(), feed.GetCommentLikes)
	router.GET("/feed/user/:user_id/save", LTMValidationMiddleware(), feed.GetSavedPosts)
	router.GET("/feed/user/:user_id/post", LTMValidationMiddleware(), feed.FetchUserCreatedPosts)
	router.GET("/feed/user/activity", LTMValidationMiddleware(), feed.GetUserActivity)
	router.GET("/feed/user/:user_id/activity", LTMValidationMiddleware(), feed.FetchUserProfileActivity)
	router.POST("/feed/user/:user_id/activity", LTMValidationMiddleware(), feed.CreateUserActivity)
	router.GET("/feed/user/activity/unread_count", LTMValidationMiddleware(), feed.GetUserActivityUnreadCount)
	router.POST("/feed/user/activity/:activity_id/mark_read", LTMValidationMiddleware(), feed.UserActivityMarkRead)
	router.GET("/feed/universal", LTMValidationMiddleware(), feed.FetchUniversalFeed)
	router.GET("/feed/group", LTMValidationMiddleware(), feed.FetchGroupFeed)
	router.POST("/feed/topic", LTMValidationMiddleware(), feed.CreateTopics)
	router.GET("/feed/topic", LTMValidationMiddleware(), feed.GetTopic)
	router.DELETE("/feed/topic", LTMValidationMiddleware(), feed.DeleteTopics)
	router.PUT("/feed/topic/:topic_id", LTMValidationMiddleware(), feed.EditTopic)
	router.GET("/feed/connection", LTMValidationMiddleware(), feed.GetConnectionFeed)
	router.POST("feed/post/pending", LTMValidationMiddleware(), feed.CreatePendingPost)

	// Utility Apis
	router.GET("/helper/url", LTMValidationMiddleware(), utility.DecodeUrl)
	router.POST("/helper/media/upload", LTMValidationMiddleware(), utility.UploadFiles)

	// Feedroom Apis
	router.POST("/feedroom", LTMValidationMiddleware(), feedroom.CreateFeedroom)
	router.PUT("/feedroom", LTMValidationMiddleware(), feedroom.EditFeedroom)
	router.DELETE("/feedroom", LTMValidationMiddleware(), feedroom.DeleteFeedroom)
	router.GET("/feedroom", LTMValidationMiddleware(), feedroom.GetFeedroom)
	router.GET("/feedroom/action", LTMValidationMiddleware(), feedroom.GetFeedroomMenu)
	router.GET("/feedroom/settings", LTMValidationMiddleware(), feedroom.GetFeedroomSettings)
	router.PUT("/feedroom/type", LTMValidationMiddleware(), feedroom.ChangeFeedroomType)
	router.GET("/feedroom/type", LTMValidationMiddleware(), feedroom.GetFeedroomTypeStatus)
	router.PUT("/feedroom/enable_member_post", LTMValidationMiddleware(), feedroom.EnableMemberPost)
	router.PUT("/feedroom/pin", LTMValidationMiddleware(), feedroom.PinFeedroom)
	router.PUT("/feedroom/auto_join_members", LTMValidationMiddleware(), feedroom.AutoJoinMembers)
	router.POST("/feedroom/participants", LTMValidationMiddleware(), feedroom.AddParticipants)
	router.GET("/feedroom/participants", LTMValidationMiddleware(), feedroom.GetParticipants)
	router.DELETE("/feedroom/participants", LTMValidationMiddleware(), feedroom.RemoveParticipants)
	router.GET("/feedroom/cohort/access", LTMValidationMiddleware(), feedroom.GetCohortAccess)
	router.PUT("/feedroom/cohort/access", LTMValidationMiddleware(), feedroom.EditCohortAccess)
	router.GET("/feedroom/mine", LTMValidationMiddleware(), feedroom.MyFeedrooms)
	router.PUT("/feedroom/follow", LTMValidationMiddleware(), feedroom.FeedroomFollow)
	router.GET("/feedroom/tag", LTMValidationMiddleware(), feedroom.GetTaggingList)
	router.GET("/feedroom/:feedroom_id/tag", LTMValidationMiddleware(), feedroom.GetTaggingList)

	// Channel Apis
	router.GET("/channel", LTMValidationMiddleware(), channel.FetchChannel)
	router.GET("/channel/invites", LTMValidationMiddleware(), channel.GetChannelInvites)
	router.PUT("/channel/invite", LTMValidationMiddleware(), channel.UpdateChannelInvite)
	router.GET("/channel/:channel_id/settings/member/:participant_uuid", LTMValidationMiddleware(), channel.GetUserChannelSettings)
	router.PUT("/channel/:channel_id/settings/member/:participant_uuid", LTMValidationMiddleware(), channel.UpdateUserChannelSettings)

	// Search Apis
	router.GET("/search/channel", LTMValidationMiddleware(), search.ChannelSearch)
	router.GET("/search/message", LTMValidationMiddleware(), search.MessageSearch)
	router.GET("/search/post", LTMValidationMiddleware(), search.PostSearch)
	router.GET("/search/post/user/:user_id", LTMValidationMiddleware(), search.UserCreatedPostSearch)
	router.GET("/search", LTMValidationMiddleware(), search.GeneralSearch)

	// Widget Apis
	router.POST("/widget", LTMValidationMiddleware(), widget.CreateWidget)
	router.GET("/widget", LTMValidationMiddleware(), widget.GetWidget)
	router.PUT("/widget/:widget_id", LTMValidationMiddleware(), widget.EditWidget)

	// Poll Apis
	router.PUT("/poll/:poll_id", LTMValidationMiddleware(), poll.AddPollOption)
	router.PUT("/poll/:poll_id/vote", LTMValidationMiddleware(), poll.CreatePollVote)
	router.GET("/poll/:poll_id/vote", LTMValidationMiddleware(), poll.GetPollVotes)

	// Webhook Apis
	router.POST("/webhook", LTMValidationMiddleware(), webhook.CreateWebhook)
	router.GET("/webhook", LTMValidationMiddleware(), webhook.GetWebhooks)
	router.GET("/webhook/:webhook_id", LTMValidationMiddleware(), webhook.GetWebhook)
	router.PATCH("/webhook/:webhook_id", LTMValidationMiddleware(), webhook.EditWebhook)
	router.DELETE("/webhook/:webhook_id", LTMValidationMiddleware(), webhook.DeleteWebhook)

	// Logging Apis
	router.POST("/logs", LTMValidationMiddleware(), logger.PushLogs)

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

			// Set API key in request header
			if otm.ApiKey != "" {
				c.Request.Header["X-Api-Key"] = []string{otm.ApiKey}
			}
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
			
			// // Set API key in request header
			if vtm.ApiKey != "" {
				c.Request.Header["X-Api-Key"] = []string{vtm.ApiKey}
			}
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

			// Set API key in request header
			if ltm.ApiKey != "" {
				c.Request.Header["X-Api-Key"] = []string{ltm.ApiKey}
			}
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
