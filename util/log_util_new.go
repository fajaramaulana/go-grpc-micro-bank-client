package util

import (
	"log"
	"time"

	"github.com/fajaramaulana/go-grpc-micro-bank-server/util/constant"

	"github.com/gofiber/fiber/v2"
)

// 0XXX -> SERVICE
// 1XXX -> REPOSITORY
// 2XXX -> VALIDATION
// 3XXX -> CONTROLLER
// 4XXX -> MAIN / JOB

// [INFO] ... [START]
func LogRequestNew(ctx *fiber.Ctx, requestData, uuidTr, uuidTc string) time.Time {
	log.Printf("[INFO]\tTR-%s\t%s\t%s\t%s\t%s\tTC-%s\t%s\t[START]",
		uuidTr, constant.APP_NAME, ctx.Route().Path, ctx.Method(), ctx.IP(), uuidTc, requestData)

	return time.Now()
}

// [INFO] ... [STOP]
func LogResponseNew(ctx *fiber.Ctx, responseData, uuidTr, uuidTc string, start time.Time) {
	log.Printf("[INFO]\tTR-%s\t%s\t%s\t%s\t%s\t%s\tTC-%s\t%s\t[STOP]",
		uuidTr, constant.APP_NAME, ctx.Route().Path, ctx.Method(), time.Since(start), ctx.IP(), uuidTc, responseData)
}

// [TRACE] | [INFO] | [WARN] | [ERROR] | [FATAL]
func Logging(ctx *fiber.Ctx, level, functionName, uniqueCode, transactionId, traceId, notes string) {
	if level != "DEBUG" {
		log.Printf("[%s]\tTR-%s\t%s\t%s\t%s\t%s\tTC-%s\t[%s]\t[%s]\t(%s)",
			level, transactionId, constant.APP_NAME, ctx.Route().Path, ctx.Method(), ctx.IP(), traceId, functionName, uniqueCode, notes)
	} else {
		log.Printf("[%s]\tTR-%s\t%s\t%s\t%s\t%s\tTC-%s\t[%s]\t[%s]\t(%s)",
			"OTHER", transactionId, constant.APP_NAME, ctx.Route().Path, ctx.Method(), ctx.IP(), traceId, functionName, uniqueCode, notes)
	}
}

// [DEBUG]
func LogDebug(ctx *fiber.Ctx, functionName, uniqueCode, transactionId, traceId, notes, data string) {
	log.Printf("[DEBUG]\tTR-%s\t%s\t%s\t%s\t%s\tTC-%s\t[%s]\t[%s]\t(%s)\t%s",
		transactionId, constant.APP_NAME, ctx.Route().Path, ctx.Method(), ctx.IP(), traceId, functionName, uniqueCode, notes, data)
}
