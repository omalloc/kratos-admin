package gormschema

import (
	"context"
	"reflect"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"gorm.io/gorm/schema"
)

var _ schema.SerializerInterface = (*StringSlice[int])(nil)

type StringSlice[T ~int | int32 | int64 | uint | uint32 | uint64] []T

// Scan implements schema.SerializerInterface.
func (g *StringSlice[T]) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) error {
	switch v := dbValue.(type) {
	case string:
		*g = lo.Map(strings.Split(v, ","), func(item string, _ int) T {
			vv, _ := strconv.ParseInt(item, 10, 64)
			return T(vv)
		})
	default:
		*g = []T{}
	}
	return nil
}

// Value implements schema.SerializerInterface.
func (g *StringSlice[T]) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	return strings.Join(lo.Map([]T(*g), func(item T, _ int) string {
		return strconv.Itoa(int(item))
	}), ","), nil
}
