package member

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
)

type MemberInfo struct {
	ID           int                    `json:"id" db:"id"`
	UserUUID     string                 `json:"user_uuid" db:"user_uuid"`
	UserName     string                 `json:"user_name" db:"user_name"`
	LoginID      string                 `json:"login_id" db:"login_id"`
	PhoneNumber  string                 `json:"phone_number" db:"phone_number"`
	ProfilePhoto *string                `json:"profile_photo" db:"profile_photo"`
	LanguageID   *int                   `json:"language_id" db:"language_id"`
	CurrencyID   int                    `json:"currency_id" db:"currency_id"`
	CoinAmount   string                 `json:"coin_amount" db:"coin_amount"`
	Balances     []MemberCurrencyBalance `json:"balances"`
	Remark       *string                `json:"remark" db:"remark"`
	Nickname     *string                `json:"nickname" db:"nickname"`
	LoginSession *string                `json:"login_session" db:"login_session"`
	LastLoginAt  *time.Time             `json:"last_login_at" db:"last_login_at"`
	IsOnline     bool                   `json:"is_online" db:"is_online"`
	StatusID     int                    `json:"status_id" db:"status_id"`
	Timezone     *string                `json:"timezone" db:"timezone"`
	Pattern      *string                `json:"pattern" db:"pattern"`
	Order        *int                   `json:"order" db:"order"`
	CreatedBy    int                    `json:"created_by" db:"created_by"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
	UpdatedBy    *int                   `json:"updated_by" db:"updated_by"`
	UpdatedAt    *time.Time             `json:"updated_at" db:"updated_at"`
}

type MemberCurrencyBalance struct {
	CurrencyID     int        `json:"currency_id" db:"currency_id"`
	CurrencyName   string     `json:"currency_name" db:"currency_name"`
	CurrencyCode   string     `json:"currency_code" db:"currency_code"`
	CurrencySymbol string     `json:"currency_symbol" db:"currency_symbol"`
	Balance        string     `json:"balance" db:"balance"`
}

type GetAllMembersResponse struct {
	Members []MemberInfo `json:"members"`
	Total   int          `json:"-"`
}

type MemberShowRequest struct {
	PageOptions types.Paging   `json:"paging_options" query:"paging_options" validate:"required"`
	Sorts       []types.Sort   `json:"sorts,omitempty" query:"sorts"`
	Filters     []types.Filter `json:"filters,omitempty" query:"filters"`
}

var memberFilterPattern = regexp.MustCompile(`^filters\[(\d+)\]\[(property|operator|value)\](?:\[(\d+)\])?$`)

func (r *MemberShowRequest) Bind(c *fiber.Ctx, v *utils.Validator) error {
	if err := c.QueryParser(r); err != nil {
		custom_log.NewCustomLog("get_all_members_failed", err.Error(), "error")
		return errors.New(utils.Translate("invalid_body", nil, c))
	}

	r.Filters = parseMemberFilters(c)

	if err := v.Validate(r, c); err != nil {
		return err
	}

	return nil
}

type memberFilterDraft struct {
	property string
	operator string
	scalar   string
	values   map[int]string
}

func parseMemberFilters(c *fiber.Ctx) []types.Filter {
	drafts := make(map[int]*memberFilterDraft)

	c.Context().QueryArgs().VisitAll(func(key, value []byte) {
		matches := memberFilterPattern.FindStringSubmatch(string(key))
		if len(matches) != 4 {
			return
		}

		idx, err := strconv.Atoi(matches[1])
		if err != nil {
			return
		}

		field := matches[2]
		pos := matches[3]
		rawValue := string(value)

		draft, ok := drafts[idx]
		if !ok {
			draft = &memberFilterDraft{
				values: make(map[int]string),
			}
			drafts[idx] = draft
		}

		switch field {
		case "property":
			draft.property = rawValue
		case "operator":
			draft.operator = strings.ToLower(rawValue)
		case "value":
			if pos == "" {
				draft.scalar = rawValue
				return
			}

			valueIdx, err := strconv.Atoi(pos)
			if err != nil {
				return
			}
			draft.values[valueIdx] = rawValue
		}
	})

	keys := make([]int, 0, len(drafts))
	for idx := range drafts {
		keys = append(keys, idx)
	}
	sort.Ints(keys)

	filters := make([]types.Filter, 0, len(keys))
	for _, idx := range keys {
		draft := drafts[idx]
		if draft.property == "" || draft.operator == "" {
			continue
		}

		var value interface{}
		switch draft.operator {
		case "in", "between":
			if len(draft.values) == 0 {
				continue
			}

			nestedKeys := make([]int, 0, len(draft.values))
			for nestedIdx := range draft.values {
				nestedKeys = append(nestedKeys, nestedIdx)
			}
			sort.Ints(nestedKeys)

			values := make([]interface{}, 0, len(nestedKeys))
			for _, nestedIdx := range nestedKeys {
				values = append(values, draft.values[nestedIdx])
			}
			value = values
		default:
			value = draft.scalar
		}

		filters = append(filters, types.Filter{
			Property: draft.property,
			Operator: draft.operator,
			Value:    value,
		})
	}

	return filters
}
