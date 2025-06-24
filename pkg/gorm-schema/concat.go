package gormschema

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"gorm.io/gorm/schema"
)

func init() {
	schema.RegisterSerializer("intslice", StringSliceSerializer{})
}

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
	case []byte:
		// Handle byte array by converting to string first
		strValue := string(v)
		*g = lo.Map(strings.Split(strValue, ","), func(item string, _ int) T {
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

// StringSliceSerializer 字符串切片序列化器
// 将 "1,2,3" 转换为 []int64{1,2,3}
type StringSliceSerializer struct{}

func (StringSliceSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	fieldValue := reflect.New(field.FieldType)

	if dbValue != nil {
		switch v := dbValue.(type) {
		case string:
			intSlice := lo.Map(strings.Split(v, ","), func(item string, _ int) int64 {
				vv, _ := strconv.ParseInt(item, 10, 64)
				return vv
			})
			fieldValue.Elem().Set(reflect.ValueOf(intSlice))
		case []byte:
			// Handle byte array by converting to string first
			strValue := string(v)
			intSlice := lo.Map(strings.Split(strValue, ","), func(item string, _ int) int64 {
				vv, _ := strconv.ParseInt(item, 10, 64)
				return vv
			})
			fieldValue.Elem().Set(reflect.ValueOf(intSlice))
		default:
			return fmt.Errorf("failed to unmarshal IntSlice value: %#v", dbValue)
		}
	}

	field.ReflectValueOf(ctx, dst).Set(fieldValue.Elem())
	return
}

// 实现 Value 方法
func (StringSliceSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	return strings.Join(lo.Map(fieldValue.([]int64), func(item int64, _ int) string {
		return strconv.Itoa(int(item))
	}), ","), nil
}
