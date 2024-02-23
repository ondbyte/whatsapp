# a go-lang package 
to write whatsapp conversation bots, based on the official WABA api

## what works

Sending/Receiving of messages 

| message type        | Sending | Receiving |
|----------------|---------|-----------|
| Text           |  yes✅   |   yes✅   |
| Location       |  yes✅   |   yes✅   |
| Contact        |  yes✅   |   yes✅   |
| Audio          |  fails❌ |   yes✅   |
| Image          |  yes✅   |   yes✅   |
| Video          |  yes✅   |   yes✅   |
| Document       |  yes✅   |   yes✅   |
| Sticker        |  yes✅   |   yes✅   |

Signature verification of the incoming events/messages on the webhook: ✅

Uploading of the media to the whatsapp servers:
- Image: ✅
- Video: ✅
- Audio: ✅
- Document: ✅
- Sticker: ✅



### example

```go
wa, _ := whatsapp.New(whatsapp.Config{/* set config */})
// mount wa.WhatsappEvents in your http server and configure that endpoint in the facebook dashboard and enable messages

events:=whatsapp.ListenToEvents()
for {
	incomingMsg := <-events
	// process the incomingMsg or start a conversation with the user
	conversation := whatsapp.StartConversation(incomingMsg.WaId())
    haveAChat(whatsapp,conversation)
}

func haveAChat(whatsapp whatsapp.Whatsapp, firstMsg event.WhatsAppBusinessAccount) {
	tm, isTextMsg := firstMsg.TextMessage()
	if !isTextMsg {
		return
	}
	conversation := whatsapp.StartConversation(firstMsg.WaId())
	defer conversation.End()
	switch tm.Body {
	case "Foo":
		conversation.SendText("Bar", false)
		var txtMsg *event.TextMessageDetails
		var ok bool
		for {
			// wait only for text messages
			event := conversation.NextEvent()
			txtMsg, ok = event.TextMessage()
			if ok {
				break
			}
		}
		if txtMsg.Body != "No bar" {
			return
		}
		conversation.SendText("Here is sticker of a bar", false)
		// you mustn't call NewSticker for everytime you require the same sticker, reuse the returned sticker
		sticker, _ := whatsapp.NewSticker("bar.webp", []byte{ /* bytes of the sticker */ })
		conversation.SendSticker(sticker)
		conversation.SendText("this conversation is finished,thank you", false)
	}
}
```

### find a full example in whatsapp_test.go