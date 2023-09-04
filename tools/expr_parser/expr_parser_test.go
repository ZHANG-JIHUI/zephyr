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

func (slf *ExprProcessor) Verify(method string, args []any) {
	handle, ok := slf.handles[method]
	if !ok {
		log.Panic("ExprProcessor.Check: method not found", log.String("method", method))
	}
	handle(args)
}

func (slf *ExprProcessor) Result(method string, args []any) bool {
	handle, ok := slf.handles[method]
	if !ok {
		log.Error("ExprProcessor.Result: method not found", log.String("method", method), log.Any("args", args))
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
		return true
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

	expr := "(a(1,2)||b(3))&&(c(sleep)&&d(5))"
	parser := NewExprParser(expr, proc)
	tree, err := parser.BuildTree()
	if err != nil {
		panic(err)
	}
	fmt.Println(tree.String())

	if err := parser.Verify(); err != nil {
		panic(err)
	}
	result := parser.Result()
	fmt.Println(result)
}
