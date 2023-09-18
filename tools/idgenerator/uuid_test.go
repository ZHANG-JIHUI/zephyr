package idgenerator_test

import (
	"github.com/google/uuid"

	"fmt"
	"testing"
)

func TestUUID(t *testing.T) {
	fmt.Println(uuid.New().String())
}
