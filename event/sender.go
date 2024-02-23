package event

type Profile struct {
	Name string `json:"name"`
}

type Sender struct {
	Profile Profile `json:"profile"`
	WaID    string  `json:"wa_id"`
}

// returns a users phone number if the event has one
func (wba *WhatsAppBusinessAccount) WaId() string {
	if len(wba.Entry) == 0 || len(wba.Entry[0].Changes) == 0 {
		return ""
	}
	value := wba.Entry[0].Changes[0].Value
	if len(value.Contacts) > 0 {
		return value.Contacts[0].WaID
	}
	if len(value.Statuses) > 0 {
		return value.Statuses[0].RecipientID
	}
	if len(value.Messages) > 0 {
		return value.Messages[0].From
	}
	return ""
}
