package message

type Media struct {
	Id string `json:"id"`
	//Link     string `json:"link"`
	Caption  string `json:"caption"`
	FileName string `json:"filename"`
}

var IMAGE MessageType = "image"

func NewImageMsg(to string, image *Image) *Message {
	return &Message{
		MessagingProduct: WHATSAPP,
		To:               to,
		Type:             IMAGE,
		Image:            image,
		RecipientType:    INDIVIDUAL,
	}
}

var VIDEO MessageType = "video"

func NewVideoMsg(to string, video *Video) *Message {
	return &Message{
		MessagingProduct: WHATSAPP,
		To:               to,
		Type:             VIDEO,
		Video:            video,
		RecipientType:    INDIVIDUAL,
	}
}

var AUDIO MessageType = "audio"

func NewAudioMsg(to string, audio *Audio) *Message {
	return &Message{
		MessagingProduct: WHATSAPP,
		To:               to,
		Type:             AUDIO,
		Audio:            audio,
		RecipientType:    INDIVIDUAL,
	}
}

var DOCUMENT MessageType = "document"

func NewDocumentMsg(to string, doc *Document) *Message {
	return &Message{
		MessagingProduct: WHATSAPP,
		To:               to,
		Type:             DOCUMENT,
		Document:         doc,
		RecipientType:    INDIVIDUAL,
	}
}

var STICKER MessageType = "sticker"

func NewStickerMsg(to string, sticker *Sticker) *Message {
	return &Message{
		MessagingProduct: WHATSAPP,
		To:               to,
		Type:             STICKER,
		Sticker:          sticker,
		RecipientType:    INDIVIDUAL,
	}
}
