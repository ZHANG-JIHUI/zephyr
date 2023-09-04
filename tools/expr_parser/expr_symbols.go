package expr_parser

var symbols = map[string]bool{
	"(":  true,
	")":  true,
	"&&": true,
	"||": true,
	",":  true,
}
