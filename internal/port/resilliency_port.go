package port

import (
	"context"

	"github.com/fajaramaulana/go-grpc-micro-bank-proto/protogen/go/resilliency"
	"google.golang.org/grpc"
)

type ResilliencyPort interface {
	Unaryresilliency(ctx context.Context, in *resilliency.ResilliencyRequest, opts ...grpc.CallOption) (*resilliency.ResilliencyResponse, error)
	ServerStreamingResilliency(ctx context.Context, in *resilliency.ResilliencyRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[resilliency.ResilliencyResponse], error)
	ClientStreamingResilliency(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[resilliency.ResilliencyRequest, resilliency.ResilliencyResponse], error)
	BidirectionalStreamingResilliency(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[resilliency.ResilliencyRequest, resilliency.ResilliencyResponse], error)
}
