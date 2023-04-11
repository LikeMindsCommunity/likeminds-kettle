package main

import (
	"log"
	"net/http"

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
)

var (
	client *redis.Client
	router *gin.Engine
)

func main() {
	var AppVersion string = "1.19.1"

	initGin()
	client = cache.InitRedis()
	router.Use(cors.New(enableCors()))
	router.Use(ApiMiddleware(client))
	router.GET("", web.Home)

	//OTP Apis
	router.GET("/otp/generate", otp.GenerateOTP)
	router.GET("/otp/verify", otp.VerifyOTP)

	//User Apis
	router.POST("/user/login", VTMValidationMiddleware(), user.Login)
	router.POST("/user/refresh", RTMValidationMiddleware(), user.Refresh)
	router.POST("/user/logout", LogoutValidationMiddleware(client), user.Logout)
	router.POST("/user/merge_account", LTMValidationMiddleware(client, true), user.MergeAccount)
	router.GET("/user/config", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), user.Config)
	router.GET("/user/bot", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), user.GetBot)
	router.POST("/user/device/push", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), user.PushUserToken)
	router.POST("/user/subscription/whatsapp", user.WASubscription)

	//Home Apis
	router.POST("/home/fetch_communities", LTMValidationMiddleware(client, true), home.FetchCommunities)
	router.GET("/home/dm/meta", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), home.DMHome)

	//SDK Apis
	router.POST("/sdk/initiate", APIKeyValidationMiddleware(), sdk.InitiateSDK)
	router.POST("/sdk/project", LTMValidationMiddleware(client, true), sdk.CreateProject)
	router.GET("/sdk/project", LTMValidationMiddleware(client, true), sdk.GetProject)
	router.PUT("/sdk/project", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), sdk.EditProject)
	router.DELETE("/sdk/project", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), sdk.DeleteProject)
	router.GET("/sdk/onboarding", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), sdk.GetScreen)
	router.POST("/sdk/onboarding", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), sdk.CreateScreen)
	router.PUT("/sdk/onboarding", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), sdk.EditScreen)
	router.DELETE("/sdk/onboarding", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), sdk.DeleteScreen)

	//Chatroom Apis
	router.GET("/chatroom", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.GetChatroom)
	router.POST("/chatroom", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.CreateChatroom)
	router.PUT("/chatroom", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.EditChatroom)
	router.DELETE("/chatroom", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.DeleteChatroom)
	router.GET("/chatroom/type", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.GetChatroomTypeStatus)
	router.PUT("/chatroom/type", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.ChangeChatroomType)
	router.POST("/chatroom/schedule_follow", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.ScheduleFollow)
	router.PUT("/chatroom/pin", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.PinChatroom)
	router.GET("/chatroom/tag", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.GetTaggingList)
	router.GET("/chatroom/participants", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.GetParticipants)
	router.POST("/chatroom/participants", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.AddParticipants)
	router.DELETE("/chatroom/participants", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.RemoveParticipants)
	router.GET("/chatroom/settings", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.GetChatroomSettings)
	router.PUT("/chatroom/settings", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.EditChatroomSettings)
	router.PUT("/chatroom/enable_member_message", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.EnableMemberMessage)
	router.PUT("/chatroom/auto_follow_members", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.AutoFollowMembers)
	router.PUT("/chatroom/files", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.UpdateFiles)
	router.GET("/chatroom/mine", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.MyChatrooms)
	router.PUT("/chatroom/seen", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.CollabcardSeen)
	router.PUT("/chatroom/follow", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.ChatroomFollow)
	router.PUT("/chatroom/mute", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.MuteChatroom)
	router.PUT("/chatroom/rename", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.RenameChatroom)
	router.GET("/chatroom/share", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.FetchShareUrl)
	router.GET("/chatroom/pending", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.FetchPendingChatroom)
	router.PUT("/chatroom/pending", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.ActionPendingChatroom)
	router.GET("/chatroom/sync", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.SyncChatrooms)
	router.POST("/chatroom/dm/block", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.ChatroomBlock)
	router.POST("/chatroom/dm/request", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.InitiatingDMRequest)
	router.POST("/chatroom/dm/create", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.CreateDM)
	router.GET("/chatroom/dm", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.ListDMChatrooms)
	router.GET("/chatroom/dm/limit", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.DMLimit)
	router.GET("/chatroom/search", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.ChatroomSearch)
	router.POST("/chatroom/cohort", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.AddCohortToChatroom)
	router.DELETE("/chatroom/cohort", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.RemoveCohortFromChatroom)
	router.GET("/chatroom/cohort/access", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.GetCohortAccess)
	router.PUT("/chatroom/cohort/access", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.EditCohortAccess)
	router.GET("/chatroom/home", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.GetChatroomHome)
	router.POST("/chatroom/mark_read", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), chatroom.ChatroomMarkRead)

	//Community Apis
	router.POST("/community/questions", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.EditQuestions)
	router.GET("/community/questions", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.GetQuestions)
	router.GET("/community/member", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.GetMember)
	router.POST("/community/member", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.AddMember)
	router.PUT("/community/member", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.EditMember)
	router.GET("/community/member/state", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.FetchMemberState)
	router.DELETE("/community/manager/remove", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.RemoveCommunityManager)
	router.DELETE("/community/member/remove", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.RemoveMember)
	router.GET("/community/management/tool", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.GetManagementTools)
	router.GET("/community/report", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.GetReport)
	router.POST("/community/report", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.PushReport)
	router.DELETE("/community/report", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.CloseReport)
	router.GET("/community/report/tag", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.GetReportTags)
	router.GET("/community/settings", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.GetCommunitySettings)
	router.PUT("/community/settings", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.UpdateCommunitySettings)
	router.PUT("/community/rights", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.EditCommunityRights)
	router.GET("/community/settings/dm", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.GetCommunityDMSettings)
	router.PUT("/community/settings/dm", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.EditCommunityDMSettings)
	router.GET("/community/feed/dm", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.DMFeed)
	router.GET("/community/dm/status", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.DMStatus)
	router.GET("/community/member/search", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.MemberSearch)
	router.GET("/community/member/profile", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.GetMemberProfile)
	router.PUT("/community/member/profile", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.EditMemberProfile)
	router.GET("/community/member/chatroom", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.MemberChatroom)
	router.GET("/community/member/:user_id/channel", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.CommunityMemberChannels)
	router.GET("/community/member/channel/status", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.GetMemberChannels)
	router.POST("/community/cohort", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.CreateCohort)
	router.GET("/community/cohort", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.GetCohort)
	router.DELETE("/community/cohort", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.DeleteCohort)
	router.PUT("/community/cohort", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.EditCohort)
	router.DELETE("/community/cohort/member", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.RemoveCohortMember)
	router.GET("/community/feed", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.GetCommunityFeed)
	router.GET("/community/settings/notification/conversation", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.GetConversationNotificationSettings)
	router.PUT("/community/settings/notification/conversation", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.EditConversationNotificationSettings)
	router.GET("/community/settings/notification/feed", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.GetFeedNotificationSettings)
	router.PUT("/community/settings/notification/feed", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.EditFeedNotificationSettings)
	router.GET("/community/tag", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), community.GetTaggingList)

	//Moderation Apis
	router.GET("/moderation/rights", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), moderation.GetRights)
	router.PUT("/moderation/rights", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), moderation.EditRights)

	//Conversation Apis
	router.GET("/conversation", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), conversation.GetConversation)
	router.POST("/conversation", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), conversation.CreateConversation)
	router.PUT("/conversation", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), conversation.EditConversation)
	router.DELETE("/conversation", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), conversation.DeleteConversation)
	router.PUT("/conversation/reaction", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), conversation.AddReaction)
	router.DELETE("/conversation/reaction", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), conversation.RemoveReaction)
	router.POST("/conversation/poll", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), conversation.AddPoll)
	router.POST("/conversation/poll/submit", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), conversation.SubmitPoll)
	router.GET("/conversation/poll/users", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), conversation.PollUsers)
	router.PUT("/conversation/topic", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), conversation.SetTopic)
	router.PUT("/conversation/event/attend", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), conversation.EventAttend)
	router.PUT("/conversation/event/attended", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), conversation.EventAttended)
	router.GET("/conversation/event/unseen", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), conversation.FetchEventUnseenCount)
	router.PUT("/conversation/event/last_seen", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), conversation.UpdateLastSeenEvent)
	router.GET("/conversation/event/link", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), conversation.FetchEventLink)
	router.GET("/conversation/event", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), conversation.FetchAllEvents)
	router.GET("/conversation/preview/unread", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), conversation.FetchUnreadPreviews)
	router.GET("/conversation/preview/unread_count", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), conversation.FetchPreviewUnreadMessagesCount)
	router.GET("/conversation/notification/unread", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), conversation.UnreadConversationNotification)
	router.GET("/conversation/sync", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), conversation.SyncConversation)
	router.GET("/conversation/search", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), conversation.ConversationSearch)

	//Feed Apis
	router.POST("/feed/post", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feed.CreatePost)
	router.GET("/feed/post/:post_id", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feed.GetPost)
	router.PUT("/feed/post/:post_id", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feed.EditPost)
	router.DELETE("/feed/post/:post_id", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feed.DeletePost)
	router.PUT("/feed/post/:post_id/like", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feed.CreatePostLike)
	router.GET("/feed/post/:post_id/like", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feed.GetPostLikes)
	router.PUT("/feed/post/:post_id/pin", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feed.PinPost)
	router.PUT("/feed/post/:post_id/save", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feed.CreateSavePost)
	router.POST("/feed/post/:post_id/comment", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feed.CommentPost)
	router.PUT("/feed/post/:post_id/comment/:comment_id", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feed.EditCommentPost)
	router.POST("/feed/post/:post_id/comment/:comment_id/comment", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feed.CreateCommentReply)
	router.GET("/feed/post/:post_id/comment/:comment_id", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feed.GetComment)
	router.DELETE("/feed/post/:post_id/comment/:comment_id", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feed.DeleteComment)
	router.PUT("/feed/post/:post_id/comment/:comment_id/like", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feed.CreateCommentLike)
	router.GET("/feed/post/:post_id/comment/:comment_id/like", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feed.GetCommentLikes)
	router.GET("/feed/user/:user_id/save", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feed.GetSavedPosts)
	router.GET("/feed/user/:user_id/post", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feed.FetchUserCreatedPosts)
	router.POST("/feed/user/:user_id/activity", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feed.CreateUserActivity)
	router.GET("/feed/universal", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feed.FetchUniversalFeed)
	router.GET("/feed/group", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feed.FetchGroupFeed)

	//Utility Apis
	router.GET("/helper/url", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), utility.DecodeUrl)
	router.POST("/helper/media/upload", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), utility.UploadFiles)

	//Feedroom Apis
	router.POST("/feedroom", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feedroom.CreateFeedroom)
	router.PUT("/feedroom", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feedroom.EditFeedroom)
	router.DELETE("/feedroom", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feedroom.DeleteFeedroom)
	router.GET("/feedroom", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feedroom.GetFeedroom)
	router.GET("/feedroom/action", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feedroom.GetFeedroomMenu)
	router.GET("/feedroom/settings", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feedroom.GetFeedroomSettings)
	router.PUT("/feedroom/type", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feedroom.ChangeFeedroomType)
	router.GET("/feedroom/type", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feedroom.GetFeedroomTypeStatus)
	router.PUT("/feedroom/enable_member_post", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feedroom.EnableMemberPost)
	router.PUT("/feedroom/pin", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feedroom.PinFeedroom)
	router.PUT("/feedroom/auto_join_members", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feedroom.AutoJoinMembers)
	router.POST("/feedroom/participants", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feedroom.AddParticipants)
	router.GET("/feedroom/participants", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feedroom.GetParticipants)
	router.DELETE("/feedroom/participants", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feedroom.RemoveParticipants)
	router.GET("/feedroom/cohort/access", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feedroom.GetCohortAccess)
	router.PUT("/feedroom/cohort/access", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feedroom.EditCohortAccess)
	router.GET("/feedroom/mine", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feedroom.MyFeedrooms)
	router.PUT("/feedroom/follow", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feedroom.FeedroomFollow)
	router.GET("/feedroom/tag", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), feedroom.GetTaggingList)

	// Channel Apis
	router.GET("/channel", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), channel.FetchChannel)
	router.GET("/channel/invites", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), channel.GetChannelInvites)
	router.PUT("/channel/invite", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), channel.UpdateChannelInvite)

	// Search Apis
	router.GET("/search/channel", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), search.ChannelSearch)
	router.GET("/search/message", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), search.MessageSearch)
	router.GET("/search/post", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), search.PostSearch)
	router.GET("/search/post/user/:user_id", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), search.UserCreatedPostSearch)
	router.GET("/search", LTMValidationMiddleware(client, true), APIKeyValidationMiddleware(), search.GeneralSearch)

	log.Printf("application version: %s", AppVersion)
	log.Fatal(router.Run(":8080"))
}

func initGin() {
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
}

// ApiMiddleware will add the db connection to the context
func ApiMiddleware(client *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(cache.ParamRedisClient, client)
		c.Next()
	}
}

func VTMValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Extract VTM from token, internally it checks if token is valid or not
		vtm, err := token.ExtractVTM(c.Request.Header.Get(token.HeaderAuthorization))
		if vtm == nil {
			log.Print(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidVTM,
			})
			return
		} else {
			//If valid, set "vtm" in context, to be used in later APIs
			c.Set(token.ParamVTM, vtm)
		}
		c.Next()
	}
}

func LTMValidationMiddleware(client *redis.Client, emptyBearerTokenCheck bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get(token.HeaderAuthorization)
		//If bearer token is empty, let it pass through
		if !emptyBearerTokenCheck && len(bearerToken) == 0 {
			c.Next()
			return
		}
		//Extract LTM from token, internally it checks if token is valid or not
		ltm, err := token.ExtractLTM(bearerToken)
		if ltm == nil {
			log.Print(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidLTM,
			})
			return
		} else {
			//Check if LTM is black listed or not
			if cache.IsLTMBlacklisted(client, ltm) {
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
			log.Print(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidRTM,
			})
			return
		} else {
			//Check if RTM is black listed or not
			if cache.IsRTMBlacklisted(client, rtm) {
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

func LogoutValidationMiddleware(client *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Extract LTM from token, internally it checks if token is valid or not
		ltm, err := token.ExtractLTM(c.Request.Header.Get(token.HeaderAuthorization))
		if ltm == nil {
			log.Print(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
				Success:      false,
				ErrorMessage: token.ErrorInvalidLTM,
			})
			return
		} else {
			//Check if LTM is black listed or not
			if cache.IsLTMBlacklisted(client, ltm) {
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
				log.Print(err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
					Success:      false,
					ErrorMessage: token.ErrorInvalidRTM,
				})
				return
			} else {
				//Check if RTM is black listed or not
				if cache.IsRTMBlacklisted(client, rtm) {
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
		//Check if request has LTM token or not
		ltm, ok := c.Get(token.ParamLTM)
		if ok && ltm.(*token.LoginTokenMeta).ApiKey != "" {
			c.Request.Header["X-Api-Key"] = []string{ltm.(*token.LoginTokenMeta).ApiKey}
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
		"x-accept-version",
		"x-username",
		"x-password",
		"x-device-id",
		"x-api-key",
		"Authorization",
	)
	return config
}
