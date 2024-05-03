package utils

import (
	"encoding/json"
	"fmt"

	"github.com/nateshr/likeminds-authentication/constants"
	"github.com/nateshr/likeminds-authentication/logging"
)

type OGTags struct {
	Title       string `json:"title"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

type AttachmentMeta struct {
	Name                 string                 `json:"name"`
	Url                  string                 `json:"url"`
	Format               string                 `json:"format"`
	Size                 int                    `json:"size"`
	Duration             int                    `json:"duration"`
	PageCount            int                    `json:"page_count"`
	ThumbnailUrl         string                 `json:"thumbnail_url"`
	OgTags               OGTags                 `json:"og_tags"`
	EntityID             string                 `json:"entity_id"`
	CoverImageUrl        string                 `json:"cover_image_url"`
	Title                string                 `json:"title"`
	Body                 string                 `json:"body"`
	Options              []string               `json:"options"`
	ExpiryTime           int64                  `json:"expiry_time"`
	PollType             string                 `json:"poll_type"`
	MultipleSelectState  string                 `json:"multiple_select_state"`
	MultipleSelectNumber int                    `json:"multiple_select_number"`
	IsAnonymous          bool                   `json:"is_anonymous"`
	AllowAddOption       bool                   `json:"allow_add_option"`
	WidgetMeta           map[string]interface{} `json:"widget_meta"`
}

type AttachmentRequest struct {
	AttachmentType int             `json:"attachment_type"`
	AttachmentMeta *AttachmentMeta `json:"attachment_meta"`
	Type           string          `json:"type"`
	MetaData       *AttachmentMeta `json:"meta_data"`
}

func ConvertAttachmentMetaForCustomWidgetAttachments(attachments []AttachmentRequest, rawData []byte) []AttachmentRequest {

	// Unmarshal widgets data for attachment type custom widget
	widgetData := make(map[string]interface{})
	err := json.Unmarshal(rawData, &widgetData)
	if err != nil {
		logging.Error(fmt.Sprint("Error in unmarshalling widgets data: ", err))
		return attachments
	}

	for i := 0; i < len(attachments); i++ {
		if attachments[i].AttachmentType == constants.CustomWidget {

			if attachments[i].AttachmentMeta != nil {
				widget_meta := widgetData["attachments"].([]interface{})[i].(map[string]interface{})["attachment_meta"].(map[string]interface{})

				if widget_meta != nil {
					delete(widget_meta, "entity_id")
					attachments[i].AttachmentMeta = &AttachmentMeta{EntityID: attachments[i].AttachmentMeta.EntityID}
					attachments[i].AttachmentMeta.WidgetMeta = widget_meta
				}
			}
		}
	}

	return attachments
}
