package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type UnaryTransport interface {
	Header() http.Header
	Send(ctx context.Context, endpoint, contentType string, body io.Reader) (http.Header, io.ReadCloser, error)
	Close() error
}

type httpTransport struct {
	host   string
	client *http.Client
	opts   *ConnectOptions

	header http.Header

	sent bool
}

func (t *httpTransport) Header() http.Header {
	return t.header
}

func (t *httpTransport) Send(ctx context.Context, endpoint, contentType string, body io.Reader) (http.Header, io.ReadCloser, error) {
	if t.sent {
		return nil, nil, errors.New("Send must be called only one time per one Request")
	}
	defer func() {
		t.sent = true
	}()

	// TODO: HTTPS support.
	scheme := "http"
	u := url.URL{Scheme: scheme, Host: t.host, Path: endpoint}
	url := u.String()
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build the API request: %w", err)
	}

	req.Header = t.Header()
	req.Header.Add("content-type", contentType)
	req.Header.Add("x-grpc-web", "1")

	res, err := t.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send the API: %w", err)
	}

	return res.Header, res.Body, nil
}

func (t *httpTransport) Close() error {
	t.client.CloseIdleConnections()
	return nil
}

var NewUnary = func(host string, opts *ConnectOptions) UnaryTransport {
	return &httpTransport{
		host:   host,
		client: http.DefaultClient,
		opts:   opts,
		header: make(http.Header),
	}
}

type ClientStreamTransport interface {
	Header() (http.Header, error)
	Trailer() http.Header

	// SetRequestHeader sets headers to send gRPC-Web server.
	// It should be called before calling Send.
	SetRequestHeader(h http.Header)
	Send(ctx context.Context, body io.Reader) error
	Receive(ctx context.Context) (io.ReadCloser, error)

	// CloseSend sends a close signal to the server.
	CloseSend() error

	// Close closes the connection.
	Close() error
}
