package event

type Conversation struct {
	ID                  string `json:"id"`
	ExpirationTimestamp string `json:"expiration_timestamp"`
	Origin              struct {
		Type string `json:"type"`
	} `json:"origin"`
}

type Pricing struct {
	PricingModel string `json:"pricing_model"`
	Billable     bool   `json:"billable"`
	Category     string `json:"category"`
}

type Status struct {
	ID           string       `json:"id"`
	RecipientID  string       `json:"recipient_id"`
	Status       string       `json:"status"`
	Timestamp    string       `json:"timestamp"`
	Conversation Conversation `json:"conversation"`
	Pricing      Pricing      `json:"pricing"`
}

type StatusDetails struct {
	// wam id of the msg
	ID string `json:"id"`
	// phone number of the end user
	RecipientID string `json:"recipient_id"`
	// sent,delivered etc
	Status string `json:"status"`
}

// returns if the event is Status event like 'delivered','sent' etc
func (wba *WhatsAppBusinessAccount) Status() (*StatusDetails, bool) {
	value := wba.value()
	if value == nil || len(value.Statuses) == 0 {
		return nil, false
	}
	return &StatusDetails{
		ID:          value.Statuses[0].ID,
		Status:      value.Statuses[0].Status,
		RecipientID: value.Statuses[0].RecipientID,
	}, true
}
