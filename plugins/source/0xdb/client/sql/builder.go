package sql

import (
	"bytes"
	"errors"
	"fmt"
)

type Select struct {
	columns    []string
	table      interface{}
	whereCond  string
	orderBy    string
	limitCount int64
}

// ToSQL creates sql statement
func (b *Select) ToSQL() (string, error) {
	if len(b.columns) == 0 {
		return "", errors.New("at least one column must be specified")
	}
	buf := new(bytes.Buffer)
	buf.WriteString("SELECT ")

	for i, col := range b.columns {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(col)
	}

	if b.table != nil {
		buf.WriteString(" FROM ")
		buf.WriteString(b.table.(string))
	}

	if b.whereCond != "" {
		buf.WriteString(" WHERE ")
		buf.WriteString(b.whereCond)
	}

	if b.orderBy != "" {
		buf.WriteString(" ORDER BY ")
		buf.WriteString(b.orderBy)
	}

	if b.limitCount > 0 {
		buf.WriteString(" LIMIT ")
		buf.WriteString(fmt.Sprint(b.limitCount))
	}
	return buf.String(), nil
}

// Select creates select part of SQL
func (b *Select) Select(columns ...string) *Select {
	b.columns = columns
	return b
}

// From creates from part of SQL
func (b *Select) From(table interface{}) *Select {
	b.table = table
	return b
}

func (b *Select) Where(arg1, op, arg2 string) *Select {
	b.whereCond = fmt.Sprint(arg1, " ", op, " ", arg2)
	return b
}

func (b *Select) Order(by string) *Select {
	b.orderBy = by
	return b
}

// Limit create limit part of SQL
func (b *Select) Limit(n int64) *Select {
	b.limitCount = n
	return b
}
