package cointransaction

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

type CoinTransactionShowRequest struct {
	PageOptions types.Paging   `json:"paging_options" query:"paging_options" validate:"required"`
	Sorts       []types.Sort   `json:"sorts,omitempty" query:"sorts"`
	Filters     []types.Filter `json:"filters,omitempty" query:"filters"`
}

var coinTransactionFilterPattern = regexp.MustCompile(`^filters\[(\d+)\]\[(property|operator|value)\](?:\[(\d+)\])?$`)

func (r *CoinTransactionShowRequest) Bind(c *fiber.Ctx, v *utils.Validator) error {
	if err := c.QueryParser(r); err != nil {
		custom_log.NewCustomLog("get_coin_transactions_failed", err.Error(), "error")
		return errors.New(utils.Translate("invalid_body", nil, c))
	}

	r.Filters = parseCoinTransactionFilters(c)

	if err := v.Validate(r, c); err != nil {
		return err
	}

	return nil
}

type CoinTransactionListResponse struct {
	Transactions []CoinTransactionResponse `json:"transactions"`
	Total        int                       `json:"-"`
}

type CoinTransactionResponse struct {
	MemberID               int       `json:"-"`
	MemberUUID             string    `json:"member_uuid"`
	UserName               string    `json:"username"`
	MemberCoinID           int       `json:"member_coin_id"`
	BeforeCoin             string    `json:"before_coin"`
	Amount                 string    `json:"amount"`
	TransactionTypeID      int       `json:"transaction_type_id"`
	TransactionGroupTypeID int       `json:"transaction_group_type_id"`
	TransactionDate        time.Time `json:"transaction_date"`
	RequireApproval        bool      `json:"require_approval"`
	Reference              string    `json:"reference"`
	Remark                 string    `json:"remark"`
	StatusID               int       `json:"status_id"`
	Order                  *int      `json:"order"`
	CreatedAt              time.Time `json:"created_at"`
}

type coinTransactionRow struct {
	MemberID               int             `db:"member_id"`
	MemberUUID             string          `db:"member_uuid"`
	UserName               string          `db:"username"`
	MemberCoinID           int             `db:"member_coin_id"`
	BeforeCoin             decimal.Decimal `db:"before_coin"`
	Amount                 decimal.Decimal `db:"amount"`
	TransactionTypeID      int             `db:"transaction_type_id"`
	TransactionGroupTypeID int             `db:"transaction_group_type_id"`
	TransactionDate        time.Time       `db:"transaction_date"`
	RequireApproval        bool            `db:"require_approval"`
	Reference              string          `db:"reference"`
	Remark                 string          `db:"remark"`
	StatusID               int             `db:"status_id"`
	Order                  *int            `db:"order"`
	CreatedAt              time.Time       `db:"created_at"`
}

func buildCoinTransactionResponse(row coinTransactionRow) CoinTransactionResponse {
	return CoinTransactionResponse{
		MemberID:               row.MemberID,
		MemberUUID:             row.MemberUUID,
		UserName:               row.UserName,
		MemberCoinID:           row.MemberCoinID,
		BeforeCoin:             row.BeforeCoin.Round(3).StringFixed(3),
		Amount:                 row.Amount.Round(3).StringFixed(3),
		TransactionTypeID:      row.TransactionTypeID,
		TransactionGroupTypeID: row.TransactionGroupTypeID,
		TransactionDate:        row.TransactionDate,
		RequireApproval:        row.RequireApproval,
		Reference:              row.Reference,
		Remark:                 row.Remark,
		StatusID:               row.StatusID,
		Order:                  row.Order,
		CreatedAt:              row.CreatedAt,
	}
}

type statementFilterDraft struct {
	property string
	operator string
	scalar   string
	values   map[int]string
}

func parseCoinTransactionFilters(c *fiber.Ctx) []types.Filter {
	drafts := make(map[int]*statementFilterDraft)

	c.Context().QueryArgs().VisitAll(func(key, value []byte) {
		matches := coinTransactionFilterPattern.FindStringSubmatch(string(key))
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
				values = append(values, normalizeCoinTransactionFilterValue(draft.values[valueIndex]))
			}
			value = values
		} else if draft.scalar != "" {
			value = normalizeCoinTransactionFilterValue(draft.scalar)
		}

		filters = append(filters, types.Filter{
			Property: draft.property,
			Operator: draft.operator,
			Value:    value,
		})
	}

	return filters
}

func normalizeCoinTransactionFilterValue(raw string) interface{} {
	if intValue, err := strconv.Atoi(raw); err == nil {
		return intValue
	}

	if boolValue, err := strconv.ParseBool(raw); err == nil {
		return boolValue
	}

	return strings.TrimSpace(raw)
}
