package message

type Reaction struct {
	MessageId string `json:"message_id"`
	Emoji     string `json:"emoji"`
}

var REACTION MessageType = "reaction"

func NewReaction(msgIdToReactTo, emoji string) *Reaction {
	return &Reaction{
		MessageId: msgIdToReactTo,
		Emoji:     emoji,
	}
}

func NewReactionMsg(to string, r *Reaction) *Message {
	return &Message{
		MessagingProduct: WHATSAPP,
		RecipientType:    INDIVIDUAL,
		To:               to,
		Type:             REACTION,
		Reaction:         r,
	}
}
