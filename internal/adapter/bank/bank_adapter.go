package bank

import (
	"context"
	"io"

	domain "github.com/fajaramaulana/go-grpc-micro-bank-client/internal/application/domain/bank"
	"github.com/fajaramaulana/go-grpc-micro-bank-client/internal/port"
	"github.com/fajaramaulana/go-grpc-micro-bank-proto/protogen/go/bank"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
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

func (a *BankAdapter) SendTransferMultiple(ctx context.Context, trx []domain.TransferTransaction) {
	trxStream, err := a.bankClient.TransferMultiple(ctx)

	if err != nil {
		log.Fatal().Msgf("Error on SendTransferMultiple - a.bankClient.TransferMultiple :: %s", err.Error())
	}

	trxChannel := make(chan struct{})

	// Goroutine to send requests
	go func() {
		for _, v := range trx {
			req := &bank.TransferRequest{
				AccountNumberSender:   v.FromAccountNumber,
				AccountNumberReciever: v.ToAccountNumber,
				Currency:              v.Currency,
				Amount:                v.Amount,
				Notes:                 "Transfer from ",
			}

			// Send each transfer request
			if err := trxStream.Send(req); err != nil {
				log.Error().Msgf("Error sending transfer request: %s", err.Error())
				return // Exit goroutine on send error
			}
		}

		// Close the send stream after all requests have been sent
		trxStream.CloseSend()
	}()

	// Goroutine to receive responses
	go func() {
		for {
			res, err := trxStream.Recv()
			if err == io.EOF {
				log.Info().Msg("EOF received, closing receive goroutine.")
				break // Normal end of the stream
			}

			if err != nil {
				handleTransferErrorGrpc(err)
				break
			} else {
				// log.Printf("Transfer status %v on %v\n", res.Status, res.Timestamp)
				log.Info().Msgf("Transfer status %v on %v\n", res.Status, res.Timestamp)
			}
		}

		close(trxChannel)
	}()

	<-trxChannel // Wait until the receive goroutine signals completion
}

func handleTransferErrorGrpc(err error) {
	st := status.Convert(err)
	log.Debug().Msgf("Error %v on TransferMultiple : %v", st.Code(), st.Message())

	for _, detail := range st.Details() {
		switch t := detail.(type) {
		case *errdetails.PreconditionFailure:
			for _, violation := range t.GetViolations() {
				log.Info().Msgf("[VIOLATION] %v", violation)
			}
		case *errdetails.ErrorInfo:
			log.Info().Msgf("Error on : %v, with reason %v\n", t.Domain, t.Reason)
			for k, v := range t.GetMetadata() {
				log.Info().Msgf("  %v : %v\n", k, v)
			}
		}
	}
}
