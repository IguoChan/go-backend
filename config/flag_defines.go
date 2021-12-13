package config

import (
	"errors"
	"strconv"
)

func numError(err error) error {
	ne, ok := err.(*strconv.NumError)
	if !ok {
		return err
	}
	if ne.Err == strconv.ErrSyntax {
		return errors.New("parse error")
	}
	if ne.Err == strconv.ErrRange {
		return errors.New("value out of range")
	}
	return err
}

// -- int32 Value
type Int32Value int32

func newint32Value(val int32, p *int32) *Int32Value {
	*p = val
	return (*Int32Value)(p)
}

func (i *Int32Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 32)
	if err != nil {
		err = numError(err)
	}
	*i = Int32Value(v)
	return err
}

func (i *Int32Value) Get() interface{} { return int32(*i) }

func (i *Int32Value) String() string { return strconv.FormatInt(int64(*i), 10) }

// -- uint32 Value
type Uint32Value uint32

func newuint32Value(val uint32, p *uint32) *Uint32Value {
	*p = val
	return (*Uint32Value)(p)
}

func (i *Uint32Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 32)
	if err != nil {
		err = numError(err)
	}
	*i = Uint32Value(v)
	return err
}

func (i *Uint32Value) Get() interface{} { return uint32(*i) }

func (i *Uint32Value) String() string { return strconv.FormatUint(uint64(*i), 10) }
