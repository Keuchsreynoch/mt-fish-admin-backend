package report

import (
	"errors"
	custom_log "fish_shooting_admin_backend/pkg/logs"
	types "fish_shooting_admin_backend/pkg/model"
	"fish_shooting_admin_backend/pkg/utils"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

type ReportShowRequest struct {
	PageOptions types.Paging   `json:"paging_options" query:"paging_options" validate:"required"`
	Sorts       []types.Sort   `json:"sorts,omitempty" query:"sorts"`
	Filters     []types.Filter `json:"filters,omitempty" query:"filters"`
}

var reportFilterPattern = regexp.MustCompile(`^filters\[(\d+)\]\[(property|operator|value)\](?:\[(\d+)\])?$`)

func (r *ReportShowRequest) Bind(c *fiber.Ctx, v *utils.Validator) error {
	if err := c.QueryParser(r); err != nil {
		custom_log.NewCustomLog("get_statements_failed", err.Error(), "error")
		return errors.New(utils.Translate("invalid_body", nil, c))
	}

	r.Filters = parseReportFilters(c)

	if err := v.Validate(r, c); err != nil {
		return err
	}

	return nil
}

type ReportMemberListResponse struct {
	Reports     []ReportMemberResponse `json:"reports"`
	TotalReport TotalReportResponse   `json:"total_report"`
	Total       int                   `json:"-"`
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

type ReportMemberResponse struct {
	MemberID         int64  `json:"member_id"`
	MemberUUID       string `json:"member_uuid"`
	MemberName       string `json:"member_name"`
	TotalBetAmount   string `json:"total_bet_amount"`
	TotalValidBet    string `json:"total_valid_bet"`
	TotalWinLose     string `json:"total_win_lose"`
	JackpotWinAmount string `json:"jackpot_win_amount"`
}

type reportMemberRow struct {
	MemberID         int64           `db:"member_id"`
	MemberUUID       string          `db:"member_uuid"`
	MemberName       string          `db:"member_name"`
	TotalBetAmount   decimal.Decimal `db:"total_bet_amount"`
	TotalValidBet    decimal.Decimal `db:"total_valid_bet"`
	TotalWinLose     decimal.Decimal `db:"total_win_lose"`
	JackpotWinAmount decimal.Decimal `db:"jackpot_win_amount"`
}

func buildReportMemberResponse(row reportMemberRow) ReportMemberResponse {
	return ReportMemberResponse{
		MemberID:         row.MemberID,
		MemberUUID:       row.MemberUUID,
		MemberName:       row.MemberName,
		TotalBetAmount:   row.TotalBetAmount.Round(3).StringFixed(3),
		TotalValidBet:    row.TotalValidBet.Round(3).StringFixed(3),
		TotalWinLose:     row.TotalWinLose.Round(3).StringFixed(3),
		JackpotWinAmount: row.JackpotWinAmount.Round(3).StringFixed(3),
	}
}

func buildReportTotalResponse(totalBet decimal.Decimal, totalValidBet decimal.Decimal, totalWinLose decimal.Decimal) TotalReportResponse {
	return TotalReportResponse{
		TotalBet:      totalBet.Round(3).StringFixed(3),
		TotalValidBet: totalValidBet.Round(3).StringFixed(3),
		TotalWinLose:  totalWinLose.Round(3).StringFixed(3),
	}
}

type reportFilterDraft struct {
	property string
	operator string
	scalar   string
	values   map[int]string
}

func parseReportFilters(c *fiber.Ctx) []types.Filter {
	drafts := make(map[int]*reportFilterDraft)

	c.Context().QueryArgs().VisitAll(func(key, value []byte) {
		matches := reportFilterPattern.FindStringSubmatch(string(key))
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
			draft = &reportFilterDraft{
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
				values = append(values, normalizeReportFilterValue(draft.values[valueIndex]))
			}
			value = values
		} else if draft.scalar != "" {
			value = normalizeReportFilterValue(draft.scalar)
		}

		filters = append(filters, types.Filter{
			Property: draft.property,
			Operator: draft.operator,
			Value:    value,
		})
	}

	return filters
}

func normalizeReportFilterValue(raw string) interface{} {
	if intValue, err := strconv.Atoi(raw); err == nil {
		return intValue
	}

	if boolValue, err := strconv.ParseBool(raw); err == nil {
		return boolValue
	}

	return strings.TrimSpace(raw)
}
