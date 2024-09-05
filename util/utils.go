package util

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ReqId(chId string) string {
	reqId := chId + time.Now().Format("20060102150405")
	return reqId
}

func ParseIntToRupiah(angka int) string {
	rupiah := ""
	angkaRev := reverseInt(angka)
	for i := 0; i < len(angkaRev); i++ {
		if i%3 == 0 && i != 0 {
			rupiah += "."
		}
		rupiah += string(angkaRev[i])
	}
	return "Rp " + reverseString(rupiah)
}

func reverseInt(angka int) string {
	angkaStr := string(angka)
	angkaRev := ""
	for i := len(angkaStr) - 1; i >= 0; i-- {
		angkaRev += string(angkaStr[i])
	}
	return angkaRev
}

func reverseString(str string) string {
	strRev := ""
	for i := len(str) - 1; i >= 0; i-- {
		strRev += string(str[i])
	}
	return strRev
}

func FormatRupiah(amount float64) string {
	strAmount := fmt.Sprintf("%.2f", amount)
	parts := strings.Split(strAmount, ".")
	intPart := parts[0]
	decimalPart := parts[1]

	// Sisipkan tanda titik setiap 3 digit dari belakang
	n := len(intPart)
	if n > 3 {
		for i := n - 3; i > 0; i -= 3 {
			intPart = intPart[:i] + "." + intPart[i:]
		}
	}

	// Gabungkan bagian integer dan desimal dengan tanda koma
	rupiah := "Rp " + intPart + "," + decimalPart
	return rupiah
}

// unaryInterceptorWithLogging wraps the retry interceptor and adds logging for each retry
func UnaryInterceptorWithLogging(retryInterceptor grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		var attempt uint = 0

		// Wrap the invoker with retry logic and add logging
		return retryInterceptor(
			ctx,
			method,
			req,
			reply,
			cc,
			func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				err := invoker(ctx, method, req, reply, cc, opts...)
				if status.Code(err) != codes.OK {
					attempt++
					log.Warn().Msgf("Retrying %s, attempt #%d, error: %v", method, attempt, err)
				}
				return err
			},
			opts...,
		)
	}
}
