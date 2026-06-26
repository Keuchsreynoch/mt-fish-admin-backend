package postgres

import (
	types "fish_shooting_admin_backend/pkg/model"
	"fmt"

	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func BuildSQLFilter(filters []types.Filter) (string, []interface{}) {

	var clauses []string
	var params []interface{}
	placeholder := 1

	// load timezone
	appTimezone := os.Getenv("APP_TIMEZONE")
	if appTimezone == "" {
		appTimezone = "Asia/Phnom_Penh"
	}

	// load location
	location, err := time.LoadLocation(appTimezone)
	if err != nil {
		return "", nil
	}

	// convert value to time.Time
	convertToTime := func(v interface{}) (time.Time, bool) {
		switch val := v.(type) {

		case string:
			parsed, err := time.ParseInLocation(
				"2006-01-02",
				val,
				location,
			)

			if err == nil {
				return parsed, true
			}

		case time.Time:
			return val.In(location), true
		}

		return time.Time{}, false
	}

	// convert boolean
	convertToBool := func(v interface{}) (bool, bool) {
		switch val := v.(type) {

		case bool:
			return val, true

		case string:
			switch strings.ToLower(val) {
			case "true":
				return true, true
			case "false":
				return false, true
			}
		}

		return false, false
	}

	toInterfaceSlice := func(v interface{}) ([]interface{}, bool) {
		switch vals := v.(type) {
		case []interface{}:
			return vals, true
		case []string:
			result := make([]interface{}, 0, len(vals))
			for _, item := range vals {
				result = append(result, item)
			}
			return result, true
		case []int:
			result := make([]interface{}, 0, len(vals))
			for _, item := range vals {
				result = append(result, item)
			}
			return result, true
		}

		return nil, false
	}

	for _, f := range filters {

		field := f.Property
		op := strings.ToLower(f.Operator)

		switch op {

		case "eq", "neq", "lt", "lte", "gt", "gte":

			sqlOp := map[string]string{
				"eq":  "=",
				"neq": "!=",
				"lt":  "<",
				"lte": "<=",
				"gt":  ">",
				"gte": ">=",
			}[op]

			if t, ok := convertToTime(f.Value); ok {
				f.Value = t
			} else if b, ok := convertToBool(f.Value); ok {
				f.Value = b
			}

			// compare date only
			if _, ok := f.Value.(time.Time); ok {
				clauses = append(
					clauses,
					fmt.Sprintf(
						"DATE(%s) %s DATE($%d)",
						field,
						sqlOp,
						placeholder,
					),
				)
			} else {
				clauses = append(
					clauses,
					fmt.Sprintf(
						"%s %s $%d",
						field,
						sqlOp,
						placeholder,
					),
				)
			}

			params = append(params, f.Value)
			placeholder++

		case "like":

			clauses = append(
				clauses,
				fmt.Sprintf(
					"%s LIKE $%d",
					field,
					placeholder,
				),
			)

			params = append(params, f.Value)
			placeholder++

		case "in":

			values, ok := toInterfaceSlice(f.Value)
			if !ok || len(values) == 0 {
				continue
			}

			var ph []string

			for _, v := range values {
				ph = append(ph, fmt.Sprintf("$%d", placeholder))
				params = append(params, v)
				placeholder++
			}

			clauses = append(
				clauses,
				fmt.Sprintf(
					"%s IN (%s)",
					field,
					strings.Join(ph, ", "),
				),
			)

		case "between":

			vals, ok := toInterfaceSlice(f.Value)
			if !ok || len(vals) != 2 {
				continue
			}

			start, ok1 := convertToTime(vals[0])
			end, ok2 := convertToTime(vals[1])

			if ok1 && ok2 {

				// inclusive end of day
				end = end.Add(24*time.Hour - time.Second)

				clauses = append(
					clauses,
					fmt.Sprintf(
						"%s BETWEEN $%d AND $%d",
						field,
						placeholder,
						placeholder+1,
					),
				)

				params = append(params, start, end)
				placeholder += 2

			} else {

				clauses = append(
					clauses,
					fmt.Sprintf(
						"%s BETWEEN $%d AND $%d",
						field,
						placeholder,
						placeholder+1,
					),
				)

				params = append(params, vals[0], vals[1])
				placeholder += 2
			}
		}
	}

	if len(clauses) == 0 {
		return "", nil
	}

	return strings.Join(clauses, " AND "), params
}

func BuildSort(sorts []types.Sort) string {

	var orderClauses []string

	for _, s := range sorts {
		field := s.Property
		direction := strings.ToUpper(s.Direction)

		// ensure the direction is either ASC or DESC
		if direction != "ASC" && direction != "DESC" {
			direction = "DESC"
		}

		// add the sort clause to the list of order clauses
		orderClauses = append(orderClauses, fmt.Sprintf("%s %s", field, direction))
	}

	if len(orderClauses) == 0 {
		return ""
	}

	// join the clauses with commas and return the final order by string
	return "ORDER BY " + strings.Join(orderClauses, ", ")
}

func BuildPaging(page int, perPage int) string {
	// var params []interface{}

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage
	limit := perPage

	// params = append(params, offset, limit)

	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}

func GetIdByUuid(space_name string, uuid_field_name string, uuid_str string, db *sqlx.Tx) (*int, error) {
	var id int

	// Parse the UUID
	uid, err := uuid.Parse(uuid_str)
	if err != nil {
		return nil, err
	}

	// Define the SQL query
	sql := fmt.Sprintf(`SELECT id FROM %s WHERE %s=$1`, space_name, uuid_field_name)

	// Execute the query
	err = db.Get(&id, sql, uid)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// GetSeqNextVal returns the next value from a sequence
// SeqResult struct to store sequence result
type SeqResult struct {
	ID int `db:"id"`
}

// Supports both normal DB connection and transactions
func GetSeqNextVal(seqName string, exec sqlx.Ext) (*int, error) {
	var result SeqResult
	sql := `SELECT nextval($1) AS id`

	// Execute query using either DB or transaction
	err := sqlx.Get(exec, &result, sql, seqName)
	if err != nil {
		return nil, fmt.Errorf("failed to get sequence value: %w", err)
	}
	return &result.ID, nil
}

// SetSeqNextVal sets and returns the next sequence value
func SetSeqNextVal(seq_name string, db *sqlx.Tx) (*int, error) {
	var id int

	// Define the SQL query - adjust to PostgreSQL sequence operations
	sql := fmt.Sprintf(`SELECT nextval('%s') as id`, seq_name)

	// Execute the query
	err := db.Get(&id, sql)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// IsExists checks if a record exists with the given field value
func IsExists(spaceName, fieldName string, value interface{}, conn interface{}) (bool, error) {
	var exists bool

	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE %s = $1 AND deleted_at IS NULL)`, spaceName, fieldName)

	switch db := conn.(type) {
	case *sqlx.DB:
		err := db.Get(&exists, query, value)
		return exists, err
	case *sqlx.Tx:
		err := db.Get(&exists, query, value)
		return exists, err
	default:
		return false, fmt.Errorf("unsupported DB type: %T", conn)
	}
}

// IsExistsWhere checks if a record exists with custom WHERE conditions
func IsExistsWhere(tbl_name string, where_sqlstr string, args []interface{}, conn interface{}) (bool, error) {
	var exists bool

	// Define the SQL query
	sql := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE %s AND deleted_at IS NULL)`, tbl_name, where_sqlstr)

	// Execute the query
	switch db := conn.(type) {
	case *sqlx.DB:
		err := db.Get(&exists, sql, args)
		return exists, err
	case *sqlx.Tx:
		err := db.Get(&exists, sql, args)
		return exists, err
	default:
		return false, fmt.Errorf("unsupported DB type: %T", conn)
	}
}
