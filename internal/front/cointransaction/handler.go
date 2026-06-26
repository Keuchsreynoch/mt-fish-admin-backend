package cointransaction

import (
	"errors"
	response "fish_shooting_admin_backend/pkg/http/response"
	"fish_shooting_admin_backend/pkg/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type CoinTransactionHandler interface {
	GetTransactions(c *fiber.Ctx) error
}

type CoinTransactionHandlerImpl struct {
	service func(c *fiber.Ctx) CoinTransactionServiceCreator
	v       *utils.Validator
}

func NewCoinTransactionHandler(db *sqlx.DB) *CoinTransactionHandlerImpl {
	return &CoinTransactionHandlerImpl{
		service: func(c *fiber.Ctx) CoinTransactionServiceCreator { return NewCoinTransactionService(db) },
		v:       utils.NewValidator(),
	}
}

func (h *CoinTransactionHandlerImpl) GetTransactions(c *fiber.Ctx) error {
	var req CoinTransactionShowRequest
	if err := req.Bind(c, h.v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("get_coin_transactions_failed", nil, c),
				-1000,
				err,
			),
		)
	}

	resp, err := h.service(c).GetTransactions(req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate(err.MessageID, nil, c),
				-1000,
				errors.New(utils.Translate(err.Err.Error(), nil, c)),
			),
		)
	}

	return c.Status(http.StatusOK).JSON(
		response.NewResponseWithPaging(
			utils.Translate("get_coin_transactions_success", nil, c),
			1000,
			resp,
			req.PageOptions.Page,
			req.PageOptions.Perpage,
			resp.Total,
		),
	)
}
