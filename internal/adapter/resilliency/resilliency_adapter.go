package resilliency

import (
	"context"
	"io"

	"github.com/fajaramaulana/go-grpc-micro-bank-client/internal/port"
	"github.com/fajaramaulana/go-grpc-micro-bank-proto/protogen/go/resilliency"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type ResilliencyAdapter struct {
	resilliencyClient port.ResilliencyPort
}

func NewResilliencyAdapter(conn *grpc.ClientConn) (*ResilliencyAdapter, error) {
	client := resilliency.NewResilliencyServiceClient(conn)
	return &ResilliencyAdapter{
		resilliencyClient: client,
	}, nil
}

func (a *ResilliencyAdapter) UnaryResilliency(ctx context.Context, minDelaySec int32, maxDelaySec int32, statusCodes []uint32) (*resilliency.ResilliencyResponse, error) {
	req := &resilliency.ResilliencyRequest{
		MinDelaySecond: minDelaySec,
		MaxDelaySecond: maxDelaySec,
		StatusCodes:    statusCodes,
	}

	res, err := a.resilliencyClient.Unaryresilliency(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("Error on Unaryresilliency")
		return nil, err
	}

	return res, nil
}

func (a *ResilliencyAdapter) ServerStreamResilliency(ctx context.Context, minDelaySec int32, maxDelaySec int32, statusCodes []uint32) {
	req := &resilliency.ResilliencyRequest{
		MinDelaySecond: minDelaySec,
		MaxDelaySecond: maxDelaySec,
		StatusCodes:    statusCodes,
	}

	stream, err := a.resilliencyClient.ServerStreamingResilliency(ctx, req)

	if err != nil {
		log.Error().Err(err).Msg("Error on ServerStreamingResilliency")
		return
	}

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			log.Info().Msg("ServerStreamingResilliency stream closed")
			break
		}

		if err != nil {
			log.Fatal().Err(err).Msg("Error on ServerStreamingResilliency: " + err.Error())
		}

		log.Info().Msgf("ServerStreamingResilliency response: %v", res)
	}
}

func (a *ResilliencyAdapter) ClientStreamResilliency(ctx context.Context, minDelaySec int32, maxDelaySec int32, statusCodes []uint32, count int) {
	stream, err := a.resilliencyClient.ClientStreamingResilliency(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error on ClientStreamingResilliency")
		return
	}

	for i := 0; i < count; i++ {
		req := &resilliency.ResilliencyRequest{
			MinDelaySecond: minDelaySec,
			MaxDelaySecond: maxDelaySec,
			StatusCodes:    statusCodes,
		}

		stream.Send(req)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Error().Err(err).Msg("Error on ClientStreamingResilliency")
		return
	}

	log.Info().Msgf("ClientStreamingResilliency response: %v", res)
}

func (a *ResilliencyAdapter) BidirectionalStreamingResilliency(ctx context.Context, minDelaySec int32, maxDelaySec int32, statusCodes []uint32, count int) {
	stream, err := a.resilliencyClient.BidirectionalStreamingResilliency(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error on BidirectionalStreamingResilliency")
		return
	}

	resilliencyChan := make(chan *resilliency.ResilliencyResponse)

	go func() {
		for i := 0; i < count; i++ {
			req := &resilliency.ResilliencyRequest{
				MinDelaySecond: minDelaySec,
				MaxDelaySecond: maxDelaySec,
				StatusCodes:    statusCodes,
			}
			stream.Send(req)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(resilliencyChan)
				break
			}

			if err != nil {
				log.Error().Err(err).Msg("Error on BidirectionalStreamingResilliency")
				return
			}
			log.Info().Msgf("BidirectionalStreamingResilliency response: %v", res)
			resilliencyChan <- res
		}
		close(resilliencyChan)
	}()

	<-resilliencyChan
}
