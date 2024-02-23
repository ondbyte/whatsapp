package message

import (
	"fmt"
)

type buttonType string

type reply struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type Button struct {
	Type  string `json:"type"`
	Reply reply  `json:"reply"`
}

func NewReplyButton(id, title string) *Button {
	return &Button{
		Type: "reply",
		Reply: reply{
			Id:    id,
			Title: title,
		},
	}
}

type Body struct {
	Text string `json:"text,omitempty"`
}

func NewBody(format string, a ...any) *Body {
	return &Body{
		Text: fmt.Sprintf(format, a...),
	}
}

type HeaderType string

type Header struct {
	Document Document `json:"document,omitempty"` //place holder
	Image    Image    `json:"image,omitempty"`    //place holder
	Video    Video    `json:"video,omitempty"`    //place holder

	Text string     `json:"text,omitemty"`
	Type HeaderType `json:"type,omitemty"`
}

var TEXT_HEADER HeaderType = "text"

func NewTextHeader(text string) *Header {
	return &Header{
		Type: TEXT_HEADER,
		Text: text,
	}
}

var IMAGE_HEADER HeaderType = "image"

func NewImageHeader(img Image) *Header {
	return &Header{
		Type:  IMAGE_HEADER,
		Image: img,
	}
}

var VIDEO_HEADER HeaderType = "video"

func NewVideoHeader(vid Video) *Header {
	return &Header{
		Type:  VIDEO_HEADER,
		Video: vid,
	}
}

var DOCUMENT_HEADER HeaderType = "document"

func NewDocumentHeader(doc Document) *Header {
	return &Header{
		Type:     DOCUMENT_HEADER,
		Document: doc,
	}
}

type Footer struct {
	Text string `json:"text"`
}

func NewFooter(text string) *Footer {
	return &Footer{
		Text: text,
	}
}

type Action struct {
	Buttons []*Button `json:"buttons"`
}

func NewActionWithButtons(buttons []*Button) *Action {
	return &Action{
		Buttons: buttons,
	}
}

type InteractiveMessageType string

var (
	INTERACTIVE MessageType = `interactive`
)

var (
	BUTTON InteractiveMessageType = "button"
)

type InteractiveMessage struct {
	Type   InteractiveMessageType `json:"type"`
	Body   *Body                  `json:"body"`
	Action *Action                `json:"action"`
	Footer *Footer                `json:"footer"`
	Header *Header                `json:"header"`
}

func NewInteractiveMessageWithButtons(
	to string,
	header *Header,
	body *Body,
	footer *Footer,
	buttons []*Button,
) *Message {
	if len(buttons) > 3 {
		fmt.Printf("maximum of 3 buttons allowed but you have passed %v for NewInteractiveMessageWithButtons, try using list messages\n", len(buttons))
	}
	return &Message{
		MessagingProduct: WHATSAPP,
		RecipientType:    INDIVIDUAL,
		To:               to,
		Type:             INTERACTIVE,
		Interactive: &InteractiveMessage{
			Type:   BUTTON,
			Header: header,
			Body:   body,
			Footer: footer,
			Action: &Action{
				Buttons: buttons,
			},
		},
	}
}
