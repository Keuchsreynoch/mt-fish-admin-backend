package jackpothistory

import (
	types "fish_shooting_admin_backend/pkg/model"
	"fish_shooting_admin_backend/pkg/responses"

	"github.com/jmoiron/sqlx"
)

type JackpotHistoryServiceCreator interface {
	GetAllJackpotHistories(req JackpotHistoryShowRequest) (*JackpotHistoryListResponse, *responses.ErrorResponse)
}

type JackpotHistoryService struct {
	DBPool             *sqlx.DB
	JackpotHistoryRepo JackpotHistoryRepo
}

func NewJackpotHistoryService(dbPool *sqlx.DB, userContext *types.UserContext) JackpotHistoryServiceCreator {
	_ = userContext
	return &JackpotHistoryService{
		DBPool:             dbPool,
		JackpotHistoryRepo: NewJackpotHistoryRepoImpl(dbPool),
	}
}

func (s *JackpotHistoryService) GetAllJackpotHistories(req JackpotHistoryShowRequest) (*JackpotHistoryListResponse, *responses.ErrorResponse) {
	return s.JackpotHistoryRepo.GetAllJackpotHistories(req)
}
