package fishtype

import (
	"fish_shooting_admin_backend/pkg/responses"

	"github.com/jmoiron/sqlx"
)

type FishTypeServiceCreator interface {
	GetFishTypes() ([]FishInfoResponse, *responses.ErrorResponse)
}

type FishTypeService struct {
	DBPool       *sqlx.DB
	FishTypeRepo FishTypeRepo
}

func NewFishTypeService(dbPool *sqlx.DB) FishTypeServiceCreator {
	return &FishTypeService{
		DBPool:       dbPool,
		FishTypeRepo: NewFishTypeRepoImpl(dbPool),
	}
}

func (s *FishTypeService) GetFishTypes() ([]FishInfoResponse, *responses.ErrorResponse) {
	return s.FishTypeRepo.GetFishTypes()
}
