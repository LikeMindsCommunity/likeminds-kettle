package community

const (
	ReportTypeMemberInt = iota
	ReportTypeChatroomInt
	ReportTypeConversationInt
	ReportTypeCommunityInt
	ReportTypeLinkInt
	ReportTypePostInt
	ReportTypeCommentInt
	ReportTypeReplyInt
)

const (
	ReportTypeMember       = "member"
	ReportTypeChatroom     = "chatroom"
	ReportTypeConversation = "conversation"
	ReportTypeCommunity    = "community"
	ReportTypeLink         = "link"
	ReportTypePost         = "post"
	ReportTypeComment      = "comment"
	ReportTypeReply        = "reply"
)

func ReportTypeStrintToInt(reportType string) int {
	switch reportType {
	case ReportTypeMember:
		return ReportTypeMemberInt
	case ReportTypeChatroom:
		return ReportTypeChatroomInt
	case ReportTypeConversation:
		return ReportTypeConversationInt
	case ReportTypeCommunity:
		return ReportTypeCommunityInt
	case ReportTypeLink:
		return ReportTypeLinkInt
	case ReportTypePost:
		return ReportTypePostInt
	case ReportTypeComment:
		return ReportTypeCommentInt
	case ReportTypeReply:
		return ReportTypeReplyInt
	default:
		return -1
	}
}
