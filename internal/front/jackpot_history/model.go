package jackpothistory

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

type JackpotHistoryShowRequest struct {
	PageOptions types.Paging   `json:"paging_options" query:"paging_options" validate:"required"`
	Sorts       []types.Sort   `json:"sorts,omitempty" query:"sorts"`
	Filters     []types.Filter `json:"filters,omitempty" query:"filters"`
}

const (
	JackpotTypePool        = 1
	JackpotTypeMemberBonus = 2
)

func jackpotTypeName(jackpotType int) string {
	switch jackpotType {
	case JackpotTypePool:
		return "Pool Jackpot"
	case JackpotTypeMemberBonus:
		return "Member Bonus"
	default:
		return "Unknown"
	}
}

var jackpotHistoryFilterPattern = regexp.MustCompile(`^filters\[(\d+)\]\[(property|operator|value)\](?:\[(\d+)\])?$`)

func (r *JackpotHistoryShowRequest) Bind(c *fiber.Ctx, v *utils.Validator) error {
	if err := c.QueryParser(r); err != nil {
		custom_log.NewCustomLog("get_jackpot_history_failed", err.Error(), "error")
		return errors.New(utils.Translate("invalid_body", nil, c))
	}

	r.Filters = parseJackpotHistoryFilters(c)

	if err := v.Validate(r, c); err != nil {
		return err
	}

	return nil
}

type JackpotHistoryListResponse struct {
	Histories []JackpotHistoryResponse `json:"histories"`
	Total     int                      `json:"-"`
}

type JackpotHistoryResponse struct {
	ID              int64      `json:"id"`
	JackpotGlobalID int64      `json:"jackpot_global_id"`
	FishTypeID      int64      `json:"fish_type_id"`
	MemberID        *int64     `json:"member_id"`
	MemberName      string     `json:"member_name"`
	JackpotType      int        `json:"jackpot_type"`
	JackpotTypeName  string     `json:"jackpot_type_name"`
	FishTypeName     string     `json:"fish_type_name"`
	PayoutCoin      string     `json:"payout_coin"`
	PoolBefore      string     `json:"pool_before"`
	PoolAfter       string     `json:"pool_after"`
	CreatedAt       time.Time  `json:"created_at"`
	CreatedBy       *int       `json:"created_by"`
	UpdatedAt       *time.Time `json:"-"`
	UpdatedBy       *int       `json:"-"`
	DeletedAt       *time.Time `json:"-"`
	DeletedBy       *int       `json:"-"`
}

type jackpotHistoryRow struct {
	ID              int64           `db:"id"`
	JackpotGlobalID int64           `db:"jackpot_global_id"`
	FishTypeID      int64           `db:"fish_type_id"`
	MemberID        *int64          `db:"member_id"`
	MemberName      string          `db:"member_name"`
	JackpotType      int             `db:"jackpot_type"`
	FishTypeName     string          `db:"fish_type_name"`
	PayoutCoin      decimal.Decimal `db:"payout_coin"`
	PoolBefore      decimal.Decimal `db:"pool_before"`
	PoolAfter       decimal.Decimal `db:"pool_after"`
	CreatedAt       time.Time       `db:"created_at"`
	CreatedBy       *int            `db:"created_by"`
	UpdatedAt       *time.Time      `db:"updated_at"`
	UpdatedBy       *int            `db:"updated_by"`
	DeletedAt       *time.Time      `db:"deleted_at"`
	DeletedBy       *int            `db:"deleted_by"`
}

func buildJackpotHistoryResponse(row jackpotHistoryRow) JackpotHistoryResponse {
	return JackpotHistoryResponse{
		ID:              row.ID,
		JackpotGlobalID: row.JackpotGlobalID,
		FishTypeID:      row.FishTypeID,
		MemberID:        row.MemberID,
		MemberName:      row.MemberName,
		JackpotType:     row.JackpotType,
		JackpotTypeName: jackpotTypeName(row.JackpotType),
		FishTypeName:    row.FishTypeName,
		PayoutCoin:      row.PayoutCoin.Round(2).StringFixed(2),
		PoolBefore:      row.PoolBefore.Round(2).StringFixed(2),
		PoolAfter:       row.PoolAfter.Round(2).StringFixed(2),
		CreatedAt:       row.CreatedAt,
		CreatedBy:       row.CreatedBy,
		UpdatedAt:       row.UpdatedAt,
		UpdatedBy:       row.UpdatedBy,
		DeletedAt:       row.DeletedAt,
		DeletedBy:       row.DeletedBy,
	}
}

type jackpotHistoryFilterDraft struct {
	property string
	operator string
	scalar   string
	values   map[int]string
}

func parseJackpotHistoryFilters(c *fiber.Ctx) []types.Filter {
	drafts := make(map[int]*jackpotHistoryFilterDraft)

	c.Context().QueryArgs().VisitAll(func(key, value []byte) {
		matches := jackpotHistoryFilterPattern.FindStringSubmatch(string(key))
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
			draft = &jackpotHistoryFilterDraft{
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
				values = append(values, normalizeJackpotHistoryFilterValue(draft.values[valueIndex]))
			}
			value = values
		} else if draft.scalar != "" {
			value = normalizeJackpotHistoryFilterValue(draft.scalar)
		}

		filters = append(filters, types.Filter{
			Property: draft.property,
			Operator: draft.operator,
			Value:    value,
		})
	}

	return filters
}

func normalizeJackpotHistoryFilterValue(raw string) interface{} {
	if intValue, err := strconv.Atoi(raw); err == nil {
		return intValue
	}

	if boolValue, err := strconv.ParseBool(raw); err == nil {
		return boolValue
	}

	return strings.TrimSpace(raw)
}
