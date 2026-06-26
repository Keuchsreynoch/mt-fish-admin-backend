package cointransaction

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type CoinTransactionRoute struct {
	App                    *fiber.App
	DBPool                 *sqlx.DB
	CoinTransactionHandler *CoinTransactionHandlerImpl
}

func NewCoinTransactionRoute(app *fiber.App, dbPool *sqlx.DB) *CoinTransactionRoute {
	return &CoinTransactionRoute{
		App:                    app,
		DBPool:                 dbPool,
		CoinTransactionHandler: NewCoinTransactionHandler(dbPool),
	}
}

func (r *CoinTransactionRoute) RegisterCoinTransactionRoute() *CoinTransactionRoute {
	adminCoinTransaction := r.App.Group("/api/v1/admin/coin-transactions")
	adminCoinTransaction.Get("/", r.CoinTransactionHandler.GetTransactions)
	return r
}
