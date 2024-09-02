package config

import (
	"github.com/fajaramaulana/go-grpc-micro-bank-server/exception"

	"github.com/gofiber/fiber/v2"
)

func NewFiberConfig() fiber.Config {
	return fiber.Config{
		ErrorHandler: exception.ErrorHandler,
		BodyLimit:    50 * 1024 * 1024,
	}
}
