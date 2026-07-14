package memberbonus

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

var jackpotMemberBonusFilterPattern = regexp.MustCompile(`^filters\[(\d+)\]\[(property|operator|value)\](?:\[(\d+)\])?$`)

type JackpotMemberBonusShowRequest struct {
	PageOptions types.Paging   `json:"paging_options" query:"paging_options" validate:"required"`
	Sorts       []types.Sort   `json:"sorts,omitempty" query:"sorts"`
	Filters     []types.Filter `json:"filters,omitempty" query:"filters"`
}

type CreateJackpotMemberBonusRequest struct {
	MemberName string `json:"member_name" validate:"required"`
	Amount     string `json:"amount" validate:"required"`
	Note       string `json:"note"`
}

func (r *JackpotMemberBonusShowRequest) Bind(c *fiber.Ctx, v *utils.Validator) error {
	if err := c.QueryParser(r); err != nil {
		custom_log.NewCustomLog("get_jackpot_member_bonuses_failed", err.Error(), "error")
		return errors.New(utils.Translate("invalid_body", nil, c))
	}

	r.Filters = parseJackpotMemberBonusFilters(c)

	if err := v.Validate(r, c); err != nil {
		return err
	}

	return nil
}

func (r *CreateJackpotMemberBonusRequest) Bind(c *fiber.Ctx, v *utils.Validator) error {
	if err := c.BodyParser(r); err != nil {
		custom_log.NewCustomLog("create_jackpot_member_bonus_failed", err.Error(), "error")
		return errors.New(utils.Translate("invalid_body", nil, c))
	}

	r.MemberName = strings.TrimSpace(r.MemberName)

	if err := v.Validate(r, c); err != nil {
		return err
	}

	return nil
}

type JackpotMemberBonusListResponse struct {
	Bonuses []JackpotMemberBonusResponse `json:"bonuses"`
	Total   int                          `json:"-"`
}

type JackpotMemberBonusResponse struct {
	ID            int64      `json:"id"`
	MemberID      int64      `json:"member_id"`
	MemberName    string     `json:"member_name"`
	MemberUUID    string     `json:"-"`
	CreatedByName string     `json:"created_by_name"`
	Amount        string     `json:"amount"`
	Note          string     `json:"note"`
	Order         int        `json:"order"`
	StatusID      int        `json:"status_id"`
	CreatedAt     time.Time  `json:"created_at"`
	CreatedBy     *int       `json:"created_by"`
	UpdatedAt     *time.Time `json:"-"`
	UpdatedBy     *int       `json:"-"`
	DeletedAt     *time.Time `json:"-"`
	DeletedBy     *int       `json:"-"`
}

type jackpotMemberBonusRow struct {
	ID            int64           `db:"id"`
	MemberID      int64           `db:"member_id"`
	MemberName    string          `db:"member_name"`
	MemberUUID    string          `db:"member_uuid"`
	CreatedByName string          `db:"created_by_name"`
	Amount        decimal.Decimal `db:"amount"`
	Note          string          `db:"note"`
	Order         int             `db:"order"`
	StatusID      int             `db:"status_id"`
	CreatedAt     time.Time       `db:"created_at"`
	CreatedBy     *int            `db:"created_by"`
	UpdatedAt     *time.Time      `db:"updated_at"`
	UpdatedBy     *int            `db:"updated_by"`
	DeletedAt     *time.Time      `db:"deleted_at"`
	DeletedBy     *int            `db:"deleted_by"`
}

type jackpotMemberBonusFilterDraft struct {
	property string
	operator string
	scalar   string
	values   map[int]string
}

func parseJackpotMemberBonusFilters(c *fiber.Ctx) []types.Filter {
	drafts := make(map[int]*jackpotMemberBonusFilterDraft)

	c.Context().QueryArgs().VisitAll(func(key, value []byte) {
		matches := jackpotMemberBonusFilterPattern.FindStringSubmatch(string(key))
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
			draft = &jackpotMemberBonusFilterDraft{
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
				values = append(values, normalizeJackpotMemberBonusFilterValue(draft.values[valueIndex]))
			}
			value = values
		} else if draft.scalar != "" {
			value = normalizeJackpotMemberBonusFilterValue(draft.scalar)
		}

		filters = append(filters, types.Filter{
			Property: draft.property,
			Operator: draft.operator,
			Value:    value,
		})
	}

	return filters
}

func normalizeJackpotMemberBonusFilterValue(raw string) interface{} {
	if intValue, err := strconv.Atoi(raw); err == nil {
		return intValue
	}

	if boolValue, err := strconv.ParseBool(raw); err == nil {
		return boolValue
	}

	return strings.TrimSpace(raw)
}

func buildJackpotMemberBonusResponse(row jackpotMemberBonusRow) JackpotMemberBonusResponse {
	return JackpotMemberBonusResponse{
		ID:            row.ID,
		MemberID:      row.MemberID,
		MemberName:    row.MemberName,
		MemberUUID:    row.MemberUUID,
		CreatedByName: row.CreatedByName,
		Amount:        row.Amount.Round(2).StringFixed(2),
		Note:          row.Note,
		Order:         row.Order,
		StatusID:      row.StatusID,
		CreatedAt:     row.CreatedAt,
		CreatedBy:     row.CreatedBy,
		UpdatedAt:     row.UpdatedAt,
		UpdatedBy:     row.UpdatedBy,
		DeletedAt:     row.DeletedAt,
		DeletedBy:     row.DeletedBy,
	}
}
