package whatsapp_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	ngrok "github.com/ondbyte/ngrok-server"
	"github.com/ondbyte/whatsapp"
	"github.com/ondbyte/whatsapp/event"
	"github.com/ondbyte/whatsapp/message"
)

type TestConfig struct {
	WhatsappCfg whatsapp.Config `json:"whatsapp_config"`
	NgrokCfg    struct {
		AuthToken string `json:"authtoken"`
		Domain    string `json:"domain"`
	} `json:"ngrok_config"`
	TestNumberToHaveTheConversationWith string `json:"test_number_to_have_the_conversation_with"`
}

func samplePngImage() (data []byte, contentType string, err error) {
	return sampleMedia("https://samplelib.com/lib/preview/png/sample-clouds2-400x300.png")
}

func sampleMp4Video() (data []byte, contentType string, err error) {
	return sampleMedia("https://samplelib.com/lib/preview/mp4/sample-5s.mp4")
}

func sampleOggAudio() (data []byte, contentType string, err error) {
	return sampleMedia("https://onlinetestcase.com/wp-content/uploads/2023/06/500-KB-OGG.ogg")
}

func samplePdfDocument() (data []byte, contentType string, err error) {
	return sampleMedia("https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf")
}

func sampleMedia(url string) (data []byte, contentType string, err error) {
	// Make a GET request to the image URL
	response, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	// Check if the request was successful (status code 200)
	if response.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to download image, status code: %d", response.StatusCode)
	}

	// Read the image content into a byte slice
	imageBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}
	contentType = response.Header.Get("Content-Type")
	return imageBytes, contentType, nil
}

func eventToMsgContacts(contacts event.Contacts) (contacts2 message.Contacts) {
	for _, c := range contacts {
		addresses2 := []*message.Address{}
		for _, addr := range c.Addresses {
			addresses2 = append(addresses2, message.NewAddress(addr.Street, addr.City, addr.State, addr.Zip, addr.Country, addr.CountryCode, addr.Type))
		}
		emails2 := []*message.Email{}
		for _, email := range c.Emails {
			emails2 = append(emails2, message.NewEmail(email.Email, email.Type))
		}
		phone2 := []*message.Phone{}
		for _, phone := range c.Phones {
			phone2 = append(phone2, message.NewPhone(phone.Phone, phone.Type, phone.WaID))
		}
		url2 := []*message.URL{}
		for _, url := range c.URLs {
			url2 = append(url2, message.NewURL(url.URL, url.Type))
		}
		c2 := message.NewContact(
			addresses2,
			c.Birthday,
			emails2,
			message.NewName(
				c.Name.FormattedName,
				c.Name.FirstName,
				c.Name.LastName,
				c.Name.MiddleName,
				c.Name.Suffix,
				c.Name.Prefix,
			),
			message.NewOrg(
				c.Org.Company,
				c.Org.Department,
				c.Org.Title,
			),
			phone2,
			url2,
		)
		contacts2 = append(contacts2, c2)
	}
	return
}

func loadCfg() (*TestConfig, error) {
	testCfgPath := "./test_cfg.json"
	testCfgFileData, err := os.ReadFile(testCfgPath)
	testCfg := &TestConfig{}
	err = json.Unmarshal(testCfgFileData, testCfg)
	return testCfg, err
}

func TestCTA(t *testing.T) {
	assert := assert.New(t)
	cfg, err := loadCfg()
	assert.NoError(err, `err while loadCfg`)
	wa, err := whatsapp.New(cfg.WhatsappCfg)
	if !assert.NoError(err) {
		return
	}
	if !assert.NotEmpty(wa.PhoneNumberId) {
		return
	}
	convo := wa.StartConversation(cfg.TestNumberToHaveTheConversationWith)

	convo.NextEvent()

	_, err = convo.SendInteractive(
		message.NewInteractiveWithCtaButton(
			message.NewTextHeader(`Hello how are you`),
			message.NewBody(`Yadhunandan haha haha`),
			nil,
			message.NewActionWithCTAButton(`Yeah`, `https://stackoverflow.com/questions/11100155/invalid-oauth-access-token-when-using-an-application-access-token-android`),
		),
	)
	assert.NoError(err, `err while sending msg`, err)
}
func TestWhatsapp(t *testing.T) {
	assert := assert.New(t)
	cfg, err := loadCfg()
	assert.NoError(err, `err while loadCfg`)
	wa, err := whatsapp.New(cfg.WhatsappCfg)
	if !assert.NoError(err) {
		return
	}
	if !assert.NotEmpty(wa.PhoneNumberId) {
		return
	}
	bi, err := wa.GetBusinessInfo()
	if assert.NoError(err) && assert.NotNil(bi) && assert.NotEmpty(bi.ID) {
		fmt.Println(bi)
	}
	numbers, err := wa.GetNumbers()
	if !assert.NoError(err) || !assert.NotNil(numbers) {
		return
	}

	handler := http.DefaultServeMux
	listeningPath := "/whatsapp_events_2"
	handler.HandleFunc(listeningPath, wa.WhatsappEvents)
	// to listen to incoming msg events start a server
	server, err := ngrok.NewWithDomain(context.TODO(), cfg.NgrokCfg.AuthToken, cfg.NgrokCfg.Domain)
	if !assert.NoError(err, "error while starting a ngrok server ") {
		return
	}
	go func() {
		fmt.Printf("test whatsapp business events webhook: %v%v\n", server.Url(), listeningPath)
		server.Serve(":8008", handler)
	}()

	usrNumber := cfg.TestNumberToHaveTheConversationWith
	conversation := wa.StartConversation(usrNumber)
	_, err = conversation.SendText(`Hi..`, false)
	assert.NoError(err, `expected no error while sending the first msg`)
	fmt.Println("expecting the message to be a 'test' to begin the test, send 'test' to ", cfg.WhatsappCfg.Number)
	for {
		evt := <-conversation.NextEvent()
		if txtMsg, ok := evt.TextMessage(); ok && strings.ToLower(txtMsg.Body) == "test" {
			// first text msg from the user is test
			// so you should begin the test
			break
		}
	}
	_, err = conversation.SendText("hi, begining of a test transaction", false)
	if !assert.NoError(err, "err while sending a message", err) {
		return
	}
	// test location
	_, err = conversation.SendText("please respond correctly so the test case passes, send me a location", false)
	if !assert.NoError(err, "err while sending a message", err) {
		return
	}

	for {
		reply := <-conversation.NextEvent()
		if _, isStatusUpdate := reply.Status(); isStatusUpdate {
			continue
		}

		location, isLocation := reply.Location()
		if isLocation {
			_, err = conversation.SendText("I recieved the following location from you, thanks", false)
			if !assert.NoError(err, "err while sending a message", err) {
				return
			}
			_, err = conversation.SendLocation(
				message.NewLocation(location.Longitude, location.Latitude, location.Name, location.Address),
			)
			if !assert.NoError(err, "err while sending a message", err) {
				return
			}
			break
		}
		_, err = conversation.SendText("please respond correctly so the test case passes, send me a location", false)
		if !assert.NoError(err, "err while sending a message", err) {
			return
		}
	}
	// test contact
	_, err = conversation.SendText("please respond correctly so the test case passes, send me a contact", false)
	if !assert.NoError(err, "err while sending a message", err) {
		return
	}
	for {
		reply := <-conversation.NextEvent()
		if _, isStatusUpdate := reply.Status(); isStatusUpdate {
			continue
		}
		contacts, isContacts := reply.Contacts()
		if isContacts {
			_, err := conversation.SendText("I received the following contacts", false)
			if !assert.NoError(err, "err while sending message", err) {
				return
			}
			contacts2 := eventToMsgContacts(contacts)
			_, err = conversation.SendContacts(contacts2)
			if !assert.NoError(err, "err while sending contacts", err) {
				return
			}
			break
		}
		_, err = conversation.SendText("please respond correctly so the test case passes, send me a contact", false)
		if !assert.NoError(err, "err while sending a message", err) {
			return
		}
	}
	// test image
	_, err = conversation.SendText("please respond correctly so the test case passes, send me a image", false)
	if !assert.NoError(err, "err while sending a message", err) {
		return
	}
	for {
		reply := <-conversation.NextEvent()
		if _, isStatusUpdate := reply.Status(); isStatusUpdate {
			continue
		}
		image, isImage := reply.Image()
		if isImage {
			imageBytes, _, err := samplePngImage()
			if !assert.NoError(err, "err while getting sample img") {
				return
			}
			// upload the image to whatsapp server, once uploaded you can send the same image to multiple users
			yaduJpg, err := wa.NewImage("yadu.png", imageBytes)
			if !assert.NoError(err, "err while uploading image to whatsapp server", err) {
				return
			}
			_, err = conversation.SendText("I received the following image, we have included one from our side as well", false)
			if !assert.NoError(err, "err while sending message", err) {
				return
			}
			_, err = conversation.SendImage(message.NewImage(image.Id))
			if !assert.NoError(err, "err while sending image", err) {
				return
			}
			_, err = conversation.SendImage(yaduJpg)
			if !assert.NoError(err, "err while sending image", err) {
				return
			}
			break
		}
		_, err = conversation.SendText("please respond correctly so the test case passes, send me a image", false)
		if !assert.NoError(err, "err while sending a message", err) {
			return
		}
	}

	// test video
	_, err = conversation.SendText("please respond correctly so the test case passes, send me a video", false)
	if !assert.NoError(err, "err while sending a message", err) {
		return
	}
	for {
		reply := <-conversation.NextEvent()
		if _, isStatusUpdate := reply.Status(); isStatusUpdate {
			continue
		}
		video, isVideo := reply.Video()
		if isVideo {
			// upload the video to whatsapp server, once uploaded you can send the same video to multiple users
			sampleMp4VidBytes, _, err := sampleMp4Video()
			if !assert.NoError(err, "err getting sample mp4 vid") {
				return
			}
			yaduMp4, err := wa.NewVideo("yadu.mp4", sampleMp4VidBytes)
			if !assert.NoError(err, "err while uploading video") {
				return
			}
			_, err = conversation.SendText("I received the following video, we have included one from our side as well", false)
			if !assert.NoError(err, "err while sending message", err) {
				return
			}
			_, err = conversation.SendVideo(message.NewVideo(video.Id))
			if !assert.NoError(err, "err while sending video", err) {
				return
			}
			_, err = conversation.SendVideo(yaduMp4)
			if !assert.NoError(err, "err while sending video", err) {
				return
			}
			break
		}
		_, err = conversation.SendText("please respond correctly so the test case passes, send me a video", false)
		if !assert.NoError(err, "err while sending a message", err) {
			return
		}
	}

	_, err = conversation.SendText("please respond correctly so the test case passes, send me a audio", false)
	if !assert.NoError(err, "err while sending a message", err) {
		return
	}
	for {
		reply := <-conversation.NextEvent()
		if _, isStatusUpdate := reply.Status(); isStatusUpdate {
			continue
		}
		_, isAudio := reply.Audio()
		if isAudio {
			// upload the audio to whatsapp server, once uploaded you can send the same audio to multiple users
			sampleOggAudio, _, err := sampleOggAudio()
			if !assert.NoError(err, "err getting sample ogg audio") {
				return
			}
			yaduOgg, err := wa.NewAudio("yadu.ogg", sampleOggAudio)
			if !assert.NoError(err, "err while uploading audio") {
				return
			}
			_, err = conversation.SendText("I received the following audio, we have included one from our side as well", false)
			if !assert.NoError(err, "err while sending message", err) {
				return
			}
			/* _, err = conversation.SendAudio(message.NewAudio(audio.Id))
			if !assert.NoError(err, "err while sending audio", err) {
				return
			} */
			_, err = conversation.SendAudio(yaduOgg)
			if !assert.NoError(err, "err while sending audio", err) {
				return
			}
			break
		}
		_, err = conversation.SendText("please respond correctly so the test case passes, send me a audio", false)
		if !assert.NoError(err, "err while sending a message", err) {
			return
		}
	}
}
