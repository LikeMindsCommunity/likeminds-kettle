package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/nateshr/likeminds-authentication/internal/cache"
	"github.com/nateshr/likeminds-authentication/internal/constants"
	"github.com/nateshr/likeminds-authentication/internal/logging"
	"github.com/nateshr/likeminds-authentication/internal/middleware"

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
)

var (
	redisClient *redis.Client
	router      *gin.Engine
)

func main() {
	var AppVersion string = "2.33.0"

	initGin()
	redisClient = cache.InitRedis()
	router.Use(cors.New(enableCors()))
	router.Use(middleware.AddResponseHeadersMiddleware())
	router.Use(middleware.ApiMiddleware(redisClient))
	router.Use(middleware.LoggingMiddleware())
	//Attach prometheus service as middleware
	prometheusService := getPrometheusMetricService()
	if prometheusService != nil {
		router.Use(monitoring.PrometheusMiddleware(prometheusService))
	}

	router.GET("", web.Home)

	// OTP Apis
	router.GET(constants.OTPGenerateRoute, otp.GenerateOTP)
	router.GET(constants.OTPVerifyRoute, otp.VerifyOTP)
	router.GET(constants.UserTokenRoute, user.CreateToken)

	// User Apis
	router.POST(constants.UserLoginRoute, middleware.OTMValidationMiddleware(), user.Login)
	router.POST(constants.UserRefreshRoute, middleware.RTMValidationMiddleware(redisClient), user.Refresh)
	router.POST(constants.UserLogoutRoute, middleware.LogoutValidationMiddleware(redisClient), user.Logout)
	router.POST(constants.UserMergeAccountRoute, middleware.LTMValidationMiddleware(redisClient, true), user.MergeAccount)
	router.GET(constants.UserConfigRoute, middleware.LTMValidationMiddleware(redisClient, true), user.Config)
	router.GET(constants.UserBotRoute, middleware.LTMValidationMiddleware(redisClient, true), user.GetBot)
	router.POST(constants.UserDevicePushRoute, middleware.LTMValidationMiddleware(redisClient, true), user.PushUserToken)
	router.POST(constants.UserSubscriptionWhatsappRoute, user.WASubscription)
	router.GET(constants.UserMetaRoute, middleware.LTMValidationMiddleware(redisClient, true), user.UserMeta)
	router.POST(constants.UserOTPRoute, middleware.OTMValidationMiddleware(), user.GenerateUserOTP)
	router.GET(constants.UserOTPVerifyRoute, middleware.OTMValidationMiddleware(), user.VerifyUserOTP)
	router.GET(constants.UserSocialLoginRoute, middleware.OTMValidationMiddleware(), user.UserSocialLogin)

	// Home Apis
	router.POST(constants.HomeFetchCommunitiesRoute, middleware.LTMValidationMiddleware(redisClient, true), home.FetchCommunities)
	router.GET(constants.HomeDmMetaRoute, middleware.LTMValidationMiddleware(redisClient, true), home.DMHome)

	// SDK Apis
	router.POST(constants.SDKInitiateRoute, middleware.VTMValidationMiddleware(false), middleware.RateLimitingMiddleware(redisClient), sdk.InitiateSDK)
	router.GET(constants.SDKInitiateRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.FetchSdkUserInfo)
	router.POST(constants.SDKProjectRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.CreateProject)
	router.GET(constants.SDKProjectRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.GetProject)
	router.PUT(constants.SDKProjectRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.EditProject)
	router.DELETE(constants.SDKProjectRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.DeleteProject)
	router.GET(constants.SDKOnboardingRoute, middleware.OTMValidationMiddleware(), middleware.RateLimitingMiddleware(redisClient), sdk.GetScreen)
	router.POST(constants.SDKOnboardingRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.CreateScreen)
	router.PUT(constants.SDKOnboardingRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.EditScreen)
	router.DELETE(constants.SDKOnboardingRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), sdk.DeleteScreen)

	// Chatroom Apis
	router.GET(constants.ChatroomRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetChatroom)
	router.POST(constants.ChatroomRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.CreateChatroom)
	router.PUT(constants.ChatroomRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.EditChatroom)
	router.DELETE(constants.ChatroomRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.DeleteChatroom)
	router.GET(constants.ChatroomTypeRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetChatroomTypeStatus)
	router.PUT(constants.ChatroomTypeRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ChangeChatroomType)
	router.POST(constants.ChatroomScheduleFollowRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ScheduleFollow)
	router.PUT(constants.ChatroomPinRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.PinChatroom)
	router.GET(constants.ChatroomTagRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetTaggingList)
	router.GET(constants.CHatroomIdTagRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetTaggingList)
	router.GET(constants.ChatroomParticipantsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetParticipants)
	router.POST(constants.ChatroomParticipantsRoute, middleware.LTMValidationMiddleware(redisClient, true), chatroom.AddParticipants)
	router.DELETE(constants.ChatroomParticipantsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.RemoveParticipants)
	router.GET(constants.ChatroomSettingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetChatroomSettings)
	router.PUT(constants.ChatroomSettingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.EditChatroomSettings)
	router.PUT(constants.ChatroomEnableMemberMessageRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.EnableMemberMessage)
	router.PUT(constants.ChatroomAutoFollowMembersRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.AutoFollowMembers)
	router.PUT(constants.ChatroomFilesRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.UpdateFiles)
	router.GET(constants.ChatroomMineRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.MyChatrooms)
	router.PUT(constants.ChatroomSeenRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.CollabcardSeen)
	router.PUT(constants.ChatroomFollowRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ChatroomFollow)
	router.PUT(constants.ChatroomMuteRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.MuteChatroom)
	router.PUT(constants.ChatroomRenameRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.RenameChatroom)
	router.GET(constants.ChatroomShareRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.FetchShareUrl)
	router.GET(constants.ChatroomPendingRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.FetchPendingChatroom)
	router.PUT(constants.ChatroomPendingRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ActionPendingChatroom)
	router.GET(constants.ChatroomSyncRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.SyncChatrooms)
	router.POST(constants.ChatroomDMBlockRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ChatroomBlock)
	router.POST(constants.ChatroomDMRequestRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.InitiatingDMRequest)
	router.POST(constants.ChatroomDMCreateRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.CreateDM)
	router.GET(constants.ChatroomDMRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ListDMChatrooms)
	router.GET(constants.ChatroomDMLimitRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.DMLimit)
	router.GET(constants.ChatroomSearchRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ChatroomSearch)
	router.POST(constants.ChatroomCohortRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.AddCohortToChatroom)
	router.DELETE(constants.ChatroomCohortRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.RemoveCohortFromChatroom)
	router.GET(constants.ChatroomCohortAccessRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetCohortAccess)
	router.PUT(constants.ChatroomCohortAccessRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.EditCohortAccess)
	router.GET(constants.ChatroomHomeRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.GetChatroomHome)
	router.POST(constants.ChatroomMarkReadRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.ChatroomMarkRead)
	router.GET(constants.ChatroomEventRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.FetchEvents)
	router.POST(constants.ChatroomEventRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.CreateEvent)
	router.PUT(constants.ChatroomEventRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.EditEvent)
	router.GET(constants.ChatroomEventMetaRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.FetchEventMeta)
	router.GET(constants.ChatroomEventLinkRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.FetchEventLinks)
	router.GET(constants.ChatroomEventUnseenCountRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.FetchEventUnseenCount)
	router.POST(constants.ChatroomEventRecordingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.UploadEventRecordings)
	router.DELETE(constants.ChatroomEventRecordingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.DeleteEventRecordings)
	router.POST(constants.ChatroomEventRecordingsMetaRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.UploadEventRecordingsMeta)
	router.DELETE(constants.ChatroomEventRecordingsMetaRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.DeleteEventRecordingsMeta)
	router.POST(constants.ChatroomEventInstructorsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.AddEventInstructors)
	router.POST(constants.ChatroomEventHighlightsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.AddEventHighlights)
	router.POST(constants.ChatroomEventTestimonialsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.AddEventTestimonials)
	router.POST(constants.ChatroomEventFAQRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), chatroom.AddEventFAQ)

	// Community Apis
	router.GET(constants.CommunityRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.Community)
	router.GET(constants.CommunityBrandingRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.Branding)
	router.POST(constants.CommunityQuestionsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditQuestions)
	router.GET(constants.CommunityQuestionsRoute, middleware.LTMorVTMValidationMiddleware(), middleware.RateLimitingMiddleware(redisClient), community.GetQuestions)
	router.GET(constants.CommunityQuestionFiltersRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetCommunityQuestionFilters)
	router.GET(constants.CommunityMemberRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetMember)
	router.POST(constants.CommunityMemberRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.AddMember)
	router.DELETE(constants.CommunityMemberRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.RemoveMembers)
	router.DELETE(constants.CommunityMemberLeaveRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.LeaveCommunity)
	router.PUT(constants.CommunityMemberRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditMember)
	router.GET(constants.CommunityMemberStateRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.FetchMemberState)
	router.GET(constants.CommunityMemberRoleRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.FetchMemberRole)
	router.DELETE(constants.CommunityManagerRemoveRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.RemoveCommunityManager)
	router.DELETE(constants.CommunityAdminRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.RemoveCommunityManager)
	router.DELETE(constants.CommunityMemberRemoveRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.RemoveMember)
	router.GET(constants.CommunityManagementToolRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetManagementTools)
	router.GET(constants.CommunityReportRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetReport)
	router.POST(constants.CommunityReportRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.PushReport)
	router.DELETE(constants.CommunityReportRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.CloseReport)
	router.PATCH(constants.CommunityReportRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.UpdateReports)
	router.GET(constants.CommunityReportTagRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetReportTags)
	router.GET(constants.CommunitySettingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetCommunitySettings)
	router.PUT(constants.CommunitySettingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.UpdateCommunitySettings)
	router.GET(constants.CommunityRightsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetCommunityRights)
	router.PUT(constants.CommunityRightsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditCommunityRights)
	router.PATCH(constants.CommunityRightsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.UpdateCommunityRights)
	router.GET(constants.CommunityDMSettingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetCommunityDMSettings)
	router.PUT(constants.CommunityDMSettingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditCommunityDMSettings)
	router.GET(constants.CommunityFeedDMRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.DMFeed)
	router.GET(constants.CommunityDMStatusRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.DMStatus)
	router.GET(constants.CommunityMemberSearchRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.MemberSearch)
	router.GET(constants.CommunityMemberProfileRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetMemberProfile)
	router.PUT(constants.CommunityMemberProfileRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditMemberProfile)
	router.GET(constants.CommunityMemberChatroomRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.MemberChatroom)
	router.GET(constants.CommunityMemberUserIdChannelsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.CommunityMemberChannels)
	router.GET(constants.CommunityMemberChannelStatusRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetMemberChannels)
	router.POST(constants.CommunityCohortRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.CreateCohort)
	router.GET(constants.CommunityCohortRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetCohort)
	router.GET(constants.CommunityCohortIdRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.FetchCohort)
	router.DELETE(constants.CommunityCohortRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.DeleteCohort)
	router.PUT(constants.CommunityCohortRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditCohort)
	router.DELETE(constants.CommunityCohortMemberRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.RemoveCohortMember)
	router.GET(constants.CommunityFeedRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetCommunityFeed)
	router.GET(constants.CommunitySettingsNotificationConversationRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetConversationNotificationSettings)
	router.PUT(constants.CommunitySettingsNotificationConversationRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditConversationNotificationSettings)
	router.GET(constants.CommunitySettingsNotificationFeedRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetFeedNotificationSettings)
	router.PUT(constants.CommunitySettingsNotificationFeedRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditFeedNotificationSettings)
	router.GET(constants.CommunitySettingsNotificationRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetNotificationSettings)
	router.PUT(constants.CommunitySettingsNotificationRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditNotificationSettings)
	router.GET(constants.CommunityTagRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetTaggingList)
	router.GET(constants.CommunityContentDownloadSettingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetContentDownloadSettings)
	router.PUT(constants.CommunityContentDownloadSettingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.EditContentDownloadSettings)
	router.GET(constants.CommunityMemberHomeMetaRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.MemberHomeMeta)
	router.PUT(constants.CommunityMemberJoin, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.AcceptRejectJoinCommunity)
	router.GET(constants.CommunityIntroExamplesRoute, middleware.LTMorVTMValidationMiddleware(), middleware.RateLimitingMiddleware(redisClient), community.GetIntroExamples)
	router.POST(constants.CommunityInviteRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.SendCommunityInvite)
	router.GET(constants.CommunityConfigurationsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetCommunityConfigurations)
	router.PATCH(constants.CommunityConfigurationsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.UpdateCommunityConfigurations)
	router.GET(constants.CommunityMemberPendingRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetPendingCommunityMembers)
	router.GET(constants.CommunityRemovalReportsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetRemovalReports)
	router.POST(constants.CommunityMemberConnectionRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.CreateMemberConnection)
	router.PATCH(constants.CommunityMemberConnectionRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.AcceptRejectMemberConnection)
	router.GET(constants.CommunityMemberConnectionRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), community.GetMemberConnection)

	// Moderation Apis
	router.GET(constants.ModerationRightsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), moderation.GetRights)
	router.PUT(constants.ModerationRightsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), moderation.EditRights)
	router.PATCH(constants.ModerationRightsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), moderation.UpdateRights)

	// Conversation Apis
	router.GET(constants.ConversationRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.GetConversation)
	router.POST(constants.ConversationRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.CreateConversation)
	router.PUT(constants.ConversationRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.EditConversation)
	router.DELETE(constants.ConversationRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.DeleteConversation)
	router.PUT(constants.ConversationReactionRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.AddReaction)
	router.DELETE(constants.ConversationReactionRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.RemoveReaction)
	router.POST(constants.ConversationPollRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.AddPoll)
	router.POST(constants.ConversationPollSubmitRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.SubmitPoll)
	router.GET(constants.ConversationPollUsersRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.PollUsers)
	router.PUT(constants.ConversationTopicRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.SetTopic)
	router.PUT(constants.ConversationEventAttendRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.EventAttend)
	router.PUT(constants.ConversationEventAttendedRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.EventAttended)
	router.GET(constants.ConversationEventUnseenRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.FetchEventUnseenCount)
	router.PUT(constants.ConversationEventLastSeenRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.UpdateLastSeenEvent)
	router.GET(constants.ConversationEventLinkRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.FetchEventLink)
	router.GET(constants.ConversationEventRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.FetchAllEvents)
	router.GET(constants.ConversationPreviewUnreadRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.FetchUnreadPreviews)
	router.GET(constants.ConversationPreviewUnreadCountRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.FetchPreviewUnreadMessagesCount)
	router.GET(constants.ConversationNotificationUnreadRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.UnreadConversationNotification)
	router.GET(constants.ConversationSyncRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.SyncConversation)
	router.GET(constants.ConversationSearchRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), conversation.ConversationSearch)

	// Feed Apis
	router.POST(constants.FeedPostRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreatePost)
	router.GET(constants.FeedPostIDRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feed.GetPost)
	router.PUT(constants.FeedPostIDRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.EditPost)
	router.DELETE(constants.FeedPostIDRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.DeletePost)
	router.PUT(constants.FeedPostIDLikeRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreatePostLike)
	router.GET(constants.FeedPostIDLikeRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetPostLikes)
	router.PUT(constants.FeedPostIDPinRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.PinPost)
	router.PUT(constants.FeedPostIDSaveRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreateSavePost)
	router.POST(constants.FeedPostIDCommentRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CommentPost)
	router.GET(constants.FeedPostIDCommentIDRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feed.GetComment)
	router.DELETE(constants.FeedPostIDCommentIDRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.DeleteComment)
	router.PUT(constants.FeedPostIDCommentIDRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.EditCommentPost)
	router.POST(constants.FeedPostIDCommentIDCommentRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreateCommentReply)
	router.PUT(constants.FeedPostIDCommentIDLikeRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreateCommentLike)
	router.GET(constants.FeedPostIDCommentIDLikeRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetCommentLikes)
	router.GET(constants.FeedUserIDSaveRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetSavedPosts)
	router.GET(constants.FeedUserIDPostRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.FetchUserCreatedPosts)
	router.GET(constants.FeedUserIDCommentRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetUserComments)
	router.GET(constants.FeedUserActivityRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetUserActivity)
	router.GET(constants.FeedUserIdActivityRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.FetchUserProfileActivity)
	router.POST(constants.FeedUserIdActivityRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreateUserActivity)
	router.GET(constants.FeedUserActivityUnreadCount, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetUserActivityUnreadCount)
	router.POST(constants.FeedUserActivityIDMarkReadRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.UserActivityMarkRead)
	router.GET(constants.FeedUserIDMetaRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetUserFeedMeta)
	router.GET(constants.FeedUserIDPostPendingRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.FetchUserCreatedPendingPosts)
	router.GET(constants.FeedUniversalRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feed.FetchUniversalFeed)
	router.GET(constants.FeedGroupRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.FetchGroupFeed)
	router.POST(constants.FeedTopicRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreateTopics)
	router.GET(constants.FeedTopicRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feed.GetTopic)
	router.DELETE(constants.FeedTopicRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.DeleteTopics)
	router.PUT(constants.FeedTopicIDRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.EditTopic)
	router.GET(constants.FeedConnectionRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.GetConnectionFeed)
	router.POST(constants.FeedPostPendingRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.CreatePendingPost)
	router.PUT(constants.FeedPostPendingPostIDRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.EditPendingPost)
	router.GET(constants.FeedPostPendingPostIDRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.FetchPendingPost)
	router.DELETE(constants.FeedPostPendingPostIDRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.DeletePendingPost)
	router.GET(constants.FeedUserTopicsRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.FetchUsersTopics)
	router.PATCH(constants.FeedUserUUIDTopicsRoute, middleware.LTMValidationMiddleware(redisClient, false), middleware.RateLimitingMiddleware(redisClient), feed.UpdateUserTopics)

	// Utility Apis
	router.GET(constants.HelperUrlRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), utility.DecodeUrl)
	router.POST(constants.HelperMediaUploadRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), utility.UploadFiles)
	router.POST(constants.HelperS3UploadRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), utility.UploadFilesToS3)

	// Feedroom Apis
	router.POST(constants.FeedroomRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.CreateFeedroom)
	router.PUT(constants.FeedroomRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.EditFeedroom)
	router.DELETE(constants.FeedroomRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.DeleteFeedroom)
	router.GET(constants.FeedroomRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetFeedroom)
	router.GET(constants.FeedroomActionRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetFeedroomMenu)
	router.GET(constants.FeedroomSettingsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetFeedroomSettings)
	router.PUT(constants.FeedroomTypeRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.ChangeFeedroomType)
	router.GET(constants.FeedroomTypeRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetFeedroomTypeStatus)
	router.PUT(constants.FeedroomEnableMemberPostRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.EnableMemberPost)
	router.PUT(constants.FeedroomPinRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.PinFeedroom)
	router.PUT(constants.FeedroomAutoJoinMembersRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.AutoJoinMembers)
	router.POST(constants.FeedroomParticipantsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.AddParticipants)
	router.GET(constants.FeedroomParticipantsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetParticipants)
	router.DELETE(constants.FeedroomParticipantsRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.RemoveParticipants)
	router.GET(constants.FeedroomCohortAccessRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetCohortAccess)
	router.PUT(constants.FeedroomCohortAccessRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.EditCohortAccess)
	router.GET(constants.FeedroomMineRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.MyFeedrooms)
	router.PUT(constants.FeedroomFollowRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.FeedroomFollow)
	router.GET(constants.FeedroomTagRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetTaggingList)
	router.GET(constants.FeedroomIDTagRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), feedroom.GetTaggingList)

	// Channel Apis
	router.GET(constants.ChannelRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), channel.FetchChannel)
	router.GET(constants.ChannelInvitesRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), channel.GetChannelInvites)
	router.PUT(constants.ChannelInviteRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), channel.UpdateChannelInvite)
	router.GET(constants.ChannelChannelIdSettingsMemberParticipantUUIDRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), channel.GetUserChannelSettings)
	router.PUT(constants.ChannelChannelIdSettingsMemberParticipantUUIDRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), channel.UpdateUserChannelSettings)

	// Search Apis
	router.GET(constants.SearchChannelRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), search.ChannelSearch)
	router.GET(constants.SearchMessageRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), search.MessageSearch)
	router.GET(constants.SearchPostRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), search.PostSearch)
	router.GET(constants.SearchPostUserUserIdRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), search.UserCreatedPostSearch)
	router.GET(constants.SearchRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), search.GeneralSearch)

	// Widget Apis
	router.POST(constants.WidgetRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), widget.CreateWidget)
	router.GET(constants.WidgetRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), widget.GetWidget)
	router.PUT(constants.WidgetIdRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), widget.EditWidget)

	// Poll Apis
	router.PUT(constants.PollIdRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), poll.AddPollOption)
	router.PUT(constants.PollIdVoteRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), poll.CreatePollVote)
	router.GET(constants.PollIdVoteRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), poll.GetPollVotes)

	// Webhook Apis
	router.POST(constants.WebhookRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), webhook.CreateWebhook)
	router.GET(constants.WebhookRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), webhook.GetWebhooks)
	router.GET(constants.WebhookIdRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), webhook.GetWebhook)
	router.PATCH(constants.WebhookIdRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), webhook.EditWebhook)
	router.DELETE(constants.WebhookIdRoute, middleware.LTMValidationMiddleware(redisClient, true), middleware.RateLimitingMiddleware(redisClient), webhook.DeleteWebhook)

	// Logging Apis
	router.POST(constants.LogsEndpoint, middleware.LTMValidationMiddleware(redisClient, true), frontendLogger.PushLogs)

	// Internal Apis
	router.DELETE(constants.CacheRoute, middleware.InternalServiceValidationMiddleware(), internalServices.DeleteCache)

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
