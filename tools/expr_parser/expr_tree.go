package expr_parser

import "encoding/json"

type ExprNode struct {
	Method    string      `json:"method"`
	Arguments []any       `json:"arguments,omitempty"`
	Children  []*ExprNode `json:"children,omitempty"`
}

func (n *ExprNode) AddChild(child *ExprNode) {
	n.Children = append(n.Children, child)
}

type ExprTree struct {
	root *ExprNode
}

func (slf *ExprTree) Bytes() []byte {
	bytes, _ := json.Marshal(slf.root)
	return bytes
}

func (slf *ExprTree) String() string {
	return string(slf.Bytes())
}
