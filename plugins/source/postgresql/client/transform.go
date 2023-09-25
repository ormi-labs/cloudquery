package client

import (
	"database/sql/driver"
	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/decimal128"
	"github.com/apache/arrow/go/v14/arrow/decimal256"
	"github.com/jackc/pgx/v5/pgtype"
)

type transformer func(any) (any, error)

func transformerForDataType(dt arrow.DataType) transformer {
	switch dt := dt.(type) {
	case *arrow.StringType:
		return func(v any) (any, error) {
			if value, ok := v.(driver.Valuer); ok {
				if value == driver.Valuer(nil) {
					return nil, nil
				}

				val, err := value.Value()
				if err != nil {
					return nil, err
				}

				if s, ok := val.(string); ok {
					return s, nil
				}
			}
			return v, nil
		}
	case *arrow.Time32Type:
		return func(v any) (any, error) {
			t, err := v.(pgtype.Time).TimeValue()
			if err != nil {
				return nil, err
			}
			return stringForTime(t, dt.Unit), nil
		}
	case *arrow.Time64Type:
		return func(v any) (any, error) {
			t, err := v.(pgtype.Time).TimeValue()
			if err != nil {
				return nil, err
			}
			return stringForTime(t, dt.Unit), nil
		}
	case *arrow.Decimal128Type:
		return func(v any) (any, error) {
			if v == nil {
				return nil, nil
			}
			t := v.(pgtype.Numeric).Int
			return decimal128.FromBigInt(t), nil
		}
	case *arrow.Decimal256Type:
		return func(v any) (any, error) {
			if v == nil {
				return nil, nil
			}
			t := v.(pgtype.Numeric).Int
			return decimal256.FromBigInt(t), nil
		}
	default:
		return func(v any) (any, error) {
			return v, nil
		}
	}
}

func transformersForSchema(schema *arrow.Schema) []transformer {
	res := make([]transformer, schema.NumFields())
	for i := range res {
		res[i] = transformerForDataType(schema.Field(i).Type)
	}
	return res
}
