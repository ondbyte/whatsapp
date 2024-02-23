package event

var BUTTON_REPLY InteractiveType = "button_reply"

type ButtonReply struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type ButtonReplyDetails struct {
	// id of the button you passed while sending
	ID string `json:"id"`
	// title of the button you passed while sending
	Title string `json:"title"`
	// contact who replied
	Contact Sender `json:"contact"`
}

func (wba *WhatsAppBusinessAccount) ButtonReply() (*ButtonReplyDetails, bool) {
	value := wba.value()
	if value == nil || len(value.Contacts) == 0 ||
		len(value.Messages) == 0 ||
		value.Messages[0].Interactive.ButtonReply.ID == "" ||
		value.Messages[0].Interactive.ButtonReply.Title == "" ||
		value.Messages[0].Type != INTERACTIVE || value.Messages[0].Interactive.Type != BUTTON_REPLY {
		return nil, false
	}
	return &ButtonReplyDetails{
		ID:      value.Messages[0].Interactive.ButtonReply.ID,
		Title:   value.Messages[0].Interactive.ButtonReply.Title,
		Contact: value.Contacts[0],
	}, true
}
