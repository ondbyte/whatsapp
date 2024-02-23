package event

var TEXT MessageType = "text"

type TextMessage struct {
	Body string `json:"body"`
}

type TextMessageDetails struct {
	Contact Sender `json:"contact"`
	Body    string `json:"body"`
}

// returns if the event is TextMessage
func (wba *WhatsAppBusinessAccount) TextMessage() (*TextMessageDetails, bool) {
	value := wba.value()
	if value == nil || len(value.Contacts) == 0 || len(value.Messages) == 0 || value.Messages[0].Type != TEXT {
		return nil, false
	}
	return &TextMessageDetails{
		Contact: value.Contacts[0],
		Body:    value.Messages[0].Text.Body,
	}, true
}
