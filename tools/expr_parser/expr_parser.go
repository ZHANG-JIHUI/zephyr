package expr_parser

import (
	"fmt"
	"github.com/ZHANG-JIHUI/zephyr/tools/log"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

// ExprParser 表达式解析器
type ExprParser struct {
	tokens  []string
	current int
	proc    Processor
	tree    *ExprTree
}

func NewExprParser(expr string, proc Processor) *ExprParser {
	return &ExprParser{
		tokens: tokenize(expr),
		proc:   proc,
	}
}

// BuildTree 构建表达式树
func (slf *ExprParser) BuildTree() (tree *ExprTree, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("build expr tree failed, err: %v", r)
			log.Error("ExprParser.BuildTree: panic", log.Any("err", err))
		}
	}()
	tree = &ExprTree{
		root: slf.parseOr(),
	}
	slf.tree = tree
	return
}

func (slf *ExprParser) parseOr() *ExprNode {
	result := slf.parseAnd()
	for slf.match("||") {
		andNode := slf.parseAnd()
		orNode := &ExprNode{Method: "||"}
		orNode.AddChild(result)
		orNode.AddChild(andNode)
		result = orNode
	}
	return result
}

func (slf *ExprParser) parseAnd() *ExprNode {
	result := slf.parsePrimary()
	for slf.match("&&") {
		primaryNode := slf.parsePrimary()
		andNode := &ExprNode{Method: "&&"}
		andNode.AddChild(result)
		andNode.AddChild(primaryNode)
		result = andNode
	}
	return result
}

func (slf *ExprParser) parsePrimary() *ExprNode {
	if slf.match("(") {
		orNode := slf.parseOr()
		slf.consume(")")
		return orNode
	}
	return slf.parseCondition()
}

func (slf *ExprParser) parseCondition() *ExprNode {
	method := slf.consumeIdentifier()
	slf.consume("(")
	args := slf.parseArguments()
	slf.consume(")")
	node := &ExprNode{Method: method, Arguments: args}
	return node
}

func (slf *ExprParser) parseArguments() []any {
	var arguments []any
	if !slf.check(")") {
		arguments = append(arguments, slf.consumeString())
		for slf.match(",") {
			argument := slf.consumeString()
			switch arg := argument.(type) {
			case string:
				if symbols[arg] {
					panic(fmt.Sprintf("Expected argument at position %d, but got %s", slf.current, arg))
				}
			default:
				arguments = append(arguments, arg)
			}
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
		str := slf.tokens[slf.current]
		num, err := strconv.Atoi(str)
		if err == nil {
			slf.advance()
			return num
		}
		slf.advance()
		return str
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

func tokenize(expr string) []string {
	for symbol := range symbols {
		expr = strings.ReplaceAll(expr, symbol, " "+symbol+" ")
	}
	return strings.Fields(expr)
}

// Verify 验证表达式处理器
func (slf *ExprParser) Verify() error {
	if slf.tree == nil {
		return errors.New("expr tree is nil")
	}
	return slf.verify(slf.tree.root)
}

func (slf *ExprParser) verify(node *ExprNode) (err error) {
	if slf.proc == nil {
		return errors.New("processor is nil")
	}
	if node.Method == "||" || node.Method == "&&" {
		for _, child := range node.Children {
			if err = slf.verify(child); err != nil {
				return err
			}
		}
		return nil
	} else {
		defer func() {
			if r := recover(); r != nil {
				err = errors.Errorf("method %s verify failed, args: %v, err: %v", node.Method, node.Arguments, r)
				log.Error("ExprParser.verify: panic",
					log.String("method", node.Method), log.Any("args", node.Arguments), log.Any("err", r))
			}
		}()
		slf.proc.Verify(node.Method, node.Arguments)
		return nil
	}
}

// Result 计算表达式结果
func (slf *ExprParser) Result() bool {
	if slf.tree == nil {
		return false
	}
	return slf.result(slf.tree.root)
}

func (slf *ExprParser) result(node *ExprNode) bool {
	if slf.proc == nil {
		return false
	}
	if node.Method == "||" {
		for _, child := range node.Children {
			if slf.result(child) {
				return true
			}
		}
		return false
	} else if node.Method == "&&" {
		for _, child := range node.Children {
			if !slf.result(child) {
				return false
			}
		}
		return true
	} else {
		defer func() {
			if err := recover(); err != nil {
				log.Error("ExprParser.result: panic",
					log.String("method", node.Method), log.Any("args", node.Arguments), log.Any("err", err))
			}
		}()
		return slf.proc.Result(node.Method, node.Arguments)
	}
}
