package memberbonus

import (
	types "fish_shooting_admin_backend/pkg/model"
	"fish_shooting_admin_backend/pkg/responses"

	"github.com/jmoiron/sqlx"
)

type JackpotMemberBonusServiceCreator interface {
	GetAllMemberBonuses(req JackpotMemberBonusShowRequest) (*JackpotMemberBonusListResponse, *responses.ErrorResponse)
	CreateMemberBonus(req CreateJackpotMemberBonusRequest) (*JackpotMemberBonusResponse, *responses.ErrorResponse)
}

type JackpotMemberBonusService struct {
	DBPool                 *sqlx.DB
	JackpotMemberBonusRepo JackpotMemberBonusRepo
	UserContext            *types.UserContext
}

func NewJackpotMemberBonusService(dbPool *sqlx.DB, userContext *types.UserContext) JackpotMemberBonusServiceCreator {
	return &JackpotMemberBonusService{
		DBPool:                 dbPool,
		JackpotMemberBonusRepo: NewJackpotMemberBonusRepoImpl(dbPool),
		UserContext:            userContext,
	}
}

func (s *JackpotMemberBonusService) GetAllMemberBonuses(req JackpotMemberBonusShowRequest) (*JackpotMemberBonusListResponse, *responses.ErrorResponse) {
	return s.JackpotMemberBonusRepo.GetAllMemberBonuses(req)
}

func (s *JackpotMemberBonusService) CreateMemberBonus(req CreateJackpotMemberBonusRequest) (*JackpotMemberBonusResponse, *responses.ErrorResponse) {
	if s.UserContext == nil {
		zero := types.UserContext{}
		return s.JackpotMemberBonusRepo.CreateMemberBonus(req, zero)
	}

	return s.JackpotMemberBonusRepo.CreateMemberBonus(req, *s.UserContext)
}
