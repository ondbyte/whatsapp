package event

var INTERACTIVE MessageType = "interactive"

type InteractiveType string

type Interactive struct {
	ButtonReply ButtonReply     `json:"button_reply"`
	ListReply   ListReply       `json:"list_reply"`
	Type        InteractiveType `json:"type"`
}
