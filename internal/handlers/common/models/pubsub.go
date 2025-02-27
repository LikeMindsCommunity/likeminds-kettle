package models

type PSRequest struct {
	TopicMessageType string `json:"topic_message_type"`
	RawData          []byte `json:"raw_data"`
}
