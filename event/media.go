package event

type media struct {
	Id       string `json:"id"`
	MimeType string `json:"mime_type"`
}

type Image media
type Audio media
type Video media
type Document media
type Sticker media

var (
	IMAGE    MessageType = "image"
	AUDIO    MessageType = "audio"
	VIDEO    MessageType = "video"
	DOCUMENT MessageType = "document"
	STICKER  MessageType = "sticker"
)

// returns if the event is image event/msg
func (w *WhatsAppBusinessAccount) Image() (*Image, bool) {
	value := w.value()
	if value == nil || len(value.Messages) == 0 || value.Messages[0].Type != IMAGE {
		return nil, false
	}
	return &value.Messages[0].Image, true
}

// returns if the event is audio event/msg
func (w *WhatsAppBusinessAccount) Audio() (*Audio, bool) {
	value := w.value()
	if value == nil || len(value.Messages) == 0 || value.Messages[0].Type != AUDIO {
		return nil, false
	}
	return &value.Messages[0].Audio, true
}

// returns if the event is video event/msg
func (w *WhatsAppBusinessAccount) Video() (*Video, bool) {
	value := w.value()
	if value == nil || len(value.Messages) == 0 || value.Messages[0].Type != VIDEO {
		return nil, false
	}
	return &value.Messages[0].Video, true
}

// returns if the event is document event/msg
func (w *WhatsAppBusinessAccount) Document() (*Document, bool) {
	value := w.value()
	if value == nil || len(value.Messages) == 0 || value.Messages[0].Type != DOCUMENT {
		return nil, false
	}
	return &value.Messages[0].Document, true
}

// returns if the event is sticker event/msg
func (w *WhatsAppBusinessAccount) Sticker() (*Sticker, bool) {
	value := w.value()
	if value == nil || len(value.Messages) == 0 || value.Messages[0].Type != STICKER {
		return nil, false
	}
	return &value.Messages[0].Sticker, true
}
