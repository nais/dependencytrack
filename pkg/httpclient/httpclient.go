package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"github.com/avast/retry-go/v4"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Response struct {
	Status  int
	Body    []byte
	Headers http.Header
}

type RequestError struct {
	StatusCode int
	Err        error
}

type HttpClient struct {
	*http.Client
	log              *log.Entry
	responseCallback func(res *http.Response, err error)
	retryOps         []retry.Option
}

type Option func(*HttpClient)

func New(opts ...Option) *HttpClient {
	c := &HttpClient{}
	for _, opt := range opts {
		opt(c)
	}

	if c.Client == nil {
		c.Client = http.DefaultClient
	}
	if c.log == nil {
		c.log = log.NewEntry(log.New())
	}
	if c.retryOps == nil {
		c.retryOps = []retry.Option{
			retry.Attempts(3),
			// delay 500ms
			retry.Delay(500),
			retry.LastErrorOnly(true),
		}
	}
	if c.responseCallback == nil {
		c.responseCallback = func(res *http.Response, err error) {}
	}
	return c
}

func WithResponseCallback(responseCallback func(res *http.Response, err error)) Option {
	return func(c *HttpClient) {
		c.responseCallback = responseCallback
	}
}

func WithLogger(log *log.Entry) Option {
	return func(c *HttpClient) {
		c.log = log
	}
}

func WithClient(client *http.Client) Option {
	return func(c *HttpClient) {
		c.Client = client
	}
}

func WithRetryOptions(opts ...retry.Option) Option {
	return func(c *HttpClient) {
		c.retryOps = opts
	}
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("status %d: err %v", r.StatusCode, r.Err)
}

func (c *HttpClient) SendRequest(ctx context.Context, httpMethod string, url string, headers map[string][]string, body []byte) (*Response, error) {
	c.log.Debugf("sending request to %s", url)
	req, err := http.NewRequestWithContext(ctx, httpMethod, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header = headers

	var returnResponse *Response
	var resp *http.Response
	err = retry.Do(func() error {
		resp, err = c.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		c.responseCallback(resp, err)
		if resp.StatusCode > 299 {
			c.log.Debugf("response status: %d", resp.StatusCode)
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				c.log.Debugf("error reading response body: %v", err)
				returnResponse = &Response{
					Status:  resp.StatusCode,
					Body:    []byte{},
					Headers: resp.Header,
				}
				return nil
			}
			return fail(resp.StatusCode, fmt.Errorf("%s", string(b)))
		}
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("reading response body: %w", err)
		}
		returnResponse = &Response{
			Status:  resp.StatusCode,
			Body:    respBody,
			Headers: resp.Header,
		}
		return nil
	}, c.retryOps...)
	if err != nil {
		return nil, err
	}
	return returnResponse, nil
}

func fail(status int, err error) *RequestError {
	return &RequestError{
		StatusCode: status,
		Err:        err,
	}
}
