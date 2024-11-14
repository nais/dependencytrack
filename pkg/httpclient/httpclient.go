package httpclient

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

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
	maxRetries       int
	retryDelay       time.Duration
	defaultTimeout   time.Duration
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
	if c.maxRetries == 0 {
		c.maxRetries = 1
		c.retryDelay = 0 * time.Second
	}
	if c.defaultTimeout == 0 {
		c.defaultTimeout = 5 * time.Second
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

func WithRetry(maxRetries int, retryDelay time.Duration) Option {
	return func(c *HttpClient) {
		c.maxRetries = maxRetries
		c.retryDelay = retryDelay
	}
}

func WithDefaultTimeout(timeout time.Duration) Option {
	return func(c *HttpClient) {
		c.defaultTimeout = timeout
	}
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("status %d: err %v", r.StatusCode, r.Err)
}

func (c *HttpClient) SendRequest(ctx context.Context, httpMethod, requestUrl string, headers map[string][]string, body []byte) (*Response, error) {
	c.log.Debugf("sending request to %s", requestUrl)

	// Set a default timeout if not already set in the context
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.defaultTimeout)
		defer cancel()
	}

	req, err := http.NewRequestWithContext(ctx, httpMethod, requestUrl, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header = headers

	var resp *http.Response
	for i := 0; i <= c.maxRetries; i++ {
		resp, err = c.Do(req)
		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && (netErr.Timeout()) {
				if i < c.maxRetries {
					c.log.Warnf("timeout network error, retrying... (%d/%d)", i+1, c.maxRetries)
					time.Sleep(c.retryDelay)
					continue
				}
				return nil, fmt.Errorf("request failed after retries due to timeout: %w", err)
			}
			return nil, fmt.Errorf("request failed: %w", err)
		}

		// Handle response and exit retry loop if successful
		defer resp.Body.Close()
		c.responseCallback(resp, err)

		// Retry only on server errors
		if resp.StatusCode != http.StatusInternalServerError && resp.StatusCode != http.StatusServiceUnavailable {
			break
		}

		if i < c.maxRetries {
			c.log.Warnf("received %d status, retrying... (%d/%d)", resp.StatusCode, i+1, c.maxRetries)
			time.Sleep(c.retryDelay)
		}
	}

	// Check if response was successful or if an error occurred after retries
	if resp == nil {
		return nil, fmt.Errorf("failed to obtain a response from server after retries")
	}

	if resp.StatusCode > 299 {
		c.log.Debugf("response status: %d", resp.StatusCode)
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			c.log.Debugf("error reading response body: %v", err)
			return &Response{
				Status:  resp.StatusCode,
				Body:    []byte{},
				Headers: resp.Header,
			}, nil
		}
		return nil, fail(resp.StatusCode, fmt.Errorf("%s", string(b)))
	}

	// Read and return the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	return &Response{
		Status:  resp.StatusCode,
		Body:    respBody,
		Headers: resp.Header,
	}, nil
}

func fail(status int, err error) *RequestError {
	return &RequestError{
		StatusCode: status,
		Err:        err,
	}
}
