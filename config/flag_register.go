package config

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"time"
)

var (
	boolType     = reflect.TypeOf(bool(true))
	intType      = reflect.TypeOf(int(0))
	uintType     = reflect.TypeOf(uint(0))
	int32Type    = reflect.TypeOf(int32(0))
	uint32Type   = reflect.TypeOf(uint32(0))
	int64Type    = reflect.TypeOf(int64(0))
	uint64Type   = reflect.TypeOf(uint64(0))
	float64Type  = reflect.TypeOf(float64(0))
	stringType   = reflect.TypeOf(string(""))
	durationType = reflect.TypeOf(time.Duration(0))
)

func joinName(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return fmt.Sprintf("%s.%s", prefix, name)
}

func registerFlags(name, usage string, v reflect.Value, flagSet *flag.FlagSet) error {
	if !v.CanAddr() {
		return errors.New("struct field can not address")
	}
	t := v.Type()
	va := v.Addr()

	switch t {
	case boolType:
		flagSet.BoolVar(va.Interface().(*bool), name, v.Bool(), usage)
	case intType:
		flagSet.IntVar(va.Interface().(*int), name, int(v.Int()), usage)
	case uintType:
		flagSet.UintVar(va.Interface().(*uint), name, uint(v.Uint()), usage)
	case int64Type:
		flagSet.Int64Var(va.Interface().(*int64), name, v.Int(), usage)
	case uint64Type:
		flagSet.Uint64Var(va.Interface().(*uint64), name, v.Uint(), usage)
	case stringType:
		flagSet.StringVar(va.Interface().(*string), name, v.String(), usage)
	case int32Type:
		flagSet.Var(newint32Value(int32(v.Int()), va.Interface().(*int32)), name, usage)
	case uint32Type:
		flagSet.Var(newuint32Value(uint32(v.Uint()), va.Interface().(*uint32)), name, usage)
	case float64Type:
		flagSet.Float64Var(va.Interface().(*float64), name, v.Float(), usage)
	case durationType:
		flagSet.DurationVar(va.Interface().(*time.Duration), name, v.Interface().(time.Duration), usage)
	default:
		if t.Kind() == reflect.Struct {
			for idx := 0; idx < v.NumField(); idx++ {
				field := reflect.Indirect(v.Field(idx))
				fieldTag := t.Field(idx).Tag
				fieldName := fieldTag.Get("config")
				if fieldName == "" {
					continue
				}
				usage = fieldTag.Get("usage")
				if err := registerFlags(joinName(name, fieldName), usage, field, flagSet); err != nil {
					return err
				}
			}
		} else {
			return fmt.Errorf("unknown field type %s %s", name, t.Name())
		}
	}
	return nil
}

func registerFlagSet(params interface{}, flagSet *flag.FlagSet) error {
	v := reflect.ValueOf(params)
	if v.Kind() != reflect.Ptr {
		return errors.New("registerFlagSet params must be a pointer to a struct")
	}
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return errors.New("registerFlagSet params must be a pointer to a struct")
	}
	return registerFlags("", "", v, flagSet)
}
