package curl

//	type HttpCaller interface {
//		GET() (*gentleman.Response, error)
//		POST() (*gentleman.Response, error)
//		PUT() (*gentleman.Response, error)
//		PATCH() (*gentleman.Response, error)
//		DELETE() (*gentleman.Response, error)
//		OPTIONS() (*gentleman.Response, error)
//		SetTimeOut(time.Duration) HttpCaller
//		SetMultipartData(fields multipart.FormData) HttpCaller
//		BasicAuth(username, password string) HttpCaller
//		Bearer(token string) HttpCaller
//		SetQuery(q map[string]string) HttpCaller
//	}
//
//	type Request struct {
//		URL    string
//		Method string
//		Body   interface{}
//		Header map[string]string
//		Client *gentleman.Client
//	}
//
// // NewRest New create the new http caller with parameter url, request body, and http header
//
//	func NewRest(url string, body interface{}, header map[string]string) HttpCaller {
//		cli := gentleman.New()
//		cli.URL(url)
//		//cli.Use(mock.Plugin)
//
//		//cli.Context.Request.Close = true
//		// Define the max timeout for the whole HTTP request
//		cli.Use(timeout.Request(30 * time.Second))
//
//		return &Request{
//			URL:    url,
//			Body:   body,
//			Header: header,
//			Client: cli,
//		}
//	}
//
// // SetClient set new http client
//
//	func (c *Request) SetClient(client *gentleman.Client) {
//		c.Client = client
//	}
//
// // GET func call server with method GET
//
//	func (c *Request) GET() (*gentleman.Response, error) {
//		c.Method = http.MethodGet
//		return invoke(*c)
//	}
//
// // POST func call server with method POST
//
//	func (c *Request) POST() (*gentleman.Response, error) {
//		c.Method = http.MethodPost
//		return invoke(*c)
//	}
//
// // PUT func call server with method PUT
//
//	func (c *Request) PUT() (*gentleman.Response, error) {
//		c.Method = http.MethodPut
//		return invoke(*c)
//	}
//
// // PATCH func call server with method PATCH
//
//	func (c *Request) PATCH() (*gentleman.Response, error) {
//		c.Method = http.MethodPatch
//		return invoke(*c)
//	}
//
// // DELETE func call server with method DELETE
//
//	func (c *Request) DELETE() (*gentleman.Response, error) {
//		c.Method = http.MethodDelete
//		return invoke(*c)
//	}
//
// // OPTIONS func call server with method DELETE
//
//	func (c *Request) OPTIONS() (*gentleman.Response, error) {
//		c.Method = http.MethodOptions
//		return invoke(*c)
//	}
//
// // SetTimeOut Set timeout
//
//	func (c *Request) SetTimeOut(duration time.Duration) HttpCaller {
//		c.Client.Use(timeout.Request(duration))
//		return &Request{
//			URL:    c.URL,
//			Body:   c.Body,
//			Header: c.Header,
//			Client: c.Client,
//		}
//	}
//
// // SetMultipartData set files into multipart/form-data
//
//	func (c *Request) SetMultipartData(fields multipart.FormData) HttpCaller {
//		c.Client.Use(multipart.Data(fields))
//		return &Request{
//			URL:    c.URL,
//			Body:   c.Body,
//			Header: c.Header,
//			Client: c.Client,
//		}
//	}
//
// // SetQuery set quert
//
//	func (c *Request) SetQuery(q map[string]string) HttpCaller {
//		for k, v := range q {
//			c.Client.Use(query.Set(k, v))
//		}
//
//		return &Request{
//			URL:    c.URL,
//			Body:   c.Body,
//			Header: c.Header,
//			Client: c.Client,
//		}
//	}
//
// // BasicAuth basic authen
//
//	func (c *Request) BasicAuth(username, password string) HttpCaller {
//		c.Client.Use(auth.Basic(username, password))
//		return &Request{
//			URL:    c.URL,
//			Body:   c.Body,
//			Header: c.Header,
//			Client: c.Client,
//		}
//	}
//
// // Bearer authen
//
//	func (c *Request) Bearer(token string) HttpCaller {
//		c.Client.Use(auth.Bearer(token))
//		return &Request{
//			URL:    c.URL,
//			Body:   c.Body,
//			Header: c.Header,
//			Client: c.Client,
//		}
//	}
//
//	func invoke(v Request) (*gentleman.Response, error) {
//		contentType, _ := v.Header["Content-Type"]
//		normalizeType := strings.ToLower(contentType)
//
//		if strings.Contains(normalizeType, "json") {
//			req := v.Client.Request()
//			req.Method(v.Method)
//			req.SetHeaders(v.Header)
//			req.Use(body.JSON(v.Body))
//
//			return req.Send()
//		} else if strings.Contains(normalizeType, "xml") {
//			req := v.Client.Request()
//			req.Method(v.Method)
//			req.SetHeaders(v.Header)
//			req.Use(body.XML(v.Body))
//
//			return req.Send()
//		} else if strings.Contains(normalizeType, "x-www-form-urlencode") || strings.Contains(normalizeType, "text/plain") {
//			if isNil(v.Body) {
//				v.Body = ""
//			}
//
//			req := v.Client.Request()
//			req.Method(v.Method)
//			req.SetHeaders(v.Header)
//
//			req.Use(body.String(v.Body.(string)))
//
//			return req.Send()
//		} else {
//			req := v.Client.Request()
//			req.Method(v.Method)
//			req.SetHeaders(v.Header)
//
//			return req.Send()
//		}
//	}
//func isNil(v interface{}) bool {
//	return v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil())
//}
