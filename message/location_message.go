package message

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
}

func NewLocation(longitude, latitude float64, name, address string) *Location {
	return &Location{
		Longitude: longitude,
		Latitude:  latitude,
		Name:      name,
		Address:   address,
	}
}

var LOCATION MessageType = "location"

func NewLocationMessage(
	to string,
	location *Location,
) *Message {
	return &Message{
		MessagingProduct: WHATSAPP,
		RecipientType:    INDIVIDUAL,
		To:               to,
		Type:             LOCATION,
		Location:         location,
	}
}
