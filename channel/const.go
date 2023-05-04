package channel

const ChatroomInvitesEndppoint = "/api/chatroom/invites"

const ParamPage = "page"
const ParamChannelType = "channel_type"
const ParamChannelId = "channel_id"
const ParamMemberUUID = "member_uuid"

const UserChannelSettingsEndpoint = "/api/chatroom/%s/settings/member/%s"

const (
	CHAT_BASED_CHANNEL = 1
	FEED_BASED_CHANNEL = 2
)
