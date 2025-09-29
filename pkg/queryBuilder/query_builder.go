package queryBuilder

import (
	"fmt"
	"reflect"
	"strings"
)

type QueryBuilder struct {
	setParts     []string
	values       []interface{}
	paramNum     int
	withNilCheck bool
}

func NewQueryBuilder(withNilCheck bool) *QueryBuilder {
	return &QueryBuilder{
		paramNum:     1,
		withNilCheck: withNilCheck,
	}
}

func (qb *QueryBuilder) Set(field string, value any) *QueryBuilder {
	if qb.withNilCheck && isNil(value) {
		return qb
	}
	qb.paramNum++
	qb.setParts = append(qb.setParts, fmt.Sprintf("%s = $%d", field, qb.paramNum))
	qb.values = append(qb.values, value)
	return qb
}

func (qb *QueryBuilder) BuildUpdateQuery(table string, whereField string, whereValue interface{}) (string, []interface{}) {
	if len(qb.setParts) == 0 {
		return "", nil
	}

	setClause := strings.Join(qb.setParts, ", ") + ", updated_at = NOW()"

	// первый параметр всегда where
	allValues := append([]interface{}{whereValue}, qb.values...)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = $1", table, setClause, whereField)

	return query, allValues
}

func isNil(value interface{}) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}
	return false
}
