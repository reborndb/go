// Copyright 2015 Reborndb Org. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

package handler

import (
	"reflect"
	"strings"

	"github.com/reborndb/go/errors"
	"github.com/reborndb/go/log"
	"github.com/reborndb/go/redis/resp"
)

type HandlerFunc func(arg0 interface{}, args ...[]byte) (resp.Resp, error)

type HandlerTable map[string]HandlerFunc

func NewHandlerTable(o interface{}, f func(string) bool) (map[string]HandlerFunc, error) {
	if o == nil {
		return nil, errors.New("handler is nil")
	}
	t := make(map[string]HandlerFunc)
	r := reflect.TypeOf(o)
	for i := 0; i < r.NumMethod(); i++ {
		m := r.Method(i)
		if f != nil && f(m.Name) {
			continue
		}

		n := strings.ToLower(m.Name)
		if h, err := createHandlerFunc(o, &m.Func); err != nil {
			return nil, err
		} else if _, exists := t[n]; exists {
			return nil, errors.Errorf("func.name = '%s' has already exists", m.Name)
		} else {
			t[n] = h
		}
	}
	return t, nil
}

func MustHandlerTable(o interface{}, f func(string) bool) map[string]HandlerFunc {
	t, err := NewHandlerTable(o, f)
	if err != nil {
		log.PanicError(err, "create redis handler map failed")
	}
	return t
}

func createHandlerFunc(o interface{}, r *reflect.Value) (HandlerFunc, error) {
	t := r.Type()
	arg0Type := reflect.TypeOf((*interface{})(nil)).Elem()
	argsType := reflect.TypeOf([][]byte{})
	if t.NumIn() != 3 || t.In(1) != arg0Type || t.In(2) != argsType {
		return nil, errors.Errorf("register with invalid func type = '%s'", t)
	}
	ret0Type := reflect.TypeOf((*resp.Resp)(nil)).Elem()
	ret1Type := reflect.TypeOf((*error)(nil)).Elem()
	if t.NumOut() != 2 || t.Out(0) != ret0Type || t.Out(1) != ret1Type {
		return nil, errors.Errorf("register with invalid func type = '%s'", t)
	}
	return func(arg0 interface{}, args ...[]byte) (resp.Resp, error) {
		var arg0Value reflect.Value
		if arg0 == nil {
			arg0Value = reflect.ValueOf((*interface{})(nil))
		} else {
			arg0Value = reflect.ValueOf(arg0)
		}
		var input, output []reflect.Value
		input = []reflect.Value{reflect.ValueOf(o), arg0Value, reflect.ValueOf(args)}
		if t.IsVariadic() {
			output = r.CallSlice(input)
		} else {
			output = r.Call(input)
		}
		var ret0 resp.Resp
		var ret1 error
		if i := output[0].Interface(); i != nil {
			ret0 = i.(resp.Resp)
		}
		if i := output[1].Interface(); i != nil {
			ret1 = i.(error)
		}
		return ret0, ret1
	}, nil
}
