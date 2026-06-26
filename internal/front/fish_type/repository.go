package fishtype

import (
	"fmt"

	custom_log "fish_shooting_admin_backend/pkg/logs"
	"fish_shooting_admin_backend/pkg/responses"

	"github.com/jmoiron/sqlx"
)

type FishTypeRepo interface {
	GetFishTypes() ([]FishInfoResponse, *responses.ErrorResponse)
}

type FishTypeRepoImpl struct {
	DBPool *sqlx.DB
}

func NewFishTypeRepoImpl(dbPool *sqlx.DB) FishTypeRepo {
	return &FishTypeRepoImpl{DBPool: dbPool}
}

func (r *FishTypeRepoImpl) GetFishTypes() ([]FishInfoResponse, *responses.ErrorResponse) {
	return r.getFishTypes()
}

func (r *FishTypeRepoImpl) getFishTypes() ([]FishInfoResponse, *responses.ErrorResponse) {
	fishRows := []manifestFishRow{}
	if err := r.DBPool.Select(&fishRows, `
		SELECT
			fish_type_name,
			is_boss,
			NULL::text AS boss_name,
			min_kill_odd::float8 AS min_kill_odd,
			max_kill_odd::float8 AS max_kill_odd,
			miss_reward_enabled,
			min_miss_reward_odd::float8 AS min_miss_reward_odd,
			max_miss_reward_odd::float8 AS max_miss_reward_odd,
			base_speed::float8 AS base_speed
		FROM tbl_fish_types
		WHERE deleted_at IS NULL
		  AND status_id = $1
		ORDER BY id
	`, ActiveStatusID); err != nil {
		custom_log.NewCustomLog("get_fish_types_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_fish_types_failed", fmt.Errorf("error_database"))
	}

	result := make([]FishInfoResponse, 0, len(fishRows))
	for _, row := range fishRows {
		result = append(result, FishInfoResponse{
			FishTypeName:      row.FishTypeName,
			IsBoss:            row.IsBoss,
			BossName:          row.BossName,
			MinKillOdd:        row.MinKillOdd,
			MaxKillOdd:        row.MaxKillOdd,
			BaseSpeed:         row.BaseSpeed,
			MissRewardEnabled: row.MissRewardEnabled,
			MinMissRewardOdd:  row.MinMissRewardOdd,
			MaxMissRewardOdd:  row.MaxMissRewardOdd,
		})
	}

	return result, nil
}
