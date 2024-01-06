package grpcweb

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/proto"
	"google.golang.org/grpc/metadata"
)

var (
	defaultDialOptions = dialOptions{}
	defaultCallOptions = callOptions{
		codec: encoding.GetCodec(proto.Name),
	}
)

type dialOptions struct {
	defaultCallOptions   callOptions
	insecure             bool
	transportCredentials credentials.TransportCredentials
}

type DialOption func(*dialOptions)

func WithInsecure() DialOption {
	return func(opt *dialOptions) {
		opt.insecure = true
	}
}

func WithTransportCredentials(creds credentials.TransportCredentials) DialOption {
	return func(opt *dialOptions) {
		opt.transportCredentials = creds
	}
}

type callOptions struct {
	codec           encoding.Codec
	header, trailer *metadata.MD
}

type grpcCodecWrapper struct {
	grpc.Codec
}

func (g grpcCodecWrapper) Name() string {
	return g.String()
}

func convertToEncodingCodec(codec grpc.Codec) encoding.Codec {
	return &grpcCodecWrapper{Codec: codec}
}

var _ encoding.Codec = &grpcCodecWrapper{nil}
