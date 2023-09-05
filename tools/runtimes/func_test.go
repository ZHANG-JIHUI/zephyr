package runtimes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCurrentFuncName(t *testing.T) {
	assert.Equal(t, CurrentFuncName(),
		"github.com/ZHANG-JIHUI/zephyr/tools/runtimes.TestCurrentFuncName")
}
