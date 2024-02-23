package message

var (
	TEXT MessageType = `text`
)

type TextMessage struct {
	Body       string `json:"body"`
	PreviewUrl bool   `json:"preview_url"`
}

func NewTextMessage(
	to string,
	body string,
	previewUrl bool,
) *Message {
	return &Message{
		MessagingProduct: WHATSAPP,
		RecipientType:    INDIVIDUAL,
		To:               to,
		Type:             TEXT,
		Text: &TextMessage{
			Body:       body,
			PreviewUrl: previewUrl,
		},
	}
}
