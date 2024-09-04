package main

// Package main is the entry point of the application.
// It initializes the necessary components and starts the gRPC server.
import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	cfg "github.com/fajaramaulana/go-grpc-micro-bank-client/config"
	"github.com/fajaramaulana/go-grpc-micro-bank-client/internal/adapter/bank"
	domain "github.com/fajaramaulana/go-grpc-micro-bank-client/internal/application/domain/bank"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	configuration := cfg.New("../.env")
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	hostPort := fmt.Sprintf("%s:%s", "localhost", configuration.Get("PORT"))
	conn, err := grpc.NewClient(hostPort, opts...)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create connection")
	}
	defer conn.Close()

	// helloAdapter, err := hello.NewHelloAdapter(conn)

	bankAdapter, err := bank.NewBankAdapter(conn)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to create bank adapter")
	}

	// context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	// runGetCurrentBalance(ctx, bankAdapter, "7835697001")
	// resGetExchangeRatesStream(ctx, bankAdapter, "IDR", "USD")
	resGetSummarizeTransactions(ctx, bankAdapter, "7835697001", 10)
}

// func runGetCurrentBalance(ctx context.Context, adapter *bank.BankAdapter, acct string) {
// 	res, err := adapter.GetBalanceByAccountNumber(ctx, acct)

// 	if err != nil {
// 		log.Fatal().Msg("Failed to call GetCurrentBalance: " + err.Error())
// 	}

// 	log.Info().Msgf("Current balance for account %s is %f", acct, res.Amount)
// }

// func resGetExchangeRatesStream(ctx context.Context, adapter *bank.BankAdapter, fromCurrency string, toCurrency string) {
// 	adapter.GetExchangeRatesStream(ctx, fromCurrency, toCurrency)
// }

func resGetSummarizeTransactions(ctx context.Context, adapter *bank.BankAdapter, account string, numDummyTransactions int) {
	var trx []domain.Transaction

	for i := 1; i <= numDummyTransactions; i++ {
		trxType := domain.TransactionTypeIn

		if i%3 == 0 {
			trxType = domain.TransactionTypeOut
		}

		t := domain.Transaction{
			Amount:          float64(rand.Intn(500) + 10),
			TransactionType: trxType,
			Notes:           "Transaction " + strconv.Itoa(i),
		}
		fmt.Printf("%# v\n", t)

		trx = append(trx, t)
	}

	adapter.GetSummarizeTransactions(ctx, account, trx)
}
