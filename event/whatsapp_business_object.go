package event

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type MessageType string

type Message struct {
	From        string       `json:"from"`
	ID          string       `json:"id"`
	Timestamp   string       `json:"timestamp"`
	Text        TextMessage  `json:"text"`
	Interactive Interactive  `json:"interactive"`
	Image       Image        `json:"image"`
	Audio       Audio        `json:"audio"`
	Video       Video        `json:"video"`
	Document    Document     `json:"document"`
	Sticker     Sticker      `json:"sticker"`
	Location    Location     `json:"location"`
	Contacts    Contacts     `json:"contacts"`
	System      NumberChange `json:"system"`
	Type        MessageType  `json:"type"`
}

func (wba *WhatsAppBusinessAccount) value() *Value {
	if len(wba.Entry) == 0 || len(wba.Entry[0].Changes) == 0 {
		return nil
	}
	return &wba.Entry[0].Changes[0].Value
}

func (wc *Sender) WaIDUint() (uint64, error) {
	ui, err := strconv.ParseUint(wc.WaID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("err while strconv.ParseUint: %v", err)
	}
	return ui, nil
}

type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type Value struct {
	MessagingProduct string    `json:"messaging_product"`
	Metadata         Metadata  `json:"metadata"`
	Contacts         []Sender  `json:"contacts"`
	Messages         []Message `json:"messages"`
	Statuses         []Status  `json:"statuses"`
}

type Change struct {
	Value Value  `json:"value"`
	Field string `json:"field"`
}

type Entry struct {
	ID      string   `json:"id"`
	Changes []Change `json:"changes"`
}

type WhatsAppBusinessAccount struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

func NewWhatsappBusinessAccountFromPayload(payload string) (*WhatsAppBusinessAccount, error) {
	wba := new(WhatsAppBusinessAccount)
	err := json.Unmarshal([]byte(payload), wba)
	if err != nil {
		return nil, fmt.Errorf("error while unamarshalling payload %v: %v", payload, err)
	}
	if wba.Object == "" {
		return nil, nil
	}
	return wba, nil
}
