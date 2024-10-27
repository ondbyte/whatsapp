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

type Parameters struct {
	DisplayText string `json:"display_text"`
	Url         string `json:"url"`
}

type Action struct {
	Buttons    []*Button  `json:"buttons,omitempty"`
	Name       string     `json:"name,omitempty"`
	Parameters Parameters `json:"parameters,omitempty"`
}

func NewActionWithCTAButton(displayText, ctaUrl string) *Action {
	return &Action{
		Name: "cta_url",
		Parameters: Parameters{
			DisplayText: displayText,
			Url:         ctaUrl,
		},
	}
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
	BUTTON  InteractiveMessageType = "button"
	CTA_URL InteractiveMessageType = "cta_url"
)

type Interactive struct {
	Type   InteractiveMessageType `json:"type"`
	Body   *Body                  `json:"body"`
	Action *Action                `json:"action"`
	Footer *Footer                `json:"footer"`
	Header *Header                `json:"header"`
}

func NewInteractivWithButtons(
	header *Header,
	body *Body,
	footer *Footer,
	buttons []*Button,
) *Interactive {
	if len(buttons) > 3 {
		fmt.Printf("maximum of 3 buttons allowed but you have passed %v for NewInteractiveMessageWithButtons, try using list messages\n", len(buttons))
	}
	return &Interactive{
		Type:   BUTTON,
		Header: header,
		Body:   body,
		Footer: footer,
		Action: &Action{
			Buttons: buttons,
		},
	}
}

func NewInteractiveWithCtaButton(
	header *Header,
	body *Body,
	footer *Footer,
	ctaAction *Action,
) *Interactive {
	return &Interactive{
		Type:   CTA_URL,
		Header: header,
		Body:   body,
		Footer: footer,
		Action: ctaAction,
	}
}

func NewInteractiveMessage(
	to string,
	interactive *Interactive,
) *Message {
	return &Message{
		MessagingProduct: WHATSAPP,
		RecipientType:    INDIVIDUAL,
		To:               to,
		Type:             INTERACTIVE,
		Interactive:      interactive,
	}
}
