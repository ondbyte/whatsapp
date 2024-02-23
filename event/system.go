package event

var SYSTEM MessageType = "system"

type NumberChange struct {
	Body    string `json:"body"`
	NewWAID string `json:"new_wa_id"`
	Type    string `json:"type"`
}

type NumberChangeDetails struct {
	Old string `json:"old"`
	New string `json:"new"`
}

// returns non nil and true if the this WBA is a system message (for example when a number of user is changed)
func (wba *WhatsAppBusinessAccount) NumberChange() (*NumberChangeDetails, bool) {
	value := wba.value()
	if value == nil || len(value.Messages) == 0 || value.Messages[0].Type != SYSTEM {
		return nil, false
	}
	return &NumberChangeDetails{
		Old: value.Messages[0].From,
		New: value.Messages[0].System.NewWAID,
	}, true
}
