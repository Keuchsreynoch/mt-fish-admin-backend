package gameconfig

import (
	"errors"
	custom_log "fish_shooting_admin_backend/pkg/logs"
	"fish_shooting_admin_backend/pkg/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

type UpdateGameConfigRequest struct {
	JackpotRate string `json:"jackpot_rate" validate:"required"`
	RtpTarget   string `json:"rtp_target" validate:"required"`
	RtpFloor    string `json:"rtp_floor" validate:"required"`
	RtpCeiling  string `json:"rtp_ceiling" validate:"required"`
	StatusID    int    `json:"status_id" validate:"required"`
}

type UpdateGameConfigResponse struct {
	ID                int64     `json:"id"`
	JackpotRate       string    `json:"jackpot_rate"`
	RtpTarget         string    `json:"rtp_target"`
	RtpFloor          string    `json:"rtp_floor"`
	RtpCeiling        string    `json:"rtp_ceiling"`
	StatusID          int       `json:"status_id"`
	UpdatedAt         time.Time `json:"updated_at"`
	UpdatedBy         *int      `json:"updated_by"`
	UpdatedByUsername string    `json:"updated_by_username"`
}

type GetGameConfigResponse struct {
	GameName          string `json:"game_name"`
	RtpTarget         string `json:"rtp_target"`
	RtpFloor          string `json:"rtp_floor"`
	RtpCeiling        string `json:"rtp_ceiling"`
	JackpotRate       string `json:"jackpot_rate"`
	StatusID          int    `json:"status_id"`
	UpdatedByUsername string `json:"updated_by_username"`
}

type gameConfigRow struct {
	ID                 int64           `db:"id"`
	GameCode           string          `db:"game_code"`
	GameName           string          `db:"game_name"`
	RtpTarget          decimal.Decimal `db:"rtp_target"`
	RtpFloor           decimal.Decimal `db:"rtp_floor"`
	RtpCeiling         decimal.Decimal `db:"rtp_ceiling"`
	JackpotRate        decimal.Decimal `db:"jackpot_rate"`
	MaxBulletPerSecond decimal.Decimal `db:"max_bullet_per_second"`
	AllowAutoFire      bool            `db:"allow_auto_fire"`
	AllowAutoLock      bool            `db:"allow_auto_lock"`
	KillRateMin        decimal.Decimal `db:"kill_rate_min"`
	KillRateMax        decimal.Decimal `db:"kill_rate_max"`
	RtpAdjustStrength  decimal.Decimal `db:"rtp_adjust_strength"`
	StatusID           int             `db:"status_id"`
	Order              int             `db:"order"`
	CreatedAt          time.Time       `db:"created_at"`
	CreatedBy          int             `db:"created_by"`
	UpdatedAt          time.Time       `db:"updated_at"`
	UpdatedBy          int             `db:"updated_by"`
	UpdatedByUsername  string          `db:"updated_by_username"`
}

func (r *UpdateGameConfigRequest) Bind(c *fiber.Ctx, v *utils.Validator) error {
	if err := c.BodyParser(r); err != nil {
		custom_log.NewCustomLog("update_game_config_failed", err.Error(), "error")
		return errors.New(utils.Translate("invalid_body", nil, c))
	}

	r.JackpotRate = strings.TrimSpace(r.JackpotRate)
	r.RtpTarget = strings.TrimSpace(r.RtpTarget)
	r.RtpFloor = strings.TrimSpace(r.RtpFloor)
	r.RtpCeiling = strings.TrimSpace(r.RtpCeiling)

	if err := v.Validate(r, c); err != nil {
		return err
	}

	if r.StatusID != 1 && r.StatusID != 2 {
		return errors.New("status_id must be 1 (natural mode) or 2 (strich mode)")
	}

	return nil
}
