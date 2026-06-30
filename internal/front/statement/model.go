package statement

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

type StatementShowRequest struct {
	PageOptions types.Paging   `json:"paging_options" query:"paging_options" validate:"required"`
	Sorts       []types.Sort   `json:"sorts,omitempty" query:"sorts"`
	Filters     []types.Filter `json:"filters,omitempty" query:"filters"`
}

var statementFilterPattern = regexp.MustCompile(`^filters\[(\d+)\]\[(property|operator|value)\](?:\[(\d+)\])?$`)

func (r *StatementShowRequest) Bind(c *fiber.Ctx, v *utils.Validator) error {
	if err := c.QueryParser(r); err != nil {
		custom_log.NewCustomLog("get_statements_failed", err.Error(), "error")
		return errors.New(utils.Translate("invalid_body", nil, c))
	}
	r.Filters = parseStatementFilters(c)

	if err := v.Validate(r, c); err != nil {
		return err
	}

	return nil
}

type StatementListResponse struct {
	Statements  []StatementResponse `json:"statements"`
	TotalReport TotalReportResponse `json:"total_report"`
	Total       int                 `json:"-"`
}

type TotalReport struct {
	TotalBet      decimal.Decimal `db:"total_bet"`
	TotalValidBet decimal.Decimal `db:"total_valid_bet"`
	TotalWinLose  decimal.Decimal `db:"total_winlose"`
}

type TotalReportResponse struct {
	TotalBet      string `json:"total_bet"`
	TotalValidBet string `json:"total_valid_bet"`
	TotalWinLose  string `json:"total_winlose"`
}

type StatementResponse struct {
	StatementUUID     zz    string     `json:"statement_uuid"`
	MemberID              int64      `json:"-"`
	MemberUUID            string     `json:"member_uuid"`
	UserName              string     `json:"username"`
	SessionID             int64      `json:"session_id"`
	SessionNo             string     `json:"session_no"`
	TicketID              int64      `json:"ticket_id"`
	TicketNo              string     `json:"ticket_no"`
	BetNo                 string     `json:"bet_no"`
	FishTypeID            int        `json:"fish_type_id"`
	FishTypeName          string     `json:"fish_type_name"`
	BetAmount             string     `json:"bet_amount"`
	BetValid              string     `json:"bet_valid"`
	BetInvalid            string     `json:"bet_invalid"`
	TotalBetAmount        string     `json:"total_bet_amount"`
	TotalBetInvalidAmount string     `json:"total_bet_invalid_amount"`
	IsKill                bool       `json:"is_kill"`
	WinLose               string     `json:"win_lose"`
	KillReward            string     `json:"kill_reward"`
	MissReward            string     `json:"miss_reward"`
	PayoutAmount          string     `json:"payout_amount"`
	JackpotWinAmount      string     `json:"jackpot_win_amount"`
	Jackpot               string     `json:"jackpot"`
	TotalWinLose          string     `json:"total_win_lose"`
	SyncID                *int       `json:"sync_id"`
	IsSynced              bool       `json:"is_synced"`
	StatementAt           time.Time  `json:"statement_at"`
	StatusID              int16      `json:"status_id"`
	Order                 *int       `json:"order"`
	CreatedAt             time.Time  `json:"created_at"`
	CreatedBy             *int       `json:"created_by"`
	UpdatedAt             *time.Time `json:"updated_at"`
	UpdatedBy             *int       `json:"updated_by"`
	DeletedAt             *time.Time `json:"deleted_at"`
	DeletedBy             *int       `json:"deleted_by"`
}

type statementRow struct {
	StatementUUID         string          `db:"statement_uuid"`
	MemberID              int64           `db:"member_id"`
	MemberUUID            string          `db:"member_uuid"`
	UserName              string          `db:"username"`
	SessionID             int64           `db:"session_id"`
	SessionNo             string          `db:"session_no"`
	TicketID              int64           `db:"ticket_id"`
	TicketNo              string          `db:"ticket_no"`
	BetNo                 string          `db:"bet_no"`
	FishTypeID            int             `db:"fish_type_id"`
	FishTypeName          string          `db:"fish_type_name"`
	BetAmount             decimal.Decimal `db:"bet_amount"`
	BetValid              decimal.Decimal `db:"bet_valid"`
	BetInvalid            decimal.Decimal `db:"bet_invalid"`
	TotalBetAmount        decimal.Decimal `db:"total_bet_amount"`
	TotalBetInvalidAmount decimal.Decimal `db:"total_bet_invalid_amount"`
	IsKill                bool            `db:"is_kill"`
	WinLose               string          `db:"win_lose"`
	KillReward            decimal.Decimal `db:"kill_reward"`
	MissReward            decimal.Decimal `db:"miss_reward"`
	PayoutAmount          decimal.Decimal `db:"payout_amount"`
	JackpotWinAmount      decimal.Decimal `db:"jackpot_win_amount"`
	Jackpot               decimal.Decimal `db:"jackpot"`
	TotalWinLose          decimal.Decimal `db:"total_win_lose"`
	SyncID                *int            `db:"sync_id"`
	IsSynced              bool            `db:"is_synced"`
	StatementAt           time.Time       `db:"statement_at"`
	StatusID              int16           `db:"status_id"`
	Order                 *int            `db:"order"`
	CreatedAt             time.Time       `db:"created_at"`
	CreatedBy             *int            `db:"created_by"`
	UpdatedAt             *time.Time      `db:"updated_at"`
	UpdatedBy             *int            `db:"updated_by"`
	DeletedAt             *time.Time      `db:"deleted_at"`
	DeletedBy             *int            `db:"deleted_by"`
}

func buildStatementResponse(row statementRow) StatementResponse {
	return StatementResponse{
		StatementUUID:         row.StatementUUID,
		MemberID:              row.MemberID,
		MemberUUID:            row.MemberUUID,
		UserName:              row.UserName,
		SessionID:             row.SessionID,
		SessionNo:             row.SessionNo,
		TicketID:              row.TicketID,
		TicketNo:              row.TicketNo,
		BetNo:                 row.BetNo,
		FishTypeID:            row.FishTypeID,
		FishTypeName:          row.FishTypeName,
		BetAmount:             row.BetAmount.Round(3).StringFixed(3),
		BetValid:              row.BetValid.Round(3).StringFixed(3),
		BetInvalid:            row.BetInvalid.Round(3).StringFixed(3),
		TotalBetAmount:        row.TotalBetAmount.Round(3).StringFixed(3),
		TotalBetInvalidAmount: row.TotalBetInvalidAmount.Round(3).StringFixed(3),
		IsKill:                row.IsKill,
		WinLose:               row.WinLose,
		KillReward:            row.KillReward.Round(3).StringFixed(3),
		MissReward:            row.MissReward.Round(3).StringFixed(3),
		PayoutAmount:          row.PayoutAmount.Round(3).StringFixed(3),
		JackpotWinAmount:      row.JackpotWinAmount.Round(3).StringFixed(3),
		Jackpot:               row.Jackpot.Round(3).StringFixed(3),
		TotalWinLose:          row.TotalWinLose.Round(3).StringFixed(3),
		SyncID:                row.SyncID,
		IsSynced:              row.IsSynced,
		StatementAt:           row.StatementAt,
		StatusID:              row.StatusID,
		Order:                 row.Order,
		CreatedAt:             row.CreatedAt,
		CreatedBy:             row.CreatedBy,
		UpdatedAt:             row.UpdatedAt,
		UpdatedBy:             row.UpdatedBy,
		DeletedAt:             row.DeletedAt,
		DeletedBy:             row.DeletedBy,
	}
}

func buildTotalReportResponse(totalBet decimal.Decimal, totalValidBet decimal.Decimal, totalWinLose decimal.Decimal) TotalReportResponse {
	return TotalReportResponse{
		TotalBet:      totalBet.Round(3).StringFixed(3),
		TotalValidBet: totalValidBet.Round(3).StringFixed(3),
		TotalWinLose:  totalWinLose.Round(3).StringFixed(3),
	}
}

type statementFilterDraft struct {
	property string
	operator string
	scalar   string
	values   map[int]string
}

func parseStatementFilters(c *fiber.Ctx) []types.Filter {
	drafts := make(map[int]*statementFilterDraft)

	c.Context().QueryArgs().VisitAll(func(key, value []byte) {
		matches := statementFilterPattern.FindStringSubmatch(string(key))
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
			draft = &statementFilterDraft{
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

	filterIndexes := make([]int, 0, len(drafts))
	for index := range drafts {
		filterIndexes = append(filterIndexes, index)
	}
	sort.Ints(filterIndexes)

	filters := make([]types.Filter, 0, len(filterIndexes))
	for _, index := range filterIndexes {
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
				values = append(values, normalizeStatementFilterValue(draft.values[valueIndex]))
			}
			value = values
		} else if draft.scalar != "" {
			value = normalizeStatementFilterValue(draft.scalar)
		}

		filters = append(filters, types.Filter{
			Property: draft.property,
			Operator: draft.operator,
			Value:    value,
		})
	}

	return filters
}

func normalizeStatementFilterValue(raw string) interface{} {
	if intValue, err := strconv.Atoi(raw); err == nil {
		return intValue
	}

	if boolValue, err := strconv.ParseBool(raw); err == nil {
		return boolValue
	}

	return strings.TrimSpace(raw)
}
