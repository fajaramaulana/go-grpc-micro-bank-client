package bank

import (
	"context"
	"io"

	domain "github.com/fajaramaulana/go-grpc-micro-bank-client/internal/application/domain/bank"
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

func (a *BankAdapter) GetExchangeRatesStream(ctx context.Context, fromCurrency string, toCurrency string) {
	req := &bank.ExchangeRateRequest{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
	}

	stream, err := a.bankClient.FetchExchangeRates(ctx, req)
	if err != nil {
		log.Fatal().Msg("Error on FetchExchangeRates :: " + err.Error())
	}

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal().Msg("Error on FetchExchangeRates :: " + err.Error())
		}

		log.Info().Msgf("Exchange rate from %s to %s is %f", res.FromCurrency, res.ToCurrency, res.Rate)
	}
}

func (a *BankAdapter) GetSummarizeTransactions(ctx context.Context, account string, trx []domain.Transaction) {
	trxStream, err := a.bankClient.SummarizeTransactions(ctx)
	if err != nil {
		log.Fatal().Msg("Error on a.bankClient.SummarizeTransactions :: " + err.Error())
	}

	for _, v := range trx {
		trxType := bank.TransactionType_TRANSACTION_TYPE_UNSPECIFIED

		if v.TransactionType == domain.TransactionTypeIn {
			trxType = bank.TransactionType_TRANSACTION_TYPE_IN
		} else {
			trxType = bank.TransactionType_TRANSACTION_TYPE_OUT
		}

		req := &bank.Transaction{
			AccountNumber: account,
			Type:          trxType,
			Amount:        v.Amount,
			Notes:         v.Notes,
		}

		trxStream.Send(req)
	}

	summary, err := trxStream.CloseAndRecv()
	if err != nil {
		stat, _ := status.FromError(err)

		log.Fatal().Msgf("Error on trxStream.CloseAndRecv GetSummarizeTransactions :: %s", stat.Message())
	}

	log.Info().Msgf("Transaction summary :: %v", summary)
}
