package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/koding/websocketproxy"
	"github.com/nateshr/likeminds-authentication/internal/cache"
	"github.com/nateshr/likeminds-authentication/internal/constants"
	"github.com/nateshr/likeminds-authentication/internal/handlers/channel"
	"github.com/nateshr/likeminds-authentication/internal/handlers/chatroom"
	"github.com/nateshr/likeminds-authentication/internal/handlers/community"
	"github.com/nateshr/likeminds-authentication/internal/handlers/conversation"
	"github.com/nateshr/likeminds-authentication/internal/handlers/feed"
	"github.com/nateshr/likeminds-authentication/internal/handlers/feedroom"
	"github.com/nateshr/likeminds-authentication/internal/handlers/home"
	"github.com/nateshr/likeminds-authentication/internal/handlers/internalServices"
	"github.com/nateshr/likeminds-authentication/internal/handlers/moderation"
	"github.com/nateshr/likeminds-authentication/internal/handlers/otp"
	"github.com/nateshr/likeminds-authentication/internal/handlers/poll"
	"github.com/nateshr/likeminds-authentication/internal/handlers/sdk"
	"github.com/nateshr/likeminds-authentication/internal/handlers/search"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/handlers/utility"
	"github.com/nateshr/likeminds-authentication/internal/handlers/utility/frontendLogger"
	"github.com/nateshr/likeminds-authentication/internal/handlers/utility/monitoring"
	"github.com/nateshr/likeminds-authentication/internal/handlers/web"
	"github.com/nateshr/likeminds-authentication/internal/handlers/webhook"
	"github.com/nateshr/likeminds-authentication/internal/handlers/widget"
	"github.com/nateshr/likeminds-authentication/internal/logging"
	"github.com/nateshr/likeminds-authentication/internal/middleware"
	"github.com/nateshr/likeminds-authentication/internal/utils/api_client"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"net/url"
)

var (
	routerA     *gin.Engine
	routerB     *gin.Engine
	routerGroup errgroup.Group
	redisClient *redis.Client
)

func main() {
	var AppVersion = "2.37.1"

	redisClient = cache.InitRedis()

	initRouterA()
	initRouterB()

	//logging.Fatal(routerA.Run(":8080"))

	routerGroup.Go(func() error {
		return routerAServer().ListenAndServe()
	})
	routerGroup.Go(func() error {
		return routerBServer().ListenAndServe()
	})
	if err := routerGroup.Wait(); err != nil {
		log.Fatal(err)
	}

	logging.Info(fmt.Sprintf("application version: %s", AppVersion))
}

func initRouterA() {
	gin.SetMode(gin.ReleaseMode)
	routerA = gin.Default()
	setRouterA()
}

func initRouterB() {
	gin.SetMode(gin.ReleaseMode)
	routerB = gin.Default()
	setRouterB()
}

func setRouterA() {
	routerA.Use(cors.New(enableCors()))
	routerA.Use(middleware.AddResponseHeadersMiddleware())
	routerA.Use(middleware.ApiMiddleware(redisClient))
	routerA.Use(middleware.LoggingMiddleware())
	//Attach prometheus service as middleware
	prometheusService := getPrometheusMetricService()
	if prometheusService != nil {
		routerA.Use(monitoring.PrometheusMiddleware(prometheusService))
	}

	routerA.GET("", web.Home)

	// OTP Apis
	routerA.GET(constants.OTPGenerateRoute, otp.GenerateOTP)
	routerA.GET(constants.OTPVerifyRoute, otp.VerifyOTP)
	routerA.GET(constants.UserTokenRoute, user.CreateToken)

	// User Apis
	routerA.POST(constants.UserLoginRoute, middleware.OTMValidationMiddleware(), user.Login)
	routerA.POST(constants.UserRefreshRoute, middleware.RTMValidationMiddleware(redisClient), user.Refresh)
	routerA.POST(constants.UserLogoutRoute, middleware.LogoutValidationMiddleware(redisClient), user.Logout)
	routerA.POST(constants.UserMergeAccountRoute, middleware.LTMValidationMiddleware(redisClient, true), user.MergeAccount)
	routerA.GET(constants.UserConfigRoute, middleware.LTMValidationMiddleware(redisClient, true), user.Config)
	routerA.GET(constants.UserBotRoute, middleware.LTMValidationMiddleware(redisClient, true), user.GetBot)
	routerA.POST(constants.UserDevicePushRoute, middleware.LTMValidationMiddleware(redisClient, true), user.PushUserToken)
	routerA.POST(constants.UserSubscriptionWhatsappRoute, user.WASubscription)
	routerA.GET(constants.UserMetaRoute, middleware.LTMValidationMiddleware(redisClient, true), user.UserMeta)
	routerA.POST(constants.UserOTPRoute, middleware.OTMValidationMiddleware(), user.GenerateUserOTP)
	routerA.GET(constants.UserOTPVerifyRoute, middleware.OTMValidationMiddleware(), user.VerifyUserOTP)
	routerA.GET(constants.UserSocialLoginRoute, middleware.OTMValidationMiddleware(), user.UserSocialLogin)

	// Home Apis
	routerA.POST(constants.HomeFetchCommunitiesRoute, middleware.LTMValidationMiddleware(redisClient, true), home.FetchCommunities)
	routerA.GET(constants.HomeDmMetaRoute, middleware.LTMValidationMiddleware(redisClient, true), home.DMHome)

	// SDK Apis
	routerA.POST(constants.SDKInitiateRoute, middleware.VTMValidationMiddleware(false), middleware.RateLimitingMiddleware(redisClient), sdk.InitiateSDK)
	routerA.GET(constants.SDKInitiateRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.FetchSdkUserInfo)
	routerA.POST(constants.SDKProjectRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.CreateProject)
	routerA.GET(constants.SDKProjectRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.GetProject)
	routerA.PUT(constants.SDKProjectRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.EditProject)
	routerA.DELETE(constants.SDKProjectRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.DeleteProject)
	routerA.GET(constants.SDKOnboardingRoute, middleware.OTMValidationMiddleware(), middleware.RateLimitingMiddleware(redisClient), sdk.GetScreen)
	routerA.POST(constants.SDKOnboardingRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.CreateScreen)
	routerA.PUT(constants.SDKOnboardingRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.EditScreen)
	routerA.DELETE(constants.SDKOnboardingRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.DeleteScreen)

	// Chatroom Apis
	routerA.GET(constants.ChatroomRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetChatroom)
	routerA.POST(constants.ChatroomRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.CreateChatroom)
	routerA.PUT(constants.ChatroomRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.EditChatroom)
	routerA.DELETE(constants.ChatroomRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.DeleteChatroom)
	routerA.GET(constants.ChatroomTypeRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetChatroomTypeStatus)
	routerA.PUT(constants.ChatroomTypeRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ChangeChatroomType)
	routerA.POST(constants.ChatroomScheduleFollowRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ScheduleFollow)
	routerA.PUT(constants.ChatroomPinRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.PinChatroom)
	routerA.GET(constants.ChatroomTagRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetTaggingList)
	routerA.GET(constants.CHatroomIdTagRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetTaggingList)
	routerA.GET(constants.ChatroomParticipantsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetParticipants)
	routerA.POST(constants.ChatroomParticipantsRoute, middleware.LTMValidationMiddleware(redisClient, true), chatroom.AddParticipants)
	routerA.DELETE(constants.ChatroomParticipantsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.RemoveParticipants)
	routerA.GET(constants.ChatroomSettingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetChatroomSettings)
	routerA.PUT(constants.ChatroomSettingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.EditChatroomSettings)
	routerA.PUT(constants.ChatroomEnableMemberMessageRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.EnableMemberMessage)
	routerA.PUT(constants.ChatroomAutoFollowMembersRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.AutoFollowMembers)
	routerA.PUT(constants.ChatroomFilesRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.UpdateFiles)
	routerA.GET(constants.ChatroomMineRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.MyChatrooms)
	routerA.PUT(constants.ChatroomSeenRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.CollabcardSeen)
	routerA.PUT(constants.ChatroomFollowRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ChatroomFollow)
	routerA.PUT(constants.ChatroomMuteRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.MuteChatroom)
	routerA.PUT(constants.ChatroomRenameRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.RenameChatroom)
	routerA.GET(constants.ChatroomShareRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.FetchShareUrl)
	routerA.GET(constants.ChatroomPendingRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.FetchPendingChatroom)
	routerA.PUT(constants.ChatroomPendingRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ActionPendingChatroom)
	routerA.GET(constants.ChatroomSyncRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.SyncChatrooms)
	routerA.POST(constants.ChatroomDMBlockRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ChatroomBlock)
	routerA.POST(constants.ChatroomDMRequestRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.InitiatingDMRequest)
	routerA.POST(constants.ChatroomDMCreateRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.CreateDM)
	routerA.GET(constants.ChatroomDMRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ListDMChatrooms)
	routerA.GET(constants.ChatroomDMLimitRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.DMLimit)
	routerA.GET(constants.ChatroomSearchRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ChatroomSearch)
	routerA.POST(constants.ChatroomCohortRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.AddCohortToChatroom)
	routerA.DELETE(constants.ChatroomCohortRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.RemoveCohortFromChatroom)
	routerA.GET(constants.ChatroomCohortAccessRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetCohortAccess)
	routerA.PUT(constants.ChatroomCohortAccessRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.EditCohortAccess)
	routerA.GET(constants.ChatroomHomeRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetChatroomHome)
	routerA.POST(constants.ChatroomMarkReadRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ChatroomMarkRead)
	routerA.GET(constants.ChatroomEventRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.FetchEvents)
	routerA.POST(constants.ChatroomEventRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.CreateEvent)
	routerA.PUT(constants.ChatroomEventRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.EditEvent)
	routerA.GET(constants.ChatroomEventMetaRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.FetchEventMeta)
	routerA.GET(constants.ChatroomEventLinkRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.FetchEventLinks)
	routerA.GET(constants.ChatroomEventUnseenCountRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.FetchEventUnseenCount)
	routerA.POST(constants.ChatroomEventRecordingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.UploadEventRecordings)
	routerA.DELETE(constants.ChatroomEventRecordingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.DeleteEventRecordings)
	routerA.POST(constants.ChatroomEventRecordingsMetaRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.UploadEventRecordingsMeta)
	routerA.DELETE(constants.ChatroomEventRecordingsMetaRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.DeleteEventRecordingsMeta)
	routerA.POST(constants.ChatroomEventInstructorsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.AddEventInstructors)
	routerA.POST(constants.ChatroomEventHighlightsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.AddEventHighlights)
	routerA.POST(constants.ChatroomEventTestimonialsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.AddEventTestimonials)
	routerA.POST(constants.ChatroomEventFAQRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.AddEventFAQ)

	// Community Apis
	routerA.GET(constants.CommunityRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.Community)
	routerA.GET(constants.CommunityBrandingRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.Branding)
	routerA.POST(constants.CommunityQuestionsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditQuestions)
	routerA.GET(constants.CommunityQuestionsRoute, middleware.LTMorVTMValidationMiddleware(), middleware.RateLimitingMiddleware(redisClient), community.GetQuestions)
	routerA.GET(constants.CommunityQuestionFiltersRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetCommunityQuestionFilters)
	routerA.GET(constants.CommunityMemberRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetMember)
	routerA.POST(constants.CommunityMemberRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.AddMember)
	routerA.DELETE(constants.CommunityMemberRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.RemoveMembers)
	routerA.DELETE(constants.CommunityMemberLeaveRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.LeaveCommunity)
	routerA.PUT(constants.CommunityMemberRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditMember)
	routerA.GET(constants.CommunityMemberStateRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.FetchMemberState)
	routerA.GET(constants.CommunityMemberRoleRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.FetchMemberRole)
	routerA.DELETE(constants.CommunityManagerRemoveRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.RemoveCommunityManager)
	routerA.DELETE(constants.CommunityAdminRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.RemoveCommunityManager)
	routerA.DELETE(constants.CommunityMemberRemoveRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.RemoveMember)
	routerA.GET(constants.CommunityManagementToolRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetManagementTools)
	routerA.GET(constants.CommunityReportRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetReport)
	routerA.POST(constants.CommunityReportRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.PushReport)
	routerA.DELETE(constants.CommunityReportRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.CloseReport)
	routerA.PATCH(constants.CommunityReportRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.UpdateReports)
	routerA.GET(constants.CommunityReportTagRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetReportTags)
	routerA.GET(constants.CommunitySettingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetCommunitySettings)
	routerA.PUT(constants.CommunitySettingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.UpdateCommunitySettings)
	routerA.GET(constants.CommunityRightsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetCommunityRights)
	routerA.PUT(constants.CommunityRightsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditCommunityRights)
	routerA.PATCH(constants.CommunityRightsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.UpdateCommunityRights)
	routerA.GET(constants.CommunityDMSettingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetCommunityDMSettings)
	routerA.PUT(constants.CommunityDMSettingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditCommunityDMSettings)
	routerA.GET(constants.CommunityFeedDMRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.DMFeed)
	routerA.GET(constants.CommunityDMStatusRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.DMStatus)
	routerA.GET(constants.CommunityMemberSearchRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.MemberSearch)
	routerA.GET(constants.CommunityMemberProfileRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetMemberProfile)
	routerA.PUT(constants.CommunityMemberProfileRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditMemberProfile)
	routerA.GET(constants.CommunityMemberChatroomRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.MemberChatroom)
	routerA.GET(constants.CommunityMemberUserIdChannelsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.CommunityMemberChannels)
	routerA.GET(constants.CommunityMemberChannelStatusRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetMemberChannels)
	routerA.POST(constants.CommunityCohortRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.CreateCohort)
	routerA.GET(constants.CommunityCohortRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetCohort)
	routerA.GET(constants.CommunityCohortIdRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.FetchCohort)
	routerA.DELETE(constants.CommunityCohortRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.DeleteCohort)
	routerA.PUT(constants.CommunityCohortRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditCohort)
	routerA.DELETE(constants.CommunityCohortMemberRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.RemoveCohortMember)
	routerA.GET(constants.CommunityFeedRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetCommunityFeed)
	routerA.GET(constants.CommunitySettingsNotificationConversationRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetConversationNotificationSettings)
	routerA.PUT(constants.CommunitySettingsNotificationConversationRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditConversationNotificationSettings)
	routerA.GET(constants.CommunitySettingsNotificationFeedRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetFeedNotificationSettings)
	routerA.PUT(constants.CommunitySettingsNotificationFeedRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditFeedNotificationSettings)
	routerA.GET(constants.CommunitySettingsNotificationRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetNotificationSettings)
	routerA.PUT(constants.CommunitySettingsNotificationRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditNotificationSettings)
	routerA.GET(constants.CommunityTagRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetTaggingList)
	routerA.GET(constants.CommunityContentDownloadSettingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetContentDownloadSettings)
	routerA.PUT(constants.CommunityContentDownloadSettingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditContentDownloadSettings)
	routerA.GET(constants.CommunityMemberHomeMetaRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.MemberHomeMeta)
	routerA.PUT(constants.CommunityMemberJoin, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.AcceptRejectJoinCommunity)
	routerA.GET(constants.CommunityIntroExamplesRoute, middleware.LTMorVTMValidationMiddleware(), middleware.RateLimitingMiddleware(redisClient), community.GetIntroExamples)
	routerA.POST(constants.CommunityInviteRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.SendCommunityInvite)
	routerA.GET(constants.CommunityConfigurationsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetCommunityConfigurations)
	routerA.PATCH(constants.CommunityConfigurationsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.UpdateCommunityConfigurations)
	routerA.GET(constants.CommunityMemberPendingRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetPendingCommunityMembers)
	routerA.GET(constants.CommunityRemovalReportsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetRemovalReports)
	routerA.POST(constants.CommunityMemberConnectionRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.CreateMemberConnection)
	routerA.PATCH(constants.CommunityMemberConnectionRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.AcceptRejectMemberConnection)
	routerA.GET(constants.CommunityMemberConnectionRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetMemberConnection)

	// Moderation Apis
	routerA.GET(constants.ModerationRightsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), moderation.GetRights)
	routerA.PUT(constants.ModerationRightsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), moderation.EditRights)
	routerA.PATCH(constants.ModerationRightsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), moderation.UpdateRights)

	// Conversation Apis
	routerA.GET(constants.ConversationRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.GetConversation)
	routerA.POST(constants.ConversationRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.CreateConversation)
	routerA.PUT(constants.ConversationRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.EditConversation)
	routerA.DELETE(constants.ConversationRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.DeleteConversation)
	routerA.PUT(constants.ConversationReactionRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.AddReaction)
	routerA.DELETE(constants.ConversationReactionRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.RemoveReaction)
	routerA.POST(constants.ConversationPollRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.AddPoll)
	routerA.POST(constants.ConversationPollSubmitRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.SubmitPoll)
	routerA.GET(constants.ConversationPollUsersRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.PollUsers)
	routerA.PUT(constants.ConversationTopicRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.SetTopic)
	routerA.PUT(constants.ConversationEventAttendRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.EventAttend)
	routerA.PUT(constants.ConversationEventAttendedRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.EventAttended)
	routerA.GET(constants.ConversationEventUnseenRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.FetchEventUnseenCount)
	routerA.PUT(constants.ConversationEventLastSeenRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.UpdateLastSeenEvent)
	routerA.GET(constants.ConversationEventLinkRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.FetchEventLink)
	routerA.GET(constants.ConversationEventRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.FetchAllEvents)
	routerA.GET(constants.ConversationPreviewUnreadRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.FetchUnreadPreviews)
	routerA.GET(constants.ConversationPreviewUnreadCountRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.FetchPreviewUnreadMessagesCount)
	routerA.GET(constants.ConversationNotificationUnreadRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.UnreadConversationNotification)
	routerA.GET(constants.ConversationSyncRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.SyncConversation)
	routerA.GET(constants.ConversationSearchRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.ConversationSearch)

	// Feed Apis
	routerA.POST(constants.FeedPostRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreatePost)
	routerA.GET(constants.FeedPostIDRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feed.GetPost)
	routerA.PUT(constants.FeedPostIDRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.EditPost)
	routerA.DELETE(constants.FeedPostIDRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.DeletePost)
	routerA.PUT(constants.FeedPostIDLikeRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreatePostLike)
	routerA.GET(constants.FeedPostIDLikeRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetPostLikes)
	routerA.PUT(constants.FeedPostIDPinRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.PinPost)
	routerA.PUT(constants.FeedPostIDSaveRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreateSavePost)
	routerA.POST(constants.FeedPostIDCommentRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CommentPost)
	routerA.GET(constants.FeedPostIDCommentIDRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feed.GetComment)
	routerA.DELETE(constants.FeedPostIDCommentIDRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.DeleteComment)
	routerA.PUT(constants.FeedPostIDCommentIDRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.EditCommentPost)
	routerA.POST(constants.FeedPostIDCommentIDCommentRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreateCommentReply)
	routerA.PUT(constants.FeedPostIDCommentIDLikeRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreateCommentLike)
	routerA.GET(constants.FeedPostIDCommentIDLikeRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetCommentLikes)
	routerA.GET(constants.FeedUserIDSaveRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetSavedPosts)
	routerA.GET(constants.FeedUserIDPostRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.FetchUserCreatedPosts)
	routerA.GET(constants.FeedUserIDCommentRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetUserComments)
	routerA.GET(constants.FeedUserActivityRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetUserActivity)
	routerA.GET(constants.FeedUserIdActivityRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.FetchUserProfileActivity)
	routerA.POST(constants.FeedUserIdActivityRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreateUserActivity)
	routerA.GET(constants.FeedUserActivityUnreadCount, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetUserActivityUnreadCount)
	routerA.POST(constants.FeedUserActivityIDMarkReadRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.UserActivityMarkRead)
	routerA.GET(constants.FeedUserIDMetaRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetUserFeedMeta)
	routerA.GET(constants.FeedUserIDPostPendingRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.FetchUserCreatedPendingPosts)
	routerA.GET(constants.FeedUniversalRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feed.FetchUniversalFeed)
	routerA.GET(constants.FeedGroupRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.FetchGroupFeed)
	routerA.POST(constants.FeedTopicRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreateTopics)
	routerA.GET(constants.FeedTopicRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feed.GetTopic)
	routerA.DELETE(constants.FeedTopicRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.DeleteTopics)
	routerA.PUT(constants.FeedTopicIDRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.EditTopic)
	routerA.GET(constants.FeedConnectionRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetConnectionFeed)
	routerA.POST(constants.FeedPostPendingRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreatePendingPost)
	routerA.PUT(constants.FeedPostPendingPostIDRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.EditPendingPost)
	routerA.GET(constants.FeedPostPendingPostIDRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.FetchPendingPost)
	routerA.DELETE(constants.FeedPostPendingPostIDRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.DeletePendingPost)
	routerA.GET(constants.FeedUserTopicsRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.FetchUsersTopics)
	routerA.PATCH(constants.FeedUserUUIDTopicsRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.UpdateUserTopics)
	routerA.GET(constants.FeedPersonalisedRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.FetchPersonalisedFeed)
	routerA.POST(constants.FeedPersonalisedRecomputeRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.RecomputePersonalisedFeed)
	routerA.POST(constants.FeedPersonalisedReorderRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.ReorderPersonalisedFeed)
	routerA.POST(constants.FeedPostSeenRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.SeenPost)

	// Utility Apis
	routerA.GET(constants.HelperUrlRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), utility.DecodeUrl)
	routerA.POST(constants.HelperMediaUploadRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), utility.UploadFiles)
	routerA.POST(constants.HelperS3UploadRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), utility.UploadFilesToS3)

	// Feedroom Apis
	routerA.POST(constants.FeedroomRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.CreateFeedroom)
	routerA.PUT(constants.FeedroomRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.EditFeedroom)
	routerA.DELETE(constants.FeedroomRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.DeleteFeedroom)
	routerA.GET(constants.FeedroomRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetFeedroom)
	routerA.GET(constants.FeedroomActionRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetFeedroomMenu)
	routerA.GET(constants.FeedroomSettingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetFeedroomSettings)
	routerA.PUT(constants.FeedroomTypeRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.ChangeFeedroomType)
	routerA.GET(constants.FeedroomTypeRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetFeedroomTypeStatus)
	routerA.PUT(constants.FeedroomEnableMemberPostRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.EnableMemberPost)
	routerA.PUT(constants.FeedroomPinRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.PinFeedroom)
	routerA.PUT(constants.FeedroomAutoJoinMembersRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.AutoJoinMembers)
	routerA.POST(constants.FeedroomParticipantsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.AddParticipants)
	routerA.GET(constants.FeedroomParticipantsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetParticipants)
	routerA.DELETE(constants.FeedroomParticipantsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.RemoveParticipants)
	routerA.GET(constants.FeedroomCohortAccessRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetCohortAccess)
	routerA.PUT(constants.FeedroomCohortAccessRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.EditCohortAccess)
	routerA.GET(constants.FeedroomMineRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.MyFeedrooms)
	routerA.PUT(constants.FeedroomFollowRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.FeedroomFollow)
	routerA.GET(constants.FeedroomTagRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetTaggingList)
	routerA.GET(constants.FeedroomIDTagRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetTaggingList)

	// Channel Apis
	routerA.GET(constants.ChannelRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), channel.FetchChannel)
	routerA.GET(constants.ChannelInvitesRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), channel.GetChannelInvites)
	routerA.PUT(constants.ChannelInviteRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), channel.UpdateChannelInvite)
	routerA.GET(constants.ChannelChannelIdSettingsMemberParticipantUUIDRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), channel.GetUserChannelSettings)
	routerA.PUT(constants.ChannelChannelIdSettingsMemberParticipantUUIDRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), channel.UpdateUserChannelSettings)

	// Search Apis
	routerA.GET(constants.SearchChannelRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), search.ChannelSearch)
	routerA.GET(constants.SearchMessageRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), search.MessageSearch)
	routerA.GET(constants.SearchPostRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), search.PostSearch)
	routerA.GET(constants.SearchPostUserUserIdRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), search.UserCreatedPostSearch)
	routerA.GET(constants.SearchRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), search.GeneralSearch)

	// Widget Apis
	routerA.POST(constants.WidgetRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), widget.CreateWidget)
	routerA.GET(constants.WidgetRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), widget.GetWidget)
	routerA.PUT(constants.WidgetIdRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), widget.EditWidget)

	// Poll Apis
	routerA.PUT(constants.PollIdRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), poll.AddPollOption)
	routerA.PUT(constants.PollIdVoteRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), poll.CreatePollVote)
	routerA.GET(constants.PollIdVoteRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), poll.GetPollVotes)

	// Webhook Apis
	routerA.POST(constants.WebhookRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), webhook.CreateWebhook)
	routerA.GET(constants.WebhookRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), webhook.GetWebhooks)
	routerA.GET(constants.WebhookIdRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), webhook.GetWebhook)
	routerA.PATCH(constants.WebhookIdRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), webhook.EditWebhook)
	routerA.DELETE(constants.WebhookIdRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), webhook.DeleteWebhook)

	// Logging Apis
	routerA.POST(constants.LogsEndpoint, middleware.LTMValidationMiddleware(redisClient, true), frontendLogger.PushLogs)

	// Internal Apis
	routerA.DELETE(constants.CacheRoute, middleware.InternalServiceValidationMiddleware(), internalServices.DeleteCache)

	routerA.GET("/metrics", gin.WrapH(promhttp.Handler()))
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

func setRouterB() {
	psURL, _ := url.Parse(api_client.GetPandemoniumServiceBaseUrl())
	psWSProxy := websocketproxy.NewProxy(psURL)
	// Pandemonium APIs
	routerB.GET(constants.SubscribeRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), gin.WrapH(psWSProxy))
}

func routerAServer() *http.Server {
	serverA := &http.Server{
		Addr:    ":8080",
		Handler: routerA,
	}
	return serverA
}
func routerBServer() *http.Server {
	serverB := &http.Server{
		Addr:    ":8083",
		Handler: routerB,
	}
	return serverB
}
