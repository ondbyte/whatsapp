package whatsapp

import (
	"crypto/hmac"
	"encoding/hex"
	"fmt"
)

// this verifies incoming events on the webhook are actually from facebook
// by using the app secret you passed as configuration
func (w *Whatsapp) verifySignAgainstFacebookSecret(shaSign, data string) error {
	w.facebookAppSecretHasher.Reset()
	shaSignBytes, err := hex.DecodeString(shaSign)
	if err != nil {
		return fmt.Errorf("error while hex.DecodeString: %v", err)
	}
	_, err = w.facebookAppSecretHasher.Write([]byte(data))
	if err != nil {
		return fmt.Errorf("error while writing data to hasher : %v", err)
	}
	sign := w.facebookAppSecretHasher.Sum([]byte{})
	equal := hmac.Equal(shaSignBytes, sign)
	if !equal {
		return fmt.Errorf("invalid data")
	}
	return nil
}
