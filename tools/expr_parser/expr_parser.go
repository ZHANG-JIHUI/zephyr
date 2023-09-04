package expr_parser

import (
	"fmt"
	"strconv"
	"strings"
)

// ExprParser 表达式递归下降解析器
type ExprParser struct {
	tokens  []string
	current int
	proc    Processor
}

type (
	Processor interface {
		Register(method string, handle func(args []any) bool)
		Check(method string, args []any) bool
		Verify(method string, args []any) bool
	}
)

func NewExprParser(expr string, proc Processor) *ExprParser {
	tokens := exprTokenize(expr)
	return &ExprParser{
		tokens:  tokens,
		current: 0,
		proc:    proc,
	}
}

func (slf *ExprParser) Parse() bool {
	return slf.parseOr()
}

func (slf *ExprParser) parseOr() bool {
	result := slf.parseAnd()
	for slf.match("||") {
		result = result || slf.parseAnd()
	}
	return result
}

func (slf *ExprParser) parseAnd() bool {
	result := slf.parsePrimary()
	for slf.match("&&") {
		result = result && slf.parsePrimary()
	}
	return result
}

func (slf *ExprParser) parsePrimary() bool {
	if slf.match("(") {
		result := slf.parseOr()
		slf.consume(")")
		return result
	}
	return slf.parseCondition()
}

func (slf *ExprParser) parseCondition() bool {
	method := slf.consumeIdentifier()
	slf.consume("(")
	args := slf.parseArguments()
	slf.consume(")")

	return slf.proc.Verify(method, args)
}

func (slf *ExprParser) parseArguments() []any {
	var arguments []any
	if !slf.check(")") {
		arguments = append(arguments, slf.consumeString())
		for slf.match(",") {
			arguments = append(arguments, slf.consumeString())
		}
	}
	return arguments
}

func (slf *ExprParser) match(expected string) bool {
	if slf.check(expected) {
		slf.advance()
		return true
	}
	return false
}

func (slf *ExprParser) check(expected string) bool {
	if slf.isAtEnd() {
		return false
	}
	return slf.tokens[slf.current] == expected
}

func (slf *ExprParser) advance() {
	if !slf.isAtEnd() {
		slf.current++
	}
}

func (slf *ExprParser) consume(expected string) {
	if slf.check(expected) {
		slf.advance()
	} else {
		panic(fmt.Sprintf("Expected '%s' at position %d", expected, slf.current))
	}
}

func (slf *ExprParser) consumeIdentifier() string {
	if slf.checkIdentifier() {
		identifier := slf.tokens[slf.current]
		slf.advance()
		return identifier
	}
	panic(fmt.Sprintf("Expected an identifier at position %d", slf.current))
}

func (slf *ExprParser) checkIdentifier() bool {
	if slf.isAtEnd() {
		return false
	}
	token := slf.tokens[slf.current]
	return !slf.isOperator(token) && token != "(" && token != ")"
}

func (slf *ExprParser) consumeString() any {
	if slf.checkString() {
		var res any
		str := slf.tokens[slf.current]
		num, err := strconv.Atoi(str)
		if err == nil {
			res = num
		} else {
			res = str
		}
		slf.advance()
		return res
	}
	panic(fmt.Sprintf("Expected a string at position %d", slf.current))
}

func (slf *ExprParser) checkString() bool {
	if slf.isAtEnd() {
		return false
	}
	return true
}

func (slf *ExprParser) isAtEnd() bool {
	return slf.current >= len(slf.tokens)
}

func (slf *ExprParser) isOperator(token string) bool {
	return token == "||" || token == "&&"
}

func exprTokenize(expr string) []string {
	expr = strings.ReplaceAll(expr, "(", " ( ")
	expr = strings.ReplaceAll(expr, ")", " ) ")
	expr = strings.ReplaceAll(expr, "||", " || ")
	expr = strings.ReplaceAll(expr, "&&", " && ")
	expr = strings.ReplaceAll(expr, ",", " , ")
	return strings.Fields(expr)
}
