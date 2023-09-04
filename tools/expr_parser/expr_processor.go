package expr_parser

type Processor interface {
	Register(method string, handle func(args []any) bool)
	Verify(method string, args []any)
	Result(method string, args []any) bool
}
