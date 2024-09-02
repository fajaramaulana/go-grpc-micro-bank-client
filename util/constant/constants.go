package constant

import (
	"github.com/gofiber/fiber/v2"
)

const (
	APP_NAME = "GO-GRPC-MICRO-BANK-SERVER"
	POST     = fiber.MethodPost
	GET      = fiber.MethodGet
	PUT      = fiber.MethodPut
	DELETE   = fiber.MethodDelete

	ContentType            = fiber.HeaderContentType
	Accept                 = fiber.HeaderAccept
	Authorization          = fiber.HeaderAuthorization
	ApplicationJSON        = fiber.MIMEApplicationJSON
	ApplicationTextXMLUtf8 = fiber.MIMETextXMLCharsetUTF8
	ApplicationXForwarded  = fiber.HeaderXForwardedFor

	AsiaJakartaZone = "Asia/Jakarta"

	TRACE = "TRACE" // Developers: Do I need to log states of variables? (NO)
	DEBUG = "DEBUG" // Developers: Do I need to log states of variables? (YES)
	INFO  = "INFO"  // Operationals: Do I log because of an unwanted state? NO
	WARN  = "WARN"  // Operationals: Do I log because of an unwanted state? (YES) && Can the process continue with the unwanted state? (YES)
	ERROR = "ERROR" // Operationals: Do I log because of an unwanted state? (YES) && Can the process continue with the unwanted state? (NO) && Can the application continue with the unwanted state? (YES)
	FATAL = "FATAL" // Operationals: Do I log because of an unwanted state? (YES) && Can the process continue with the unwanted state? (NO) && Can the application continue with the unwanted state? (NO

	FuncGetDataPIL          = "GET DATA PIL="
	FuncUpdateInsertDataPIL = "INSERT/UPDATE DATA PIL"
)
