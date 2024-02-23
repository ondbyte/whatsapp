package event

type Address struct {
	City        string `json:"city"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	Street      string `json:"street"`
	Type        string `json:"type"`
	Zip         string `json:"zip"`
}

type Email struct {
	Email string `json:"email"`
	Type  string `json:"type"`
}

type IM struct {
	Service string `json:"service"`
	UserID  string `json:"user_id"`
}

type Name struct {
	FormattedName string `json:"formatted_name"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	MiddleName    string `json:"middle_name"`
	Suffix        string `json:"suffix"`
	Prefix        string `json:"prefix"`
}

type Org struct {
	Company    string `json:"company"`
	Department string `json:"department"`
	Title      string `json:"title"`
}

type Phone struct {
	Phone string `json:"phone"`
	Type  string `json:"type"`
	WaID  string `json:"wa_id,omitempty"`
}

type URL struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

type Contact struct {
	Addresses []Address `json:"addresses"`
	Birthday  string    `json:"birthday"`
	Emails    []Email   `json:"emails"`
	IMS       []IM      `json:"ims"`
	Name      Name      `json:"name"`
	Org       Org       `json:"org"`
	Phones    []Phone   `json:"phones"`
	URLs      []URL     `json:"urls"`
}

type Contacts []Contact

var CONTACTS MessageType = "contacts"

// returns if the event is contacts event/msg i e a user sent contacts to you
func (w *WhatsAppBusinessAccount) Contacts() (Contacts, bool) {
	value := w.value()
	if value == nil || len(value.Messages) == 0 || value.Messages[0].Type != CONTACTS {
		return nil, false
	}
	return value.Messages[0].Contacts, true
}
