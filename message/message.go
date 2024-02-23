package message

import (
	"fmt"
)

type MessageType string
type RecipientType string
type MessagingProduct string

var INDIVIDUAL RecipientType = "individual"

var WHATSAPP MessagingProduct = "whatsapp"

type context struct {
	MessageId string `json:"message_id"`
}

func NewReplyMsg(wamIdOfTheMsgToReplyTo string, msg Message) Message {
	if msg.Type == REACTION {
		fmt.Println("reaction message itself is a reply msg, cannot make it a reply msg")
		return msg
	}
	msg.Context = &context{
		MessageId: wamIdOfTheMsgToReplyTo,
	}
	return msg
}

// Message represents the overall structure for the provided JSON
type Message struct {
	MessagingProduct MessagingProduct    `json:"messaging_product,omitempty"`
	RecipientType    RecipientType       `json:"recipient_type,omitempty"`
	To               string              `json:"to,omitempty"`
	Context          *context            `json:"context,omitempty"`
	Type             MessageType         `json:"type,omitempty"`
	Text             *TextMessage        `json:"text,omitempty"`
	Reaction         *Reaction           `json:"reaction,omitempty"`
	Image            *Image              `json:"image,omitempty"`
	Video            *Video              `json:"video,omitempty"`
	Audio            *Audio              `json:"audio,omitempty"`
	Document         *Document           `json:"document,omitempty"`
	Sticker          *Sticker            `json:"sticker,omitempty"`
	Location         *Location           `json:"location,omitempty"`
	Contacts         Contacts            `json:"contacts,omitempty"`
	Interactive      *InteractiveMessage `json:"interactive,omitempty"`
}

// NewMessage is a constructor function for creating a new message instance.
func NewMessage(
	messagingProduct MessagingProduct,
	recipientType RecipientType,
	to string,
	messageType MessageType,
	text *TextMessage,
	reaction *Reaction,
	image *Image,
	video *Video,
	audio *Audio,
	document *Document,
	sticker *Sticker,
	location *Location,
	contacts Contacts,
	interactive *InteractiveMessage,
) *Message {
	return &Message{
		MessagingProduct: messagingProduct,
		RecipientType:    recipientType,
		To:               to,
		Type:             messageType,
		Text:             text,
		Reaction:         reaction,
		Image:            image,
		Video:            video,
		Audio:            audio,
		Document:         document,
		Sticker:          sticker,
		Location:         location,
		Contacts:         contacts,
		Interactive:      interactive,
	}
}
