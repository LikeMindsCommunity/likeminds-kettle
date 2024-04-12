package utility

const DecodeUrlEndPoint = "/api/decode_url"
const UploadFilesEndpoint = "/api/v1/upload_files"
const UserMetaInfoInternalEndpoint = "/api/community/users"

const ParamUrl = "url"
const ParamMemberIDs = "member_ids"
const ParamCommunityID = "community_id"

// Categories for file uploads
const (
	CategoryFeed = "feed"
)

// Entities for file uploads
const (
	EntityPost   = "post"
	EntityWidget = "widget"
)

// source for file uploads
const (
	SourceGDrive = "gdrive"
)

// file paths for S3 uploads
const (
	FeedPostFilePath   = "files/post/%s"
	FeedWidgetFilePath = "files/widget/%s"
)

// helper constants for S3 file upload
const (
	MaxFilesPerUpload             = 10
	FilterFileIDFromDriveUrlRegex = `\/(?:file\/d\/|open\?id=|uc\?export=download&id=)([a-zA-Z0-9_-]+)`
)
