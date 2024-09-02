package bank

import (
	"context"

	"github.com/fajaramaulana/go-grpc-micro-bank-client/internal/port"
	"github.com/fajaramaulana/go-grpc-micro-bank-proto/protogen/go/bank"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type BankAdapter struct {
	bankClient port.BankClientPort
}

func NewBankAdapter(conn *grpc.ClientConn) (*BankAdapter, error) {
	client := bank.NewBankServiceClient(conn)
	return &BankAdapter{
		bankClient: client,
	}, nil
}

func (a *BankAdapter) GetBalanceByAccountNumber(ctx context.Context, accNum string) (*bank.CurrentBalanceResponse, error) {
	req := &bank.CurrentBalanceRequest{
		AccountNumber: accNum,
	}

	res, err := a.bankClient.GetCurrentBalance(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		log.Fatal().Msgf("Error on GetCurrentBalance :: %s", st.Message())
	}

	return res, nil
}
