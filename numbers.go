package whatsapp

import (
	"fmt"
	"strings"

	"github.com/ondbyte/whatsapp/util"
)

type Number struct {
	ID                 string `json:"id"`
	DisplayPhoneNumber string `json:"display_phone_number"`
	VerifiedName       string `json:"verified_name"`
	QualityRating      string `json:"quality_rating"`
}

type getNumbersResponse struct {
	Data []*Number `json:"data"`
}

// Valid implements util.Body.
func (g *getNumbersResponse) Valid() bool {
	return len(g.Data) > 0
}

func (w *Whatsapp) GetNumbers() (numbers map[string]*Number, err error) {
	bi := &getNumbersResponse{
		Data: []*Number{},
	}
	url, err := w.urlFor([]string{w.businessAccountId, "phone_numbers"}, nil)
	if err != nil {
		return nil, err
	}
	errBody := &util.ErrorBody{}
	_, err = w.client.Get(url, util.CONTENT_TYPE_JSON, nil, bi, errBody)
	if err != nil {
		return nil, fmt.Errorf("err while w.client.HttpGet: URL: %v \n%v", url, err)
	}
	if errBody.Valid() {
		return nil, fmt.Errorf("got err response : URL: %v\n%v", url, errBody)
	}
	numbers = map[string]*Number{}
	for _, v := range bi.Data {
		number := strings.ReplaceAll(strings.ReplaceAll(v.DisplayPhoneNumber, "-", ""), " ", "")
		numbers[number] = v
	}
	return
}

// you must call this before calling any other methods/apis,
// fetches all the numbers available for the Whatsapp.BusinessAccountId,
// sets the one with the number's id as Whatsapp.PhoneNumberId which will be used when we use the message api
func (w *Whatsapp) useNumber(number string) error {
	numbers, err := w.GetNumbers()
	if err != nil {
		return fmt.Errorf("error while w.GetNumbers %v", err)
	}
	number = strings.ReplaceAll(strings.ReplaceAll(number, "-", ""), " ", "")
	num, ok := numbers[number]
	if !ok {
		return fmt.Errorf("number %v isn't available with WABA", number)
	}
	w.PhoneNumberId = num.ID
	fmt.Printf("using phone number %v with id %v\n", number, num.ID)
	return nil
}
