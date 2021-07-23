package utils

import "github.com/spf13/cast"

type TypeTransform struct {
	Value interface{}
}

func Transform(v interface{}) *TypeTransform {
	return &TypeTransform{Value: v}
}

func (t *TypeTransform) String() string {
	return cast.ToString(t.Value)
}

func (t *TypeTransform) Int() int {
	return cast.ToInt(t.Value)
}

func (t *TypeTransform) Int64() int64 {
	return cast.ToInt64(t.Value)
}

func (t *TypeTransform) Float64() float64 {
	return cast.ToFloat64(t.Value)
}

func (t *TypeTransform) Bool() bool {
	return cast.ToBool(t.Value)
}

func (t *TypeTransform) Slice() []interface{} {
	return cast.ToSlice(t.Value)
}
