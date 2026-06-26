package gameconfig

import (
	types "fish_shooting_admin_backend/pkg/model"
	"fish_shooting_admin_backend/pkg/responses"

	"github.com/jmoiron/sqlx"
)

type GameConfigServiceCreator interface {
	GetGameConfig() (*GetGameConfigResponse, *responses.ErrorResponse)
	UpdateGameConfig(req UpdateGameConfigRequest) (*UpdateGameConfigResponse, *responses.ErrorResponse)
}

type GameConfigService struct {
	DBPool        *sqlx.DB
	GameConfigRepo GameConfigRepo
	UserContext    *types.UserContext
}

func NewGameConfigService(dbPool *sqlx.DB, userContext *types.UserContext) GameConfigServiceCreator {
	return &GameConfigService{
		DBPool:         dbPool,
		GameConfigRepo: NewGameConfigRepoImpl(dbPool, userContext),
		UserContext:    userContext,
	}
}

func (s *GameConfigService) GetGameConfig() (*GetGameConfigResponse, *responses.ErrorResponse) {
	return s.GameConfigRepo.GetGameConfig()
}

func (s *GameConfigService) UpdateGameConfig(req UpdateGameConfigRequest) (*UpdateGameConfigResponse, *responses.ErrorResponse) {

	return s.GameConfigRepo.UpdateGameConfig(req)
}
