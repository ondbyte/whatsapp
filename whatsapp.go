package whatsapp

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"hash"
	"net/url"
	"sync"

	event "github.com/ondbyte/whatsapp/event"
	"github.com/ondbyte/whatsapp/util"
)

var GraphApiVersion18_0 = "v18.0"

type TokenVerifier func(string) error

type Whatsapp struct {
	businessAccountId, PhoneNumberId string
	client                           *util.AuthorizedHttpClient
	eventToken                       string
	facebookAppSecretHasher          hash.Hash
	userMessageChannelsLock          *sync.Mutex
	conversations                    map[string]*Conversation
	nonUserMessageChannel            chan *event.WhatsAppBusinessAccount
	graphApiUrl                      string
}

type Config struct {
	// go here https://developers.facebook.com/apps,
	// click on your app and on the dashboard of your app select >> App settings >> Basic to find your app secret
	FacebookAppSecret string `json:"facebook_app_secret"`
	// go here https://developers.facebook.com/apps,
	// then WebHooks >> select whatsapp business account >> subscribe to this object >>
	// add your url where you will be listening for incoming messages (where this whatsapp listener is running,see events_handler.go) and
	// enter verify token which should be same as this
	VerifyToken string `json:"verify_token"`
	// you can (possibly) find it in whatsapp account / whatsapp manager at
	// https://business.facebook.com/latest/settings/whatsapp_account
	WhatsappBusinessAccountId string `json:"whatsapp_business_account_id"`
	// temp token can be obtained from here https://developers.facebook.com/tools/explorer/
	// while in production you must use any valid token defined here https://developers.facebook.com/docs/facebook-login/guides/access-tokens/
	AccessToken string `json:"access_token"`
	// this number should be added in the whatsapp manager, this number will be used by the whatsapp
	// instance to send and receive messages
	Number string `json:"number"`
}

func (c Config) Valid() bool {
	return c.VerifyToken != "" && c.AccessToken != "" && c.WhatsappBusinessAccountId != "" && c.Number != "" && c.FacebookAppSecret != ""
}

// returns new Whatsapp business api
// error if unable to get the whatsapp business details
func New(cfg Config) (*Whatsapp, error) {
	if !cfg.Valid() {
		return nil, fmt.Errorf("whatsapp configuration is not valid some fields might be missing")
	}
	w := &Whatsapp{
		businessAccountId:       cfg.WhatsappBusinessAccountId,
		client:                  util.NewAuthorizedHttpClient(cfg.AccessToken),
		eventToken:              cfg.VerifyToken,
		facebookAppSecretHasher: hmac.New(sha256.New, []byte(cfg.FacebookAppSecret)),
		conversations:           make(map[string]*Conversation),
		userMessageChannelsLock: &sync.Mutex{},
		graphApiUrl:             GraphUrlv19,
	}
	err := w.useNumber(cfg.Number)
	if err != nil {
		return nil, fmt.Errorf("unable to use number %v with this whatsapp instance, make sure its available in the WABA settings: %v", err, cfg.Number)
	}
	return w, nil
}
func (wa *Whatsapp) verifyWhatsappToken(s string) error {
	if s == wa.eventToken {
		return nil
	}
	return fmt.Errorf("invalid token")
}

var GraphUrlv19 = "https://graph.facebook.com/v18.0"

func (w *Whatsapp) SetApiUrl(url string) {
	w.graphApiUrl = url
}

func (w *Whatsapp) urlFor(pathFragments []string, params map[string]string) (string, error) {
	host := w.graphApiUrl
	for _, fragment := range pathFragments {
		host += "/" + fragment
	}
	host += "?"
	for k, v := range params {
		host += k + "=" + v + "&"
	}
	u, err := url.Parse(host)
	if err != nil {
		return "", fmt.Errorf("err while parsing '%v' as url", host)
	}
	return u.String(), nil
}

func (w *Whatsapp) GetBusinessInfo() (bi *BusinessInfo, err error) {
	bi = &BusinessInfo{}
	url, err := w.urlFor([]string{w.businessAccountId}, nil)
	if err != nil {
		return nil, err
	}
	errBody := &util.ErrorBody{}
	_, err = w.client.Get(url, util.CONTENT_TYPE_JSON, nil, bi, errBody)
	if err != nil {
		return nil, fmt.Errorf("err while w.client.HttpGet: %v", err)
	}
	if errBody.Valid() {
		return nil, fmt.Errorf("gor err response \n%v", errBody)
	}
	return bi, nil
}
