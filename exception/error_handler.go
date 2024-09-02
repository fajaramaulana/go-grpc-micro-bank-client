package exception

import (
	"github.com/fajaramaulana/go-grpc-micro-bank-server/model"

	"github.com/gofiber/fiber/v2"
)

//001 = Database
//002 = General

func ErrorHandler(ctx *fiber.Ctx, err error) error {

	_, databaseError := err.(DatabaseError)
	if databaseError {
		return ctx.Status(500).JSON(model.WebResponse{
			StatusCode: 500,
			Message:    "Terjadi kesalahan, silahkan ulangi beberapa saat lagi, code=[001]",
		})
	}

	_, dataNotFoundError := err.(DataNotFoundError)
	if dataNotFoundError {
		return ctx.Status(404).JSON(model.WebResponse{
			StatusCode: 404,
			Message:    err.Error(),
		})
	}

	_, generalError := err.(GeneralError)
	if generalError {
		return ctx.Status(400).JSON(model.WebResponse{
			StatusCode: 400,
			Message:    err.Error(),
		})
	}

	return ctx.Status(400).JSON(model.WebResponse{
		StatusCode: 400,
		Message:    "Terjadi kesalahan, silahkan ulangi beberapa saat lagi, code=[002]",
	})
}
