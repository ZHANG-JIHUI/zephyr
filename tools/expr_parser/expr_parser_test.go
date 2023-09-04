package expr_parser

import (
	"fmt"
	"github.com/ZHANG-JIHUI/zephyr/tools/fn"
	"github.com/ZHANG-JIHUI/zephyr/tools/log"
	"testing"
)

type ExprProcessor struct {
	handles map[string]fn.HandleReturnBool
}

func (slf *ExprProcessor) Register(method string, handle fn.HandleReturnBool) {
	slf.handles[method] = handle
}

func (slf *ExprProcessor) Check(method string, args []any) bool {
	_, ok := slf.handles[method]
	if !ok {
		log.Panic("ExprProcessor.Check: method not found", log.String("method", method))
		return false
	}
	return true
}

func (slf *ExprProcessor) Verify(method string, args []any) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Error("ExprProcessor.Verify: panic", log.String("method", method), log.Any("args", args), log.Any("err", err))
		}
	}()

	handle, ok := slf.handles[method]
	if !ok {
		log.Error("ExprProcessor.Verify: method not found", log.String("method", method), log.Any("args", args))
		return false
	}
	return handle(args)
}

func TestExprParser(t *testing.T) {
	proc := &ExprProcessor{make(map[string]fn.HandleReturnBool)}
	proc.Register("a", func(args []any) bool {
		arg1 := args[0].(int)
		arg2 := args[1].(int)
		fmt.Println("method A", arg1+arg2)
		return false
	})
	proc.Register("b", func(args []any) bool {
		fmt.Println("method B")
		return true
	})
	proc.Register("c", func(args []any) bool {
		fmt.Println("method C")
		return true
	})
	proc.Register("d", func(args []any) bool {
		fmt.Println("method D")
		return true
	})

	expr := "(a(1,2)||b(3))&&(c(4)&&d(5))"
	parser := NewExprParser(expr, proc)
	result := parser.Parse()
	fmt.Println(result)
}
