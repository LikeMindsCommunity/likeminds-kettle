package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/cache"
	"github.com/nateshr/likeminds-authentication/chatroom"
	"github.com/nateshr/likeminds-authentication/conversation"
	"github.com/nateshr/likeminds-authentication/home"
	"github.com/nateshr/likeminds-authentication/otp"
	"github.com/nateshr/likeminds-authentication/sdk"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
	"github.com/nateshr/likeminds-authentication/web"
)

var (
	client *redis.Client
	router = gin.Default()
)

func main() {
	client = cache.InitRedis()
	router.Use(ApiMiddleware(client))
	router.GET("", web.Home)
	router.GET("/otp/generate", otp.GenerateOTP)
	router.GET("/otp/verify", otp.VerifyOTP)
	router.POST("/user/login", VTMValidationMiddleware(), user.Login)
	router.POST("/user/refresh", RTMValidationMiddleware(), user.Refresh)
	router.POST("/user/logout", LogoutValidationMiddleware(client), user.Logout)
	router.POST("/user/merge_account", LTMValidationMiddleware(client), user.MergeAccount)
	router.POST("/user/create_bot", APIKeyValidationMiddleware(), user.CreateBot)
	router.POST("/home/fetch_communities", LTMValidationMiddleware(client), home.FetchCommunities)
	router.POST("/sdk/initiate", APIKeyValidationMiddleware(), sdk.InitiateSDK)
	router.POST("/sdk/create", LTMValidationMiddleware(client), sdk.CreateSDK)
	router.POST("/chatroom/schedule_follow", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), chatroom.ScheduleFollow)
	router.POST("/chatroom/create", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), chatroom.CreateChatroom)
	router.GET("/chatroom/fetch", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), chatroom.FetchChatroom)
	router.POST("/chatroom/edit", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), chatroom.EditChatroom)
	router.POST("/chatroom/pin", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), chatroom.PinChatroom)
	router.GET("/chatroom/get_tagging_list", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), chatroom.GetTaggingList)
	router.POST("/chatroom/auto_follow_for_all_members", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), chatroom.AutoFollowForAllMembers)
	router.GET("/chatroom/fetch_participants_meta", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), chatroom.FetchParticipantsMeta)
	router.POST("/chatroom/add", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), chatroom.AddChatroom)
	router.POST("/chatroom/enable_member_message", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), chatroom.EnableMemberMessage)
	router.POST("/chatroom/update_files", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), chatroom.UpdateFiles)
	router.GET("/v2/fetch_chatroom", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), chatroom.FetchChatroomV2)
	router.GET("/v1/fetch_chatroom_feed", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), chatroom.FetchChatroomFeed)
	router.GET("/fetch_community_chatroom_feed", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), chatroom.FetchCommunityChatroomFeed)
	router.GET("/v1/my_chatrooms", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), chatroom.MyChatrooms)
	router.POST("/collabcard_seen", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), chatroom.CollabcardSeen)
	router.GET("/collabcard_follow", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), chatroom.CollabcardFollow)
	// router.POST("/chatroom_mute", LTMValidationMiddleware(client), chatroom.MuteChatroom)
	// router.POST("/chatroom_rename", LTMValidationMiddleware(client), chatroom.RenameChatroom)
	// router.POST("/chatroom_delete", LTMValidationMiddleware(client), chatroom.DeleteChatroom)
	router.GET("/fetch_share_url", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), chatroom.FetchShareUrl)
	router.GET("/fetch_pending_chatroom", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), chatroom.FetchPendingChatroom)
	// router.POST("/action_pending_chatroom", LTMValidationMiddleware(client), chatroom.ActionPendingChatroom)
	router.GET("/sync_chatrooms", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), chatroom.SyncChatrooms)
	router.GET("/sync_chatrooms_diff", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), chatroom.SyncChatroomsDiff)
	router.GET("/conversation/fetch", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), conversation.FetchConversation)
	// router.POST("/conversation/create", LTMValidationMiddleware(client), conversation.CreateConversation)
	// router.POST("/conversation/add_reaction", LTMValidationMiddleware(client), conversation.AddReaction)
	// router.POST("/conversation/remove_reaction", LTMValidationMiddleware(client), conversation.RemoveReaction)
	router.POST("/conversation/add_poll", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), conversation.AddPoll)
	router.POST("/conversation/submit_poll", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), conversation.SubmitPoll)
	router.GET("/conversation/poll_users", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), conversation.PollUsers)
	router.POST("/conversation/set_topic", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), conversation.SetTopic)
	router.POST("/conversation/event/attend", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), conversation.EventAttend)
	router.POST("/conversation/event/attended", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), conversation.EventAttended)
	router.GET("/conversation/event/fetch_unseen_count", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), conversation.FetchEventUnseenCount)
	router.POST("/conversation/event/update_last_seen_event", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), conversation.UpdateLastSeenEvent)
	router.GET("/conversation/event/fetch_link", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), conversation.FetchEventLink)
	router.GET("/conversation/event/fetch_all", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), conversation.FetchAllEvents)
	router.GET("/conversation/fetch_unread_previews", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), conversation.FetchUnreadPreviews)
	router.GET("/conversation/fetch_preview_unread_messages_count", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), conversation.FetchPreviewUnreadMessagesCount)
	router.GET("/conversation_meta", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), conversation.ConversationMeta)
	router.GET("/unread_conversation_notification", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), conversation.UnreadConversationNotification)
	router.POST("/delete_conversation", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), conversation.DeleteConversation)
	// router.POST("/edit_conversation", LTMValidationMiddleware(client), conversation.EditConversation)
	router.GET("/sync_conversation", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), conversation.SyncConversation)
	router.GET("/sync_conversation_diff", LTMValidationMiddleware(client), APIKeyValidationMiddleware(), conversation.SyncConversationDiff)

	log.Fatal(router.Run(":8080"))
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

func LTMValidationMiddleware(client *redis.Client) gin.HandlerFunc {
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
		//Create internal API client
		client := api_client.NewAPIClient()
		options := api_client.GetRequestOptions{
			Url:           client.CoreServiceBaseURL + SDKAuthenticateEndPoint,
			CustomHeaders: utils.CreateHeaders(c),
		}
		//Send request
		respBytes, err := client.GetRequest(&options)
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
