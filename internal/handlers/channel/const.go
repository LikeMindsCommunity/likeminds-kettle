package channel

const ChatroomInvitesEndppoint = "/api/chatroom/invites"
const SyncChannelDetailEndppoint = "/api/sync/channel_detail"

const ParamPage = "page"
const ParamChannelType = "channel_type"
const ParamChannelId = "channel_id"
const ParamChannelActionTypes = "channel_action_types"
const ParamParticipantUUID = "participant_uuid"

const UserChannelSettingsEndpoint = "/api/chatroom/%s/settings/member/%s"

const (
	CHAT_BASED_CHANNEL = 1
	FEED_BASED_CHANNEL = 2
)
