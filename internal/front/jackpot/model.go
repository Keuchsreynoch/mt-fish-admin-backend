package jackpot

import (
	"errors"
	custom_log "fish_shooting_admin_backend/pkg/logs"
	types "fish_shooting_admin_backend/pkg/model"
	"fish_shooting_admin_backend/pkg/utils"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

var jackpotFilterPattern = regexp.MustCompile(`^filters\[(\d+)\]\[(property|operator|value)\](?:\[(\d+)\])?$`)

type JackpotCurrentRequest struct{}

type JackpotCompanyTopupShowRequest struct {
	PageOptions types.Paging   `json:"paging_options" query:"paging_options" validate:"required"`
	Sorts       []types.Sort   `json:"sorts,omitempty" query:"sorts"`
	Filters     []types.Filter `json:"filters,omitempty" query:"filters"`
}

type JackpotLedgerShowRequest struct {
	PageOptions types.Paging   `json:"paging_options" query:"paging_options" validate:"required"`
	Sorts       []types.Sort   `json:"sorts,omitempty" query:"sorts"`
	Filters     []types.Filter `json:"filters,omitempty" query:"filters"`
}

type CreateJackpotCompanyTopupRequest struct {
	Amount string `json:"amount" validate:"required"`
	Note   string `json:"note"`
}

type UpdateJackpotRateRequest struct {
	GameCode   string `json:"game_code" validate:"required,min=1,max=30"`
	JackpotRate string `json:"jackpot_rate" validate:"required"`
}

type UpdateJackpotGlobalRequest struct {
	ThresholdAmount          string `json:"threshold_amount" validate:"required"`
	ChanceDenom              int64  `json:"chance_denom" validate:"required,min=1"`
	PayoutPercent            string `json:"payout_percent" validate:"required"`
	MinEligibleBetAmount     string `json:"min_eligible_bet_amount" validate:"required"`
	JackpotFixedPayoutAmount string `json:"jackpot_fixed_payout_amount" validate:"required"`
}

func (r *JackpotCompanyTopupShowRequest) Bind(c *fiber.Ctx, v *utils.Validator) error {
	if err := c.QueryParser(r); err != nil {
		custom_log.NewCustomLog("get_jackpot_company_topups_failed", err.Error(), "error")
		return errors.New(utils.Translate("invalid_body", nil, c))
	}

	r.Filters = parseJackpotFilters(c)

	if err := v.Validate(r, c); err != nil {
		return err
	}

	return nil
}

func (r *JackpotLedgerShowRequest) Bind(c *fiber.Ctx, v *utils.Validator) error {
	if err := c.QueryParser(r); err != nil {
		custom_log.NewCustomLog("get_jackpot_ledger_failed", err.Error(), "error")
		return errors.New(utils.Translate("invalid_body", nil, c))
	}

	r.Filters = parseJackpotFilters(c)

	if err := v.Validate(r, c); err != nil {
		return err
	}

	return nil
}

func (r *CreateJackpotCompanyTopupRequest) Bind(c *fiber.Ctx, v *utils.Validator) error {
	if err := c.BodyParser(r); err != nil {
		custom_log.NewCustomLog("create_jackpot_company_topup_failed", err.Error(), "error")
		return errors.New(utils.Translate("invalid_body", nil, c))
	}

	if err := v.Validate(r, c); err != nil {
		return err
	}

	return nil
}

func (r *UpdateJackpotRateRequest) Bind(c *fiber.Ctx, v *utils.Validator) error {
	if err := c.BodyParser(r); err != nil {
		custom_log.NewCustomLog("update_jackpot_rate_failed", err.Error(), "error")
		return errors.New(utils.Translate("invalid_body", nil, c))
	}

	r.GameCode = strings.TrimSpace(strings.ToUpper(r.GameCode))
	r.JackpotRate = strings.TrimSpace(r.JackpotRate)

	if err := v.Validate(r, c); err != nil {
		return err
	}

	return nil
}

func (r *UpdateJackpotGlobalRequest) Bind(c *fiber.Ctx, v *utils.Validator) error {
	if err := c.BodyParser(r); err != nil {
		custom_log.NewCustomLog("update_jackpot_global_failed", err.Error(), "error")
		return errors.New(utils.Translate("invalid_body", nil, c))
	}

	r.ThresholdAmount = strings.TrimSpace(r.ThresholdAmount)
	r.PayoutPercent = strings.TrimSpace(r.PayoutPercent)
	r.MinEligibleBetAmount = strings.TrimSpace(r.MinEligibleBetAmount)
	r.JackpotFixedPayoutAmount = strings.TrimSpace(r.JackpotFixedPayoutAmount)

	if err := v.Validate(r, c); err != nil {
		return err
	}

	return nil
}

type JackpotGlobalResponse struct {
	ID                       int64      `json:"id"`
	CurrentAmount            string     `json:"current_amount"`
	ThresholdAmount          string     `json:"threshold_amount"`
	ChanceDenom              int64      `json:"chance_denom"`
	PayoutPercent            string     `json:"payout_percent"`
	MinEligibleBetAmount     string     `json:"min_eligible_bet_amount"`
	CompanyTopupAmount       string     `json:"company_topup_amount"`
	JackpotFixedPayoutAmount string     `json:"jackpot_fixed_payout_amount"`
	StatusID                 int16      `json:"status_id"`
	Order                    *int       `json:"order"`
	CreatedAt                time.Time  `json:"created_at"`
	CreatedBy                *int       `json:"created_by"`
	UpdatedAt                *time.Time `json:"updated_at"`
	UpdatedBy                *int       `json:"updated_by"`
	DeletedAt                *time.Time `json:"-"`
	DeletedBy                *int       `json:"-"`
}

type UpdateJackpotRateResponse struct {
	ID         int64      `json:"id"`
	GameCode   string     `json:"game_code"`
	JackpotRate string     `json:"jackpot_rate"`
	UpdatedAt  time.Time  `json:"updated_at"`
	UpdatedBy  *int       `json:"updated_by"`
}

type jackpotGlobalRow struct {
	ID                       int64           `db:"id"`
	CurrentAmount            decimal.Decimal `db:"current_amount"`
	ThresholdAmount          decimal.Decimal `db:"threshold_amount"`
	ChanceDenom              int64           `db:"chance_denom"`
	PayoutPercent            decimal.Decimal `db:"payout_percent"`
	MinEligibleBetAmount     decimal.Decimal `db:"min_eligible_bet_amount"`
	CompanyTopupAmount       decimal.Decimal `db:"company_topup_amount"`
	JackpotFixedPayoutAmount decimal.Decimal `db:"jackpot_fixed_payout_amount"`
	StatusID                 int16           `db:"status_id"`
	Order                    *int            `db:"order"`
	CreatedAt                time.Time       `db:"created_at"`
	CreatedBy                *int            `db:"created_by"`
	UpdatedAt                *time.Time      `db:"updated_at"`
	UpdatedBy                *int            `db:"updated_by"`
	DeletedAt                *time.Time      `db:"deleted_at"`
	DeletedBy                *int            `db:"deleted_by"`
}

type jackpotGameConfigRow struct {
	ID              int64           `db:"id"`
	GameCode        string          `db:"game_code"`
	JackpotRate     decimal.Decimal `db:"jackpot_rate"`
	UpdatedAt       *time.Time      `db:"updated_at"`
	UpdatedBy       *int            `db:"updated_by"`
	DeletedAt       *time.Time      `db:"deleted_at"`
	DeletedBy       *int            `db:"deleted_by"`
}

type JackpotCompanyTopupListResponse struct {
	Topups []JackpotCompanyTopupResponse `json:"topups"`
	Total  int                           `json:"-"`
}

type JackpotCompanyTopupResponse struct {
	ID                  int64      `json:"id"`
	Username            string     `json:"username"`
	Amount              string     `json:"amount"`
	CurrentAmountBefore string     `json:"current_amount_before"`
	CurrentAmountAfter  string     `json:"current_amount_after"`
	Note                string     `json:"note"`
	Order               int        `json:"order"`
	CreatedAt           time.Time  `json:"created_at"`
	CreatedBy           *int       `json:"created_by"`
	UpdatedAt           *time.Time `json:"-"`
	UpdatedBy           *int       `json:"-"`
	DeletedAt           *time.Time `json:"-"`
	DeletedBy           *int       `json:"-"`
}

type JackpotMemberBonusListResponse struct {
	Bonuses []JackpotMemberBonusResponse `json:"bonuses"`
	Total   int                          `json:"-"`
}

type JackpotMemberBonusResponse struct {
	ID                  int64      `json:"id"`
	JackpotGlobalID     int64      `json:"jackpot_global_id"`
	MemberID            int64      `json:"member_id"`
	MemberUUID          string     `json:"member_uuid"`
	UserName            string     `json:"user_name"`
	Amount              string     `json:"amount"`
	CurrentAmountBefore string     `json:"current_amount_before"`
	CurrentAmountAfter  string     `json:"current_amount_after"`
	Note                string     `json:"note"`
	Order               int        `json:"order"`
	CreatedAt           time.Time  `json:"created_at"`
	CreatedBy           *int       `json:"created_by"`
	UpdatedAt           *time.Time `json:"-"`
	UpdatedBy           *int       `json:"-"`
	DeletedAt           *time.Time `json:"-"`
	DeletedBy           *int       `json:"-"`
}

type jackpotMemberBonusRow struct {
	ID                  int64           `db:"id"`
	JackpotGlobalID     int64           `db:"jackpot_global_id"`
	MemberID            int64           `db:"member_id"`
	MemberUUID          string          `db:"member_uuid"`
	UserName            string          `db:"user_name"`
	Amount              decimal.Decimal `db:"amount"`
	CurrentAmountBefore decimal.Decimal `db:"current_amount_before"`
	CurrentAmountAfter  decimal.Decimal `db:"current_amount_after"`
	Note                string          `db:"note"`
	Order               int             `db:"order"`
	CreatedAt           time.Time       `db:"created_at"`
	CreatedBy           *int            `db:"created_by"`
	UpdatedAt           *time.Time      `db:"updated_at"`
	UpdatedBy           *int            `db:"updated_by"`
	DeletedAt           *time.Time      `db:"deleted_at"`
	DeletedBy           *int            `db:"deleted_by"`
}

type jackpotCompanyTopupRow struct {
	ID                  int64           `db:"id"`
	JackpotGlobalID     int64           `db:"jackpot_global_id"`
	Username            string          `db:"username"`
	Amount              decimal.Decimal `db:"amount"`
	CurrentAmountBefore decimal.Decimal `db:"current_amount_before"`
	CurrentAmountAfter  decimal.Decimal `db:"current_amount_after"`
	Note                string          `db:"note"`
	Order               int             `db:"order"`
	CreatedAt           time.Time       `db:"created_at"`
	CreatedBy           *int            `db:"created_by"`
	UpdatedAt           *time.Time      `db:"updated_at"`
	UpdatedBy           *int            `db:"updated_by"`
	DeletedAt           *time.Time      `db:"deleted_at"`
	DeletedBy           *int            `db:"deleted_by"`
}

type JackpotLedgerListResponse struct {
	Ledgers []JackpotLedgerResponse `json:"ledgers"`
	Total   int                     `json:"-"`
}

type JackpotLedgerResponse struct {
	ID                     int64      `json:"id"`
	JackpotGlobalID        int64      `json:"jackpot_global_id"`
	FishTypeID             int64      `json:"fish_type_id"`
	FishTypeName           string     `json:"fish_type_name"`
	MemberID               *int       `json:"member_id"`
	BetID                  *int64     `json:"bet_id"`
	TicketID               *int64     `json:"ticket_id"`
	StatementID            *int64     `json:"statement_id"`
	SourceType             string     `json:"source_type"`
	GlobalContributionCoin string     `json:"global_contribution_coin"`
	PoolBefore             string     `json:"pool_before"`
	PoolAfter              string     `json:"pool_after"`
	CreatedAt              time.Time  `json:"created_at"`
	CreatedBy              *int       `json:"created_by"`
	UpdatedAt              *time.Time `json:"-"`
	UpdatedBy              *int       `json:"-"`
	DeletedAt              *time.Time `json:"-"`
	DeletedBy              *int       `json:"-"`
}

type jackpotLedgerRow struct {
	ID                     int64           `db:"id"`
	JackpotGlobalID        int64           `db:"jackpot_global_id"`
	FishTypeID             int64           `db:"fish_type_id"`
	FishTypeName           string          `db:"fish_type_name"`
	MemberID               *int            `db:"member_id"`
	BetID                  *int64          `db:"bet_id"`
	TicketID               *int64          `db:"ticket_id"`
	StatementID            *int64          `db:"statement_id"`
	SourceType             string          `db:"source_type"`
	GlobalContributionCoin decimal.Decimal `db:"global_contribution_coin"`
	PoolBefore             decimal.Decimal `db:"pool_before"`
	PoolAfter              decimal.Decimal `db:"pool_after"`
	CreatedAt              time.Time       `db:"created_at"`
	CreatedBy              *int            `db:"created_by"`
	UpdatedAt              *time.Time      `db:"updated_at"`
	UpdatedBy              *int            `db:"updated_by"`
	DeletedAt              *time.Time      `db:"deleted_at"`
	DeletedBy              *int            `db:"deleted_by"`
}

type jackpotFilterDraft struct {
	property string
	operator string
	scalar   string
	values   map[int]string
}

func parseJackpotFilters(c *fiber.Ctx) []types.Filter {
	drafts := make(map[int]*jackpotFilterDraft)

	c.Context().QueryArgs().VisitAll(func(key, value []byte) {
		matches := jackpotFilterPattern.FindStringSubmatch(string(key))
		if len(matches) != 4 {
			return
		}

		index, err := strconv.Atoi(matches[1])
		if err != nil {
			return
		}

		field := matches[2]
		nestedIndex := matches[3]

		draft := drafts[index]
		if draft == nil {
			draft = &jackpotFilterDraft{
				values: make(map[int]string),
			}
			drafts[index] = draft
		}

		switch field {
		case "property":
			draft.property = string(value)
		case "operator":
			draft.operator = string(value)
		case "value":
			if nestedIndex == "" {
				draft.scalar = string(value)
				return
			}

			valueIndex, err := strconv.Atoi(nestedIndex)
			if err != nil {
				return
			}

			draft.values[valueIndex] = string(value)
		}
	})

	indexes := make([]int, 0, len(drafts))
	for index := range drafts {
		indexes = append(indexes, index)
	}
	sort.Ints(indexes)

	filters := make([]types.Filter, 0, len(indexes))
	for _, index := range indexes {
		draft := drafts[index]
		if draft == nil || draft.property == "" || draft.operator == "" {
			continue
		}

		var value interface{}
		if len(draft.values) > 0 {
			valueIndexes := make([]int, 0, len(draft.values))
			for valueIndex := range draft.values {
				valueIndexes = append(valueIndexes, valueIndex)
			}
			sort.Ints(valueIndexes)

			values := make([]interface{}, 0, len(valueIndexes))
			for _, valueIndex := range valueIndexes {
				values = append(values, normalizeJackpotFilterValue(draft.values[valueIndex]))
			}
			value = values
		} else if draft.scalar != "" {
			value = normalizeJackpotFilterValue(draft.scalar)
		}

		filters = append(filters, types.Filter{
			Property: draft.property,
			Operator: draft.operator,
			Value:    value,
		})
	}

	return filters
}

func normalizeJackpotFilterValue(raw string) interface{} {
	if intValue, err := strconv.Atoi(raw); err == nil {
		return intValue
	}

	if boolValue, err := strconv.ParseBool(raw); err == nil {
		return boolValue
	}

	return strings.TrimSpace(raw)
}

func buildJackpotGlobalResponse(row jackpotGlobalRow) JackpotGlobalResponse {
	return JackpotGlobalResponse{
		ID:                       row.ID,
		CurrentAmount:            row.CurrentAmount.Round(2).StringFixed(2),
		ThresholdAmount:          row.ThresholdAmount.Round(2).StringFixed(2),
		ChanceDenom:              row.ChanceDenom,
		PayoutPercent:            row.PayoutPercent.Round(2).StringFixed(2),
		MinEligibleBetAmount:     row.MinEligibleBetAmount.Round(2).StringFixed(2),
		CompanyTopupAmount:       row.CompanyTopupAmount.Round(2).StringFixed(2),
		JackpotFixedPayoutAmount: row.JackpotFixedPayoutAmount.Round(2).StringFixed(2),
		StatusID:                 row.StatusID,
		Order:                    row.Order,
		CreatedAt:                row.CreatedAt,
		CreatedBy:                row.CreatedBy,
		UpdatedAt:                row.UpdatedAt,
		UpdatedBy:                row.UpdatedBy,
		DeletedAt:                row.DeletedAt,
		DeletedBy:                row.DeletedBy,
	}
}

func buildJackpotCompanyTopupResponse(row jackpotCompanyTopupRow) JackpotCompanyTopupResponse {
	return JackpotCompanyTopupResponse{
		ID:                  row.ID,
		Username:            row.Username,
		Amount:              row.Amount.Round(2).StringFixed(2),
		CurrentAmountBefore: row.CurrentAmountBefore.Round(2).StringFixed(2),
		CurrentAmountAfter:  row.CurrentAmountAfter.Round(2).StringFixed(2),
		Note:                row.Note,
		Order:               row.Order,
		CreatedAt:           row.CreatedAt,
		CreatedBy:           row.CreatedBy,
		UpdatedAt:           row.UpdatedAt,
		UpdatedBy:           row.UpdatedBy,
		DeletedAt:           row.DeletedAt,
		DeletedBy:           row.DeletedBy,
	}
}

func buildJackpotMemberBonusResponse(row jackpotMemberBonusRow) JackpotMemberBonusResponse {
	return JackpotMemberBonusResponse{
		ID:                  row.ID,
		JackpotGlobalID:     row.JackpotGlobalID,
		MemberID:            row.MemberID,
		MemberUUID:          row.MemberUUID,
		UserName:            row.UserName,
		Amount:              row.Amount.Round(2).StringFixed(2),
		CurrentAmountBefore: row.CurrentAmountBefore.Round(2).StringFixed(2),
		CurrentAmountAfter:  row.CurrentAmountAfter.Round(2).StringFixed(2),
		Note:                row.Note,
		Order:               row.Order,
		CreatedAt:           row.CreatedAt,
		CreatedBy:           row.CreatedBy,
		UpdatedAt:           row.UpdatedAt,
		UpdatedBy:           row.UpdatedBy,
		DeletedAt:           row.DeletedAt,
		DeletedBy:           row.DeletedBy,
	}
}

func buildJackpotLedgerResponse(row jackpotLedgerRow) JackpotLedgerResponse {
	return JackpotLedgerResponse{
		ID:                     row.ID,
		JackpotGlobalID:        row.JackpotGlobalID,
		FishTypeID:             row.FishTypeID,
		FishTypeName:           row.FishTypeName,
		MemberID:               row.MemberID,
		BetID:                  row.BetID,
		TicketID:               row.TicketID,
		StatementID:            row.StatementID,
		SourceType:             row.SourceType,
		GlobalContributionCoin: row.GlobalContributionCoin.Round(2).StringFixed(2),
		PoolBefore:             row.PoolBefore.Round(2).StringFixed(2),
		PoolAfter:              row.PoolAfter.Round(2).StringFixed(2),
		CreatedAt:              row.CreatedAt,
		CreatedBy:              row.CreatedBy,
		UpdatedAt:              row.UpdatedAt,
		UpdatedBy:              row.UpdatedBy,
		DeletedAt:              row.DeletedAt,
		DeletedBy:              row.DeletedBy,
	}
}
