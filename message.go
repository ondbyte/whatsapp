package whatsapp

import (
	"fmt"

	event "github.com/ondbyte/whatsapp/event"
	"github.com/ondbyte/whatsapp/message"
	"github.com/ondbyte/whatsapp/util"
)

type SendMessageResponse struct {
	MessagingProduct string        `json:"messaging_product"`
	Contacts         []Contact     `json:"contacts"`
	Messages         []SentMessage `json:"messages"`
}

// Valid implements util.Body.
func (sr *SendMessageResponse) Valid() bool {
	return len(sr.Messages) > 0
}

type Contact struct {
	Input string `json:"input"`
	WaID  string `json:"wa_id"`
}

type SentMessage struct {
	ID string `json:"id"`
}

func (wa *Whatsapp) sendMessage(m *message.Message) (*SendMessageResponse, error) {
	url, err := wa.urlFor([]string{wa.PhoneNumberId, "messages"}, nil)
	if err != nil {
		return nil, err
	}
	resp := &SendMessageResponse{}
	errBody := &util.ErrorBody{}
	_, err = wa.client.Do("POST", url, util.CONTENT_TYPE_JSON, m, resp, errBody)
	if err != nil {
		return nil, err
	}
	if errBody.Valid() {
		return nil, fmt.Errorf("err response while POST:%v\n%v", url, errBody)
	}
	return resp, nil
}

// returns a channel which delivers events one after another, use the result in a for loop,
// if you are listening for user specific events for a user those events will not be delivered here
// you must call StopListeningToEvents() to properly dispose Whatsapp instance
//
//	 allEvents:=wa.ListenToEvents()
//		for {
//			newMsg := <-allEvents
//			fmt.Println("whatsapp business delivered us a message", msg)
//		}
func (wa *Whatsapp) ListenToEvents() chan *event.WhatsAppBusinessAccount {
	if wa.nonUserMessageChannel == nil {
		wa.nonUserMessageChannel = make(chan *event.WhatsAppBusinessAccount)
	}
	return wa.nonUserMessageChannel
}

// stops any message delivery to the channel you got through ListenToEvents()
func (wa *Whatsapp) StopListeningToEvents() {
	close(wa.nonUserMessageChannel)
	wa.nonUserMessageChannel = nil
}
