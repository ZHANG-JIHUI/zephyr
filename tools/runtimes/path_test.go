package runtimes

import (
	"fmt"
	"testing"
)

func TestGetWorkingDir(t *testing.T) {
	fmt.Println(GetWorkingDir())
}
