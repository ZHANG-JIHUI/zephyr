package reflects

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"testing"
)

func Sum[T constraints.Integer](a, b T) T {
	return a + b
}

func TestGetFuncInfo(t *testing.T) {
	fmt.Println(GetFuncInfo(Sum[int]))
}
