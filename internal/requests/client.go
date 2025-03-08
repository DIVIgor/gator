package requests

import (
	"context"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient http.Client
}

func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

func newRequest(ctx context.Context, method, url string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {return}

	req.Header.Set("User-Agent", "gator")

	return req, err
}

func (c *Client) MakeRequest(ctx context.Context, method, url string, body io.Reader) (resp *http.Response, err error) {
	req, err := newRequest(ctx, method, url, body)
	if err != nil {return}

	return c.httpClient.Do(req)
}