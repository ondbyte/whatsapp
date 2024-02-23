package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"strings"
)

var DEBUG_MODE = false

func init() {
	DEBUG_MODE = !(os.Getenv("DEBUG") == "FALSE")
	if DEBUG_MODE {
		fmt.Println("DEBUG MODE ENABLED")
	} else {
		fmt.Println("DEBUG MODE DISABLED")
	}
}

type AuthorizedHttpClient struct {
	// default params we pass to the each requests
	defaultReqParams url.Values
	defaultHeaders   http.Header
	client           http.Client
}

func NewAuthorizedHttpClient(authorizationToken string) *AuthorizedHttpClient {
	a := &AuthorizedHttpClient{
		client:           http.Client{},
		defaultHeaders:   http.Header{},
		defaultReqParams: url.Values{},
	}
	a.defaultHeaders.Add("Authorization", fmt.Sprintf("Bearer %v", authorizationToken))
	return a
}

type ErrorBody struct {
	Err Error `json:"error"`
}

// returns whether this err is non nil
func (e *ErrorBody) Valid() bool {
	return e.Err.Message != ""
}

func (e *ErrorBody) Error() string {
	b, _ := json.Marshal(e)
	return fmt.Sprintf("%v", string(b))
}

type ErrorData struct {
	MessagingProduct string `json:"messaging_product"`
	Details          string `json:"details"`
}

type Error struct {
	Message        string `json:"message"`
	Type           string `json:"type"`
	Code           int64  `json:"code"`
	ErrorSubcode   int64  `json:"error_subcode"`
	ErrorUserTitle string `json:"error_user_title"`
	ErrorUserMsg   string `json:"error_user_msg"`
	FbtraceID      string `json:"fbtrace_id"`

	ErrorData ErrorData `json:"error_data"`
}

func DebugPrintF(format string, a ...any) {
	if DEBUG_MODE {
		fmt.Printf(format, a...)
	}
}

type Body interface {
	Valid() bool
}

func (fc *AuthorizedHttpClient) httpReq(method, url string, contentType string, reqBody any, respBodies ...Body) (body []byte, err error) {
	var bodyReader io.Reader
	bs, isBytes := reqBody.([]byte)
	bb, isBuffer := reqBody.(*bytes.Buffer)
	if isBuffer {
		bodyReader = bb
	} else if isBytes {
		bodyReader = bytes.NewBuffer(bs)
	} else if reqBody != nil {
		b, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("error while json.Marshal: %v", err)
		}
		s := string(b)
		bodyReader = strings.NewReader(s)
		DebugPrintF("REQUEST:\n\t%v\n\t%v\n\tBODY: %v", method, url, s)
	}
	url = fmt.Sprintf("%v?%v", url, fc.defaultReqParams.Encode())
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("err while http.NewRequest: %v", err)
	}
	req.Header.Add("Content-Type", contentType)
	for k, v := range fc.defaultHeaders {
		req.Header[k] = v
	}
	resp, err := fc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error doing a get request due to err: %v", err)
	}
	defer resp.Body.Close()

	// Read the resp body
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading the body of the get request due to err: %v", err)
	}
	for _, respBody := range respBodies {
		err = json.Unmarshal(b, respBody)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling the body to type %t due to err: %v", respBody, err)
		}
		if respBody.Valid() {
			return b, nil
		}
	}
	return b, nil
}

func (fc *AuthorizedHttpClient) Get(url string, contentType string, reqBody any, respBodies ...Body) ([]byte, error) {
	return fc.httpReq("GET", url, contentType, reqBody, respBodies...)
}

func (fc *AuthorizedHttpClient) Post(url string, contentType string, reqBody any, respBodies ...Body) ([]byte, error) {
	return fc.httpReq("POST", url, contentType, reqBody, respBodies...)
}

var (
	CONTENT_TYPE_FORM_DATA_FILE = "multipart/form-data"
	CONTENT_TYPE_JSON           = "application/json"
)

type CustomForm struct {
	multipart.Writer
}

func NewCustomFormWriter(buffer *bytes.Buffer) *CustomForm {
	return &CustomForm{
		Writer: *multipart.NewWriter(buffer),
	}
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func (cf *CustomForm) CreateFormFileWithCustomMediaType(fieldname string, filename string, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", contentType)
	return cf.CreatePart(h)
}

// does a request and unmarshalls response to the bodies
func (fc *AuthorizedHttpClient) Do(
	method string,
	url string,
	contentType string,
	body any,
	respBodies ...Body,
) (response *http.Response, err error) {
	bodyB := new(bytes.Buffer)
	b, ok := body.([]byte)
	if ok {
		bodyB = bytes.NewBuffer(b)
	} else if buf, ok := body.(*bytes.Buffer); ok {
		bodyB = buf
	} else {
		b, err = json.Marshal(body)
		bodyB = bytes.NewBuffer(b)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal while adding body to a request: %v", err)
	}
	req, err := http.NewRequest(method, url, bodyB)
	if err != nil {
		return nil, fmt.Errorf("unable to create new request because err: %v", err)
	}
	req.Header.Add("Content-Type", contentType)
	for k, v := range fc.defaultHeaders {
		req.Header[k] = v
	}
	resp, err := fc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("err while http.Client.Do: %v", err)
	}
	// Read the resp body
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return resp, fmt.Errorf("error reading the body of the get request due to err: %v", err)
	}
	for _, respBody := range respBodies {
		err = json.Unmarshal(b, respBody)
		if err != nil {
			return resp, fmt.Errorf("error unmarshalling the body to type %t due to err: %v", respBody, err)
		}
		if respBody.Valid() {
			return resp, nil
		}
	}
	return resp, nil
}
