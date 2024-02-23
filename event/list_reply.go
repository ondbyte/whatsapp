package event

var LIST_REPLY InteractiveType

type ListReply struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ListReplyDetails struct {
	// id of the list item you passed while sending
	ID string `json:"id"`
	// title of the list item you passed while sending
	Title string `json:"title"`
	// description of the list item you passed while sending
	Description string `json:"description"`
	// contact who replied
	Contact Sender `json:"contact"`
}

func (wba *WhatsAppBusinessAccount) ListRepy() (*ListReplyDetails, bool) {
	value := wba.value()
	if value == nil || len(value.Contacts) == 0 ||
		len(value.Messages) == 0 ||
		value.Messages[0].Interactive.ListReply.ID == "" ||
		value.Messages[0].Interactive.ListReply.Title == "" ||
		value.Messages[0].Type != INTERACTIVE || value.Messages[0].Interactive.Type != LIST_REPLY {
		return nil, false
	}
	return &ListReplyDetails{
		ID:          value.Messages[0].Interactive.ListReply.ID,
		Title:       value.Messages[0].Interactive.ListReply.Title,
		Description: value.Messages[0].Interactive.ListReply.Description,
		Contact:     value.Contacts[0],
	}, true
}
