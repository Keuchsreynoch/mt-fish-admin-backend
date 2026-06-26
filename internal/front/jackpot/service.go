package jackpot

import (
	types "fish_shooting_admin_backend/pkg/model"
	"fish_shooting_admin_backend/pkg/responses"

	"github.com/jmoiron/sqlx"
)

type JackpotServiceCreator interface {
	GetCurrentJackpotGlobal() (*JackpotGlobalResponse, *responses.ErrorResponse)
	GetJackpotGlobalUpdateForm() (*JackpotGlobalResponse, *responses.ErrorResponse)
	GetAllCompanyTopups(req JackpotCompanyTopupShowRequest) (*JackpotCompanyTopupListResponse, *responses.ErrorResponse)
	CreateCompanyTopup(req CreateJackpotCompanyTopupRequest) (*JackpotCompanyTopupResponse, *responses.ErrorResponse)
	GetAllLedgers(req JackpotLedgerShowRequest) (*JackpotLedgerListResponse, *responses.ErrorResponse)
	UpdateJackpotRate(req UpdateJackpotRateRequest) (*UpdateJackpotRateResponse, *responses.ErrorResponse)
	UpdateJackpotGlobal(req UpdateJackpotGlobalRequest) (*JackpotGlobalResponse, *responses.ErrorResponse)
}

type JackpotService struct {
	DBPool      *sqlx.DB
	JackpotRepo JackpotRepo
	UserContext *types.UserContext
}

func NewJackpotService(dbPool *sqlx.DB, userContext *types.UserContext) JackpotServiceCreator {
	return &JackpotService{
		DBPool:      dbPool,
		JackpotRepo: NewJackpotRepoImpl(dbPool),
		UserContext: userContext,
	}
}

func (s *JackpotService) GetCurrentJackpotGlobal() (*JackpotGlobalResponse, *responses.ErrorResponse) {
	return s.JackpotRepo.GetCurrentJackpotGlobal()
}

func (s *JackpotService) GetJackpotGlobalUpdateForm() (*JackpotGlobalResponse, *responses.ErrorResponse) {
	return s.JackpotRepo.GetCurrentJackpotGlobal()
}

func (s *JackpotService) GetAllCompanyTopups(req JackpotCompanyTopupShowRequest) (*JackpotCompanyTopupListResponse, *responses.ErrorResponse) {
	return s.JackpotRepo.GetAllCompanyTopups(req)
}

func (s *JackpotService) CreateCompanyTopup(req CreateJackpotCompanyTopupRequest) (*JackpotCompanyTopupResponse, *responses.ErrorResponse) {
	if s.UserContext == nil {
		zero := types.UserContext{}
		return s.JackpotRepo.CreateCompanyTopup(req, zero)
	}

	return s.JackpotRepo.CreateCompanyTopup(req, *s.UserContext)
}

func (s *JackpotService) GetAllLedgers(req JackpotLedgerShowRequest) (*JackpotLedgerListResponse, *responses.ErrorResponse) {
	return s.JackpotRepo.GetAllLedgers(req)
}

func (s *JackpotService) UpdateJackpotRate(req UpdateJackpotRateRequest) (*UpdateJackpotRateResponse, *responses.ErrorResponse) {
	if s.UserContext == nil {
		zero := types.UserContext{}
		return s.JackpotRepo.UpdateJackpotRate(req, zero)
	}

	return s.JackpotRepo.UpdateJackpotRate(req, *s.UserContext)
}

func (s *JackpotService) UpdateJackpotGlobal(req UpdateJackpotGlobalRequest) (*JackpotGlobalResponse, *responses.ErrorResponse) {
	if s.UserContext == nil {
		zero := types.UserContext{}
		return s.JackpotRepo.UpdateJackpotGlobal(req, zero)
	}

	return s.JackpotRepo.UpdateJackpotGlobal(req, *s.UserContext)
}
