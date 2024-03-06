package curl

import (
	"encoding/base64"
	"encoding/xml"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/valyala/fasthttp"
)

type HttpCaller interface {
	GET() ([]byte, int, http.Header, error)
	POST() ([]byte, int, http.Header, error)
	PUT() ([]byte, int, http.Header, error)
	PATCH() ([]byte, int, http.Header, error)
	DELETE() ([]byte, int, http.Header, error)
	OPTIONS() ([]byte, int, http.Header, error)
	SetTimeOut(time.Duration) HttpCaller
	Bearer(token string) HttpCaller
	BasicAuth(username, password string) HttpCaller
	// SetQuery(q map[string]string) HttpCaller
}

var Sec3 = 3 * time.Second
var Sec5 = 5 * time.Second
var Sec15 = 15 * time.Second

var Cli = &fasthttp.Client{
	MaxConnsPerHost: 10000,
}

type Request struct {
	Timeout time.Duration
	URL     string
	Method  string
	Body    interface{}
	Header  map[string]string
	Client  *fasthttp.Client
	Writer  *multipart.Writer
}

// New create the new http caller with parameter url, request body, and http header
func NewRest(url string, body interface{}, header map[string]string) HttpCaller {
	return &Request{
		Timeout: Sec5,
		URL:     url,
		Body:    body,
		Header:  header,
		Client:  Cli,
	}
}

func NewMultipartData(url string, body []byte, header map[string]string, writer *multipart.Writer) HttpCaller {
	if len(header) == 0 {
		header = make(map[string]string)
	}
	return &Request{
		Timeout: Sec15,
		URL:     url,
		Body:    body,
		Header:  header,
		Client:  Cli,
		Writer:  writer,
	}
}

// New create the new http caller with parameter url, request body, and http header
func NewWithHeaderWithTimeout(url string, body interface{}, header map[string]string, timeoutDuration time.Duration) HttpCaller {
	return &Request{
		Timeout: timeoutDuration,
		URL:     url,
		Body:    body,
		Header:  header,
		Client:  Cli,
	}
}

// GET func call server with method GET
func (c *Request) GET() ([]byte, int, http.Header, error) {
	c.Method = http.MethodGet
	header := http.Header{}
	resp, err := invoke(*c)
	if err != nil {
		return nil, 0, nil, err
	}

	resp.Header.VisitAll(func(k, v []byte) {
		header[string(k)] = append(header[string(k)], string(v))
	})

	return resp.Body(), resp.Header.StatusCode(), header, err
}

// POST func call server with method POST
func (c *Request) POST() ([]byte, int, http.Header, error) {
	c.Method = http.MethodPost
	header := http.Header{}
	resp, err := invoke(*c)
	if err != nil {
		return nil, 0, nil, err
	}
	resp.Header.VisitAll(func(k, v []byte) {
		header[string(k)] = append(header[string(k)], string(v))
	})

	return resp.Body(), resp.Header.StatusCode(), header, err
}

// PUT func call server with method PUT
func (c *Request) PUT() ([]byte, int, http.Header, error) {
	c.Method = http.MethodPut
	header := http.Header{}
	resp, err := invoke(*c)
	if err != nil {
		return nil, 0, nil, err
	}
	resp.Header.VisitAll(func(k, v []byte) {
		header[string(k)] = append(header[string(k)], string(v))
	})

	return resp.Body(), resp.Header.StatusCode(), header, err
}

// PATCH func call server with method PATCH
func (c *Request) PATCH() ([]byte, int, http.Header, error) {
	c.Method = http.MethodPatch
	header := http.Header{}
	resp, err := invoke(*c)
	if err != nil {
		return nil, 0, nil, err
	}
	resp.Header.VisitAll(func(k, v []byte) {
		header[string(k)] = append(header[string(k)], string(v))
	})

	return resp.Body(), resp.Header.StatusCode(), header, err
}

// DELETE func call server with method DELETE
func (c *Request) DELETE() ([]byte, int, http.Header, error) {
	c.Method = http.MethodDelete
	header := http.Header{}
	resp, err := invoke(*c)
	if err != nil {
		return nil, 0, nil, err
	}
	resp.Header.VisitAll(func(k, v []byte) {
		header[string(k)] = append(header[string(k)], string(v))
	})

	return resp.Body(), resp.Header.StatusCode(), header, err
}

// OPTIONS func call server with method DELETE
func (c *Request) OPTIONS() ([]byte, int, http.Header, error) {
	c.Method = http.MethodOptions
	header := http.Header{}
	resp, err := invoke(*c)
	if err != nil {
		return nil, 0, nil, err
	}
	resp.Header.VisitAll(func(k, v []byte) {
		header[string(k)] = append(header[string(k)], string(v))
	})

	return resp.Body(), resp.Header.StatusCode(), header, err
}

// Bearer authen
func (c *Request) Bearer(token string) HttpCaller {
	c.Header["Authorization"] = "Bearer " + token
	return &Request{
		Timeout: c.Timeout,
		URL:     c.URL,
		Method:  c.Method,
		Body:    c.Body,
		Header:  c.Header,
		Client:  c.Client,
		Writer:  c.Writer,
	}
}

// SetTimeOut Set timeout
func (c *Request) SetTimeOut(duration time.Duration) HttpCaller {
	return &Request{
		Timeout: duration,
		URL:     c.URL,
		Body:    c.Body,
		Header:  c.Header,
		Client:  c.Client,
	}
}

// SetMultipartData set files into multipart/form-data
//func (c *Request) SetMultipartData(body []byte, writer *multipart.Writer) HttpCaller {
//	return &Request{
//		URL:    c.URL,
//		Body:   body,
//		Header: c.Header,
//		Client: c.Client,
//	}
//}

// SetQueryParam set files into multipart/form-data
// func (c *Request) SetQuery(q map[string]string) HttpCaller {
// 	for k, v := range q {
// 		c.Client.Use(query.Set(k, v))
// 	}

// 	return &Request{
// 		URL:    c.URL,
// 		Body:   c.Body,
// 		Header: c.Header,
// 		Client: c.Client,
// 	}
// }

// BasicAuth basic authen
func (c *Request) BasicAuth(username, password string) HttpCaller {

	c.Header["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))

	return &Request{
		URL:    c.URL,
		Body:   c.Body,
		Header: c.Header,
		Client: c.Client,
	}
}

func invoke(v Request) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()

	for k, v := range v.Header {
		req.Header.Set(k, v)
	}

	req.SetRequestURI(v.URL)
	req.Header.SetMethod(v.Method)
	if v.Writer != nil {
		req.Header.SetContentType(v.Writer.FormDataContentType())
	}
	contentType := v.Header["Content-Type"]
	//if !ok {
	//	return &gentleman.Response{}, errors.New("please provide Content-Type")
	//}
	normalizeType := strings.ToLower(contentType)
	resp := fasthttp.AcquireResponse()
	if strings.Contains(normalizeType, "json") {
		body, err := sonic.Marshal(v.Body)
		if err != nil {
			return resp, err
		}

		req.SetBody(body)

		err = v.Client.DoTimeout(req, resp, v.Timeout)
		if err != nil {

			return resp, err
		}

		return resp, nil
	} else if strings.Contains(normalizeType, "xml") {
		body, err := xml.Marshal(v.Body)
		if err != nil {
			return resp, err
		}

		req.SetBody(body)

		err = v.Client.DoTimeout(req, resp, v.Timeout)
		if err != nil {
			return resp, err
		}

		return resp, nil
	} else if strings.Contains(normalizeType, "x-www-form-urlencode") || strings.Contains(normalizeType, "text/plain") {
		if isNil(v.Body) {
			v.Body = ""
		}

		req.SetBodyRaw([]byte(v.Body.(string)))

		err := v.Client.DoTimeout(req, resp, v.Timeout)
		if err != nil {
			return resp, err
		}

		return resp, nil
	} else {
		switch body := v.Body.(type) {
		case []byte:
			req.SetBody(body)
		case string:
			req.SetBody([]byte(body))
		default:
			bData, _ := sonic.Marshal(body)
			req.SetBody(bData)
		}
		err := v.Client.DoTimeout(req, resp, v.Timeout)
		if err != nil {
			return resp, err
		}

		return resp, nil
	}
}

func isNil(v interface{}) bool {
	return v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil())
}
