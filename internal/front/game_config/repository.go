package gameconfig

import (
	"database/sql"
	custom_log "fish_shooting_admin_backend/pkg/logs"
	types "fish_shooting_admin_backend/pkg/model"
	"fish_shooting_admin_backend/pkg/responses"
	"fish_shooting_admin_backend/pkg/utils"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type GameConfigRepo interface {
	GetGameConfig() (*GetGameConfigResponse, *responses.ErrorResponse)
	UpdateGameConfig(req UpdateGameConfigRequest) (*UpdateGameConfigResponse, *responses.ErrorResponse)
}

type GameConfigRepoImpl struct {
	DBPool      *sqlx.DB
	UserContext *types.UserContext
}

func NewGameConfigRepoImpl(dbPool *sqlx.DB, userContext *types.UserContext) GameConfigRepo {
	return &GameConfigRepoImpl{
		DBPool:      dbPool,
		UserContext: userContext,
	}
}

func (r *GameConfigRepoImpl) GetGameConfig() (*GetGameConfigResponse, *responses.ErrorResponse) {
	row := gameConfigRow{}
	if err := r.DBPool.Get(&row, `
		SELECT
			g.id,
			g.game_code,
			g.game_name,
			COALESCE(g.rtp_target, 0) AS rtp_target,
			COALESCE(g.rtp_floor, 0) AS rtp_floor,
			COALESCE(g.rtp_ceiling, 0) AS rtp_ceiling,
			COALESCE(g.jackpot_rate, 0) AS jackpot_rate,
			COALESCE(g.company_profit_rate, 0) AS company_profit_rate,
			COALESCE(g.max_bullet_per_second, 0) AS max_bullet_per_second,
			g.allow_auto_fire,
			g.allow_auto_lock,
			COALESCE(g.kill_rate_min, 0) AS kill_rate_min,
			COALESCE(g.kill_rate_max, 0) AS kill_rate_max,
			COALESCE(g.rtp_adjust_strength, 0) AS rtp_adjust_strength,
			g.status_id,
			g."order",
			g.created_at,
			COALESCE(g.created_by, 0) AS created_by,
			COALESCE(g.updated_at, g.created_at) AS updated_at,
			COALESCE(g.updated_by, g.created_by, 0) AS updated_by,
			COALESCE(u.user_name, '') AS updated_by_username
		FROM tbl_game_configs g
		LEFT JOIN tbl_users u
			ON u.id = COALESCE(g.updated_by, g.created_by, 0)
		WHERE g.deleted_at IS NULL
		ORDER BY g.id ASC
		LIMIT 1
	`); err != nil {
		custom_log.NewCustomLog("get_game_config_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		if err == sql.ErrNoRows {
			return nil, e.NewErrorResponse("not_found", fmt.Errorf("not_found"))
		}
		return nil, e.NewErrorResponse("get_game_config_failed", fmt.Errorf("error_database"))
	}

	return &GetGameConfigResponse{
		GameName:          row.GameName,
		RtpTarget:         row.RtpTarget.StringFixed(4),
		RtpFloor:          row.RtpFloor.StringFixed(4),
		RtpCeiling:        row.RtpCeiling.StringFixed(4),
		JackpotRate:       row.JackpotRate.StringFixed(6),
		CompanyProfitRate: row.CompanyProfitRate.StringFixed(6),
		StatusID:          row.StatusID,
		UpdatedByUsername: row.UpdatedByUsername,
		UpdatedAt:         row.UpdatedAt,
	}, nil
}

func (r *GameConfigRepoImpl) UpdateGameConfig(req UpdateGameConfigRequest) (*UpdateGameConfigResponse, *responses.ErrorResponse) {
	jackpotRate, err := decimal.NewFromString(strings.TrimSpace(req.JackpotRate))
	if err != nil {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("invalid_jackpot_rate", err)
	}
	companyProfitRate, err := decimal.NewFromString(strings.TrimSpace(req.CompanyProfitRate))
	if err != nil {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("invalid_company_profit_rate", err)
	}
	rtpTarget, err := decimal.NewFromString(strings.TrimSpace(req.RtpTarget))
	if err != nil {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("invalid_rtp_target", err)
	}
	rtpFloor, err := decimal.NewFromString(strings.TrimSpace(req.RtpFloor))
	if err != nil {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("invalid_rtp_floor", err)
	}
	rtpCeiling, err := decimal.NewFromString(strings.TrimSpace(req.RtpCeiling))
	if err != nil {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("invalid_rtp_ceiling", err)
	}

	if jackpotRate.LessThanOrEqual(decimal.Zero) || jackpotRate.GreaterThan(decimal.NewFromInt(1)) {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("jackpot_rate_must_be_between_0_and_1", fmt.Errorf("jackpot rate must be between 0 and 1"))
	}
	if companyProfitRate.LessThan(decimal.Zero) || companyProfitRate.GreaterThan(decimal.NewFromInt(1)) {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("company_profit_rate_must_be_between_0_and_1", fmt.Errorf("company profit rate must be between 0 and 1"))
	}
	if rtpFloor.LessThanOrEqual(decimal.Zero) || rtpFloor.GreaterThan(decimal.NewFromInt(1)) {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("rtp_floor_must_be_between_0_and_1", fmt.Errorf("rtp floor must be between 0 and 1"))
	}
	if rtpCeiling.LessThanOrEqual(decimal.Zero) || rtpCeiling.GreaterThan(decimal.NewFromInt(1)) {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("rtp_ceiling_must_be_between_0_and_1", fmt.Errorf("rtp ceiling must be between 0 and 1"))
	}
	if rtpTarget.LessThanOrEqual(decimal.Zero) || rtpTarget.GreaterThan(decimal.NewFromInt(1)) {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("rtp_target_must_be_between_0_and_1", fmt.Errorf("rtp target must be between 0 and 1"))
	}
	if rtpTarget.Add(jackpotRate).Add(companyProfitRate).GreaterThan(decimal.NewFromInt(1)) {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("rtp_target_jackpot_rate_company_profit_rate_must_be_less_than_or_equal_to_1", fmt.Errorf("rtp target, jackpot rate, and company profit rate must be less than or equal to 1"))
	}
	if rtpFloor.GreaterThan(rtpCeiling) {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("rtp_floor_must_be_less_than_or_equal_to_rtp_ceiling", fmt.Errorf("rtp floor must be less than or equal to rtp ceiling"))
	}
	if rtpTarget.LessThan(rtpFloor) || rtpTarget.GreaterThan(rtpCeiling) {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("rtp_target_must_be_between_floor_and_ceiling", fmt.Errorf("rtp target must be between floor and ceiling"))
	}
	if req.StatusID != 1 && req.StatusID != 2 {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("status_id_must_be_1_or_2", fmt.Errorf("status id must be 1 or 2"))
	}

	tx, err := r.DBPool.Beginx()
	if err != nil {
		custom_log.NewCustomLog("update_game_config_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("update_game_config_failed", fmt.Errorf("error_database"))
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	now, err := utils.GetCurrentAppTime()
	if err != nil {
		custom_log.NewCustomLog("update_game_config_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("update_game_config_failed", fmt.Errorf("error_database"))
	}

	row := gameConfigRow{}
	if err := tx.Get(&row, `
		SELECT
			g.id
		FROM tbl_game_configs g
		WHERE g.deleted_at IS NULL
		ORDER BY g.id ASC
		LIMIT 1
		FOR UPDATE
	`); err != nil {
		custom_log.NewCustomLog("update_game_config_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		if err == sql.ErrNoRows {
			return nil, e.NewErrorResponse("not_found", fmt.Errorf("not_found"))
		}
		return nil, e.NewErrorResponse("update_game_config_failed", fmt.Errorf("error_database"))
	}

	if _, err := tx.Exec(`
		UPDATE tbl_game_configs
		SET jackpot_rate = $1,
			company_profit_rate = $2,
			rtp_target = $3,
			rtp_floor = $4,
			rtp_ceiling = $5,
			status_id = $6,
			updated_at = $7,
			updated_by = $8
		WHERE id = $9
	`, jackpotRate, companyProfitRate, rtpTarget, rtpFloor, rtpCeiling, req.StatusID, now, r.UserContext.Id, row.ID); err != nil {
		custom_log.NewCustomLog("update_game_config_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("update_game_config_failed", fmt.Errorf("error_database"))
	}

	if err := tx.Commit(); err != nil {
		custom_log.NewCustomLog("update_game_config_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("update_game_config_failed", fmt.Errorf("error_database"))
	}
	tx = nil

	return &UpdateGameConfigResponse{
		ID:                row.ID,
		JackpotRate:       jackpotRate.StringFixed(6),
		CompanyProfitRate: companyProfitRate.StringFixed(6),
		RtpTarget:         rtpTarget.StringFixed(4),
		RtpFloor:          rtpFloor.StringFixed(4),
		RtpCeiling:        rtpCeiling.StringFixed(4),
		StatusID:          req.StatusID,
		UpdatedAt:         now,
		UpdatedBy:         &r.UserContext.Id,
		UpdatedByUsername: r.UserContext.UserName,
	}, nil
}
