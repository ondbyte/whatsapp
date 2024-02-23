package whatsapp

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ondbyte/whatsapp/event"
)

// this is where Whatsapp receives events/incoming messages from users
// wherever you mount this should be added as webhook for whatsapp business.
//
// you can take advantage of ngrok package available alongside this package to
// test this package like the following example.
//
// ```go
//
// ```
func (wa *Whatsapp) WhatsappEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		q := r.URL.Query()
		challenge := q.Get("hub.challenge")
		verifyToken := q.Get("hub.verify_token")
		err := wa.verifyWhatsappToken(verifyToken)
		if err != nil {
			fmt.Printf("you tried to add me as a webhook but its failed because err: %v\n", err)
			return
		}
		w.Write([]byte(challenge))
		fmt.Println("you have successfully added the webhook on facebook dashboard, listening for events/messages")
		return
	}
	if r.Method == "POST" {
		defer w.WriteHeader(http.StatusOK)
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println(fmt.Errorf("err while reading body: %v", err))
			return
		}
		shaSign := r.Header.Get("X-Hub-Signature-256")
		shaSign = strings.Replace(shaSign, "sha256=", "", 1)
		payload := string(bytes)
		err = wa.verifySignAgainstFacebookSecret(shaSign, payload)
		if err != nil {
			fmt.Println(fmt.Errorf("Whatsapp.VerifySignAgainstWhatsappSecret returned err: %v", err))
			return
		}
		wba, err := event.NewWhatsappBusinessAccountFromPayload(payload)
		if err != nil {
			fmt.Println(fmt.Errorf("events.NewWhatsappBusinessAccountFromPayload returned err: %v", err))
			return
		}
		if wba == nil {
			fmt.Printf("unsupported event rcd : \n%v\n", payload)
			return
		}

		if wba == nil {
			fmt.Printf("un-handled event rcd : \n%v\n", wba)
			return
		}
		waId := wba.WaId()
		if waId == "" {
			// failed to get waid ie the number from which we rcd the wba
			// this shouldnt be the case
			fmt.Println("unable to get wa_id for events.WhatsAppBusinessAccount: payload as follows: \n", payload)
		}
		var convo *Conversation

		if waId != "" {
			wa.userMessageChannelsLock.Lock()
			convo = wa.conversations[waId]
			wa.userMessageChannelsLock.Unlock()
		}
		if convo != nil {
			convo.receive <- wba
		} else if wa.nonUserMessageChannel != nil {
			wa.nonUserMessageChannel <- wba
		} else {
			fmt.Println("no listeners found, discarding event \n", payload)
		}
	}
}
