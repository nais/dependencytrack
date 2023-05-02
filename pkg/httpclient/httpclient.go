package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type RequestError struct {
	StatusCode int
	Err        error
}

type HttpClient struct {
	*http.Client
	log              *log.Entry
	responseCallback func(res *http.Response, err error)
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

func (r *RequestError) Error() string {
	return fmt.Sprintf("status %d: err %v", r.StatusCode, r.Err)
}

func (c *HttpClient) SendRequest(ctx context.Context, httpMethod string, url string, headers map[string][]string, body []byte) ([]byte, error) {
	c.log.Debugf("sending request to %s", url)
	req, err := http.NewRequestWithContext(ctx, httpMethod, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header = headers

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	c.responseCallback(resp, err)

	if resp.StatusCode > 299 {
		c.log.Debugf("response status: %d", resp.StatusCode)
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("reading response body: %w", err)
		}
		return nil, fail(resp.StatusCode, fmt.Errorf("%s", string(b)))
	}
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	return resBody, err
}

func fail(status int, err error) *RequestError {
	return &RequestError{
		StatusCode: status,
		Err:        err,
	}
}
