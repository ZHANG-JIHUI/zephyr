package reflects

import (
	"errors"
	"fmt"
	"reflect"
)

type FuncInfo struct {
	Type       reflect.Type
	Value      reflect.Value
	InArgs     []reflect.Type
	InArgsLen  int
	OutArgs    []reflect.Type
	OutArgsLen int
}

func (slf *FuncInfo) String() string {
	return fmt.Sprintf("FuncInfo{Type:%v, Value:%v, InArgs:%v, InArgsLen:%d, OutArgs:%v, OutArgsLen:%d}", slf.Type, slf.Value.String(), slf.InArgs, slf.InArgsLen, slf.OutArgs, slf.OutArgsLen)
}

func GetFuncInfo(fn any) (*FuncInfo, error) {
	if fn == nil {
		return nil, errors.New("fn is nil")
	}

	typ := reflect.TypeOf(fn)
	if typ.Kind() != reflect.Func {
		return nil, errors.New("fn is not func type")
	}

	var inArgs []reflect.Type
	for i := 0; i < typ.NumIn(); i++ {
		in := typ.In(i)
		inArgs = append(inArgs, in)
	}

	var outArgs []reflect.Type
	for i := 0; i < typ.NumOut(); i++ {
		out := typ.Out(i)
		outArgs = append(outArgs, out)
	}

	funcInfo := &FuncInfo{
		Type:       typ,
		Value:      reflect.ValueOf(fn),
		InArgs:     inArgs,
		InArgsLen:  typ.NumIn(),
		OutArgs:    outArgs,
		OutArgsLen: typ.NumOut(),
	}
	return funcInfo, nil
}
