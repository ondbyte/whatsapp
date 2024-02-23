package message

type Address struct {
	Street      string `json:"street"`
	City        string `json:"city"`
	State       string `json:"state"`
	Zip         string `json:"zip"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Type        string `json:"type"`
}

func NewAddress(street, city, state, zip, country, countryCode, addressType string) *Address {
	return &Address{
		Street:      street,
		City:        city,
		State:       state,
		Zip:         zip,
		Country:     country,
		CountryCode: countryCode,
		Type:        addressType,
	}
}

type Email struct {
	Email string `json:"email"`
	Type  string `json:"type"`
}

func NewEmail(eMail, emailType string) *Email {
	return &Email{
		Email: eMail,
		Type:  emailType,
	}
}

type Name struct {
	FormattedName string `json:"formatted_name"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	MiddleName    string `json:"middle_name"`
	Suffix        string `json:"suffix"`
	Prefix        string `json:"prefix"`
}

func NewName(formattedName, firstName, lastName, middleName, suffix, prefix string) *Name {
	return &Name{
		FormattedName: formattedName,
		FirstName:     firstName,
		LastName:      lastName,
		MiddleName:    middleName,
		Suffix:        suffix,
		Prefix:        prefix,
	}
}

type Org struct {
	Company    string `json:"company"`
	Department string `json:"department"`
	Title      string `json:"title"`
}

func NewOrg(company, department, title string) *Org {
	return &Org{
		Company:    company,
		Department: department,
		Title:      title,
	}
}

type Phone struct {
	Phone string `json:"phone"`
	Type  string `json:"type"`
	WaID  string `json:"wa_id,omitempty"`
}

func NewPhone(_phone, phoneType, waID string) *Phone {
	return &Phone{
		Phone: _phone,
		Type:  phoneType,
		WaID:  waID,
	}
}

type URL struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

func NewURL(_url, urlType string) *URL {
	return &URL{
		URL:  _url,
		Type: urlType,
	}
}

type Contact struct {
	Addresses []*Address `json:"addresses"`
	Birthday  string     `json:"birthday"`
	Emails    []*Email   `json:"emails"`
	Name      *Name      `json:"name"`
	Org       *Org       `json:"org"`
	Phones    []*Phone   `json:"phones"`
	URLs      []*URL     `json:"urls"`
}

type Contacts []*Contact

func NewContact(addresses []*Address, birthday string, emails []*Email, name *Name, org *Org, phones []*Phone, urls []*URL) *Contact {
	return &Contact{
		Addresses: addresses,
		Birthday:  birthday,
		Emails:    emails,
		Name:      name,
		Org:       org,
		Phones:    phones,
		URLs:      urls,
	}
}

func NewContacts(contacts ...*Contact) Contacts {
	return contacts
}

var CONTACTS MessageType = "contacts"

func NewContactsMsg(to string, contacts Contacts) *Message {
	return &Message{
		MessagingProduct: WHATSAPP,
		To:               to,
		Type:             CONTACTS,
		Contacts:         contacts,
		RecipientType:    INDIVIDUAL,
	}
}
