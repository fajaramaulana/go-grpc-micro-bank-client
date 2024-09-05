package main

// Package main is the entry point of the application.
// It initializes the necessary components and starts the gRPC server.
import (
	"context"
	"fmt"
	"os"
	"time"

	cfg "github.com/fajaramaulana/go-grpc-micro-bank-client/config"
	"github.com/fajaramaulana/go-grpc-micro-bank-client/internal/adapter/resilliency"
	"github.com/fajaramaulana/go-grpc-micro-bank-client/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	domainResil "github.com/fajaramaulana/go-grpc-micro-bank-client/internal/application/domain/resilliency"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	configuration := cfg.New("../.env")
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// Adding a retry interceptor for unary calls
	opts = append(opts, grpc.WithUnaryInterceptor(util.CreateRetryInterceptor()))
	// Adding a retry interceptor for stream calls
	opts = append(opts, grpc.WithStreamInterceptor(util.CreateStreamRetryInterceptor()))

	hostPort := fmt.Sprintf("%s:%s", "localhost", configuration.Get("PORT"))
	conn, err := grpc.NewClient(hostPort, opts...)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create connection")
	}
	defer conn.Close()

	// helloAdapter, err := hello.NewHelloAdapter(conn)

	/* bankAdapter, err := bank.NewBankAdapter(conn)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to create bank adapter")
	} */
	resilliencyAdapter, err := resilliency.NewResilliencyAdapter(conn)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create resilliency adapter")
	}

	// context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// runGetCurrentBalance(bankAdapter, "7835697001xxxxx")
	// runFetchExchangeRates(bankAdapter, "USD", "GBP")
	// runSummarizeTransactions(bankAdapter, "7835697002yyyyy", 10)
	// runTransferMultiple(bankAdapter, "7835697004", "7835697003", 200)

	// runUnaryResiliencyWithTimeout(ctx, resilliencyAdapter, 2, 8, []uint32{domainResil.OK}, 5*time.Second)
	// runServerStreamingResiliencyWithTimeout(ctx, resilliencyAdapter, 0, 3, []uint32{domainResil.OK}, 15*time.Second)
	// runClientStreamingResiliencyWithTimeout(ctx, resilliencyAdapter, 0, 3, []uint32{domainResil.OK}, 10, 10*time.Second)
	// runBiDirectionalResiliencyWithTimeout(ctx, resilliencyAdapter, 0, 3, []uint32{domainResil.OK}, 10, 10*time.Second)
	// runUnaryResiliency(ctx, resilliencyAdapter, 0, 3, []uint32{domainResil.UNKNOWN, domainResil.OK})
	// runServerStreamingResilliency(ctx, resilliencyAdapter, 0, 3, []uint32{domainResil.UNKNOWN})
	// runClientStreamResiliency(ctx, resilliencyAdapter, 0, 3, []uint32{domainResil.UNKNOWN}, 10)
	runUnaryResiliency(ctx, resilliencyAdapter, 0, 3, []uint32{domainResil.UNKNOWN, domainResil.OK})
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

// func resGetSummarizeTransactions(ctx context.Context, adapter *bank.BankAdapter, account string, numDummyTransactions int) {
// 	var trx []domain.Transaction

// 	for i := 1; i <= numDummyTransactions; i++ {
// 		trxType := domain.TransactionTypeIn

// 		if i%3 == 0 {
// 			trxType = domain.TransactionTypeOut
// 		}

// 		t := domain.Transaction{
// 			Amount:          float64(rand.Intn(500) + 10),
// 			TransactionType: trxType,
// 			Notes:           "Transaction " + strconv.Itoa(i),
// 		}
// 		fmt.Printf("%# v\n", t)

// 		trx = append(trx, t)
// 	}

// 	adapter.GetSummarizeTransactions(ctx, account, trx)
// }

// func resSendTransferMultiple(ctx context.Context, adapter *bank.BankAdapter, accountSender string, accountReciever string, numTransactions int) {
// 	var reqs []domain.TransferTransaction

// 	for i := 1; i <= numTransactions; i++ {
// 		t := domain.TransferTransaction{
// 			FromAccountNumber: accountSender,
// 			ToAccountNumber:   accountReciever,
// 			Currency:          "USD",
// 			Amount:            float64(rand.Intn(50) + 10),
// 		}

// 		reqs = append(reqs, t)
// 	}

// 	adapter.SendTransferMultiple(ctx, reqs)
// }

func runUnaryResiliency(ctx context.Context, adapter *resilliency.ResilliencyAdapter, minDelaySecond int32,
	maxDelaySecond int32, statusCodes []uint32) {
	res, err := adapter.UnaryResilliency(ctx, minDelaySecond, maxDelaySecond, statusCodes)

	if err != nil {
		log.Fatal().Msg("Failed to call UnaryResiliency: " + err.Error())
	}

	log.Info().Msgf("UnaryResiliency response: %v", res)
}

func runClientStreamResiliency(ctx context.Context, adapter *resilliency.ResilliencyAdapter, minDelaySecond int32,
	maxDelaySecond int32, statusCodes []uint32, count int) {
	adapter.ClientStreamResilliency(ctx, minDelaySecond, maxDelaySecond, statusCodes, count)
}

func runServerStreamingResilliency(ctx context.Context, adapter *resilliency.ResilliencyAdapter, minDelaySecond int32,
	maxDelaySecond int32, statusCodes []uint32) {
	adapter.ServerStreamResilliency(ctx, minDelaySecond, maxDelaySecond, statusCodes)
}
