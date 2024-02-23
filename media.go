package whatsapp

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"

	"github.com/ondbyte/whatsapp/message"
	"github.com/ondbyte/whatsapp/util"
)

// creates a new document on whatsapp server
// which you can share with any number of users by persisting the returned object in the database may be.
// fileName is used to identify the file format
func (wa *Whatsapp) NewDocument(fileName string, data []byte) (*message.Document, error) {
	a, err := uploadMedia[message.Document](wa, fileName, documentLimitations, data)
	if err != nil {
		return nil, err
	}
	return a.(*message.Document), nil
}

// creates a new audio on whatsapp server
// which you can share with any number of users by persisting the returned object in the database may be.
// fileName is used to identify the file format
func (wa *Whatsapp) NewAudio(fileName string, data []byte) (*message.Audio, error) {
	a, err := uploadMedia[message.Audio](wa, fileName, audioLimitations, data)
	if err != nil {
		return nil, err
	}
	return a.(*message.Audio), nil
}

// creates a new image on whatsapp server
// which you can share with any number of users by persisting the returned object in the database may be.
// fileName is used to identify the file format
func (wa *Whatsapp) NewImage(fileName string, data []byte) (*message.Image, error) {
	a, err := uploadMedia[message.Image](wa, fileName, imageLimitations, data)
	if err != nil {
		return nil, err
	}
	return a.(*message.Image), nil
}

// creates a new video on whatsapp server
// which you can share with any number of users by persisting the returned object in the database may be.
// fileName is used to identify the file format
func (wa *Whatsapp) NewVideo(fileName string, data []byte) (*message.Video, error) {
	a, err := uploadMedia[message.Video](wa, fileName, videoLimitations, data)
	if err != nil {
		return nil, err
	}
	return a.(*message.Video), nil
}

// creates a new sticker on whatsapp server
// which you can share with any number of users by persisting the returned object in the database may be.
// fileName is used to identify the file format
func (wa *Whatsapp) NewSticker(fileName string, data []byte) (*message.Sticker, error) {
	a, err := uploadMedia[message.Sticker](wa, fileName, stickerLimitations, data)
	if err != nil {
		return nil, err
	}
	return a.(*message.Sticker), nil
}

type uploadMediaResponse struct {
	Id string `json:"id"`
}

// Valid implements util.Body.
func (u *uploadMediaResponse) Valid() bool {
	return u.Id != ""
}

func uploadMedia[T any](wa *Whatsapp, fileName string, mediaType mediaLimitations[T], mediaData []byte) (any any, err error) {
	extension := filepath.Ext(fileName)
	mimeType, ok := mediaType.mimeTypes[extension]
	if !ok {
		return nil, fmt.Errorf("whatsapp doesn't support media type %v to upload", extension)
	}
	url, err := wa.urlFor(
		[]string{wa.PhoneNumberId, "media"},
		nil,
	)
	if err != nil {
		return nil, err
	}
	body := new(bytes.Buffer)
	formWriter := util.NewCustomFormWriter(body)
	partWriter, err := formWriter.CreateFormFileWithCustomMediaType("file", fileName, mimeType)
	if err != nil {
		return nil, fmt.Errorf("error while util.CustomForm.CreateFormFileWithCustomMediaType: %v", err)
	}
	_, err = io.Copy(partWriter, bytes.NewReader(mediaData))
	if err != nil {
		return nil, fmt.Errorf("error while io.Copy: %v", err)
	}
	// add some fields required by whatsapp
	err = formWriter.WriteField("messaging_product", "whatsapp")
	if err != nil {
		return nil, fmt.Errorf("error while util.CustomForm.WriteField: %v", err)
	}
	err = formWriter.WriteField("type", mimeType)
	if err != nil {
		return nil, fmt.Errorf("error while util.CustomForm.WriteField: %v", err)
	}
	err = formWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("error while closing formWriter: %v", err)
	}
	resp, errResp := &uploadMediaResponse{}, &util.ErrorBody{}
	_, err = wa.client.Do("POST", url, formWriter.FormDataContentType(), body, resp, errResp)
	if err != nil {
		return nil, err
	}
	if errResp.Valid() {
		return nil, fmt.Errorf("got error response from endpoint %v, \nerr response is: %v", url, errResp)
	}
	return mediaType.wrap(resp.Id), nil
}
