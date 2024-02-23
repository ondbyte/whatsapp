package event

var LOCATION MessageType = "location"

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
}

// returns if the event is a Location event/msg
func (wba *WhatsAppBusinessAccount) Location() (*Location, bool) {
	value := wba.value()
	if value == nil || len(value.Messages) == 0 || value.Messages[0].Type != LOCATION {
		return nil, false
	}
	return &value.Messages[0].Location, true
}
