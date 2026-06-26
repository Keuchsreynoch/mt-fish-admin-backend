package cointransaction

import (
	"fish_shooting_admin_backend/pkg/responses"

	"github.com/jmoiron/sqlx"
)

type CoinTransactionServiceCreator interface {
	GetTransactions(req CoinTransactionShowRequest) (*CoinTransactionListResponse, *responses.ErrorResponse)
}

type CoinTransactionService struct {
	DBPool              *sqlx.DB
	CoinTransactionRepo CoinTransactionRepo
}

func NewCoinTransactionService(dbPool *sqlx.DB) CoinTransactionServiceCreator {
	return &CoinTransactionService{
		DBPool:              dbPool,
		CoinTransactionRepo: NewCoinTransactionRepoImpl(dbPool),
	}
}

func (s *CoinTransactionService) GetTransactions(req CoinTransactionShowRequest) (*CoinTransactionListResponse, *responses.ErrorResponse) {
	return s.CoinTransactionRepo.GetTransactions(req)
}
