package whatsapp

import (
	"github.com/ondbyte/whatsapp/event"
	"github.com/ondbyte/whatsapp/message"
)

type msgSender func(*message.Message) (*SendMessageResponse, error)
type ender func()

type Conversation struct {
	waId    string
	receive chan *event.WhatsAppBusinessAccount
	send    msgSender
	end     ender
}

// wait for the next msg from the user
func (c *Conversation) NextEvent() chan *event.WhatsAppBusinessAccount {
	return c.receive
}

// send text msg to the user
func (c *Conversation) SendText(text string, previewUrlInTheText bool) (*SendMessageResponse, error) {
	return c.send(message.NewTextMessage(c.waId, text, previewUrlInTheText))
}

// send location msg to the user
func (c *Conversation) SendLocation(l *message.Location) (*SendMessageResponse, error) {
	return c.send(message.NewLocationMessage(c.waId, l))
}

// send contacts msg to the user
func (c *Conversation) SendContacts(cs message.Contacts) (*SendMessageResponse, error) {
	return c.send(message.NewContactsMsg(c.waId, cs))
}

// send image msg to the user
func (c *Conversation) SendImage(i *message.Image) (*SendMessageResponse, error) {
	return c.send(message.NewImageMsg(c.waId, i))
}

// send audio msg to the user
func (c *Conversation) SendAudio(a *message.Audio) (*SendMessageResponse, error) {
	return c.send(message.NewAudioMsg(c.waId, a))
}

// send video msg to the user
func (c *Conversation) SendVideo(cs *message.Video) (*SendMessageResponse, error) {
	return c.send(message.NewVideoMsg(c.waId, cs))
}

// send document msg to the user
func (c *Conversation) SendDocument(cs *message.Document) (*SendMessageResponse, error) {
	return c.send(message.NewDocumentMsg(c.waId, cs))
}

// send sticker msg to the user
func (c *Conversation) SendSticker(cs *message.Sticker) (*SendMessageResponse, error) {
	return c.send(message.NewStickerMsg(c.waId, cs))
}

// ends the conversation, after this any messages will be diverted to the main whatsapp messages listener
func (c *Conversation) End() {
	c.end()
}

// starts a conversation or returns the existing one becuase only one conversation with a user can exist with a specific WaID/number,
// returns a Conversation with which you can receive msgs and send msgs
func (wa *Whatsapp) StartConversation(waId string) *Conversation {
	wa.userMessageChannelsLock.Lock()
	convo, ok := wa.conversations[waId]
	if !ok {
		ch := make(chan *event.WhatsAppBusinessAccount)
		convo = &Conversation{
			waId:    waId,
			receive: ch,
			send:    wa.sendMessage,
			end: func() {
				delete(wa.conversations, waId)
				close(ch)
			},
		}
		wa.conversations[waId] = convo
	}
	wa.userMessageChannelsLock.Unlock()
	return convo
}
