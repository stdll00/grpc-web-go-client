package grpcweb

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/stdll00/grpc-web-go-client/grpcweb/parser"
	"github.com/stdll00/grpc-web-go-client/grpcweb/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
)

type ClientConn struct {
	host        string
	dialOptions *dialOptions
}

func (c *ClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented client side")
}

func Dial(host string, opts ...DialOption) (*ClientConn, error) {
	opt := defaultDialOptions
	// Validate host!
	for _, o := range opts {
		o(&opt)
	}
	return &ClientConn{
		host:        host,
		dialOptions: &opt,
	}, nil
}

func (c *ClientConn) Invoke(ctx context.Context, method string, args, reply any, opts ...grpc.CallOption) error {
	callOptions, err := c.applyCallOptions(opts)
	if err != nil {
		return err
	}
	codec := callOptions.codec

	tr := transport.NewUnary(c.host, nil)
	defer tr.Close()

	r, err := encodeRequestBody(codec, args)
	if err != nil {
		return fmt.Errorf("failed to build the request body: %w", err)
	}

	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		for k, v := range md {
			for _, vv := range v {
				tr.Header().Add(k, vv)
			}
		}
	}
	contentType := "application/grpc-web+" + codec.Name()
	header, rawBody, err := tr.Send(ctx, method, contentType, r)
	if err != nil {
		return fmt.Errorf("failed to send the request: %w", err)
	}
	defer rawBody.Close()

	if callOptions.header != nil {
		*callOptions.header = toMetadata(header)
	}

	resHeader, err := parser.ParseResponseHeader(rawBody)
	if err != nil {
		return status.Error(codes.Unavailable, fmt.Sprintf("failed to parse response header: %s", err.Error()))
	}

	if resHeader.IsMessageHeader() {
		resBody, err := parser.ParseLengthPrefixedMessage(rawBody, resHeader.ContentLength)
		if err != nil {
			return fmt.Errorf("failed to parse the response body: %w", err)
		}
		if err := codec.Unmarshal(resBody, reply); err != nil {
			return fmt.Errorf("failed to unmarshal response body by codec %s: %w", codec.Name(), err)
		}

		resHeader, err = parser.ParseResponseHeader(rawBody)
		if err != nil {
			return fmt.Errorf("failed to parse response header: %w", err)
		}
	}
	if !resHeader.IsTrailerHeader() {
		return errors.New("unexpected header")
	}

	st, trailer, err := parser.ParseStatusAndTrailer(rawBody, resHeader.ContentLength)
	if err != nil {
		return fmt.Errorf("failed to parse status and trailer: %w", err)
	}
	if callOptions.trailer != nil {
		*callOptions.trailer = trailer
	}

	return st.Err()
}

func (c *ClientConn) applyCallOptions(opts []grpc.CallOption) (*callOptions, error) {
	callOptions := defaultCallOptions
	for _, o := range opts {
		switch v := o.(type) {
		case grpc.CustomCodecCallOption:
			callOptions.codec = convertToEncodingCodec(v.Codec)
		case grpc.HeaderCallOption:
			callOptions.header = v.HeaderAddr
		case grpc.TrailerCallOption:
			callOptions.trailer = v.TrailerAddr
		case grpc.EmptyCallOption:
		default:
			return nil, fmt.Errorf("%v is not supported type", o)
		}
	}
	return &callOptions, nil
}

// copied from rpc_util.go#msgHeader
const headerLen = 5

func header(body []byte) []byte {
	h := make([]byte, 5)
	h[0] = byte(0)
	binary.BigEndian.PutUint32(h[1:], uint32(len(body)))
	return h
}

// header (compressed-flag(1) + message-length(4)) + body
// TODO: compressed message
func encodeRequestBody(codec encoding.Codec, in interface{}) (io.Reader, error) {
	body, err := codec.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal the request body: %w", err)
	}
	buf := bytes.NewBuffer(make([]byte, 0, headerLen+len(body)))
	buf.Write(header(body))
	buf.Write(body)
	return buf, nil
}

func toMetadata(h http.Header) metadata.MD {
	if len(h) == 0 {
		return nil
	}
	md := metadata.New(nil)
	for k, v := range h {
		md.Append(k, v...)
	}
	return md
}
