package files_test

import (
	"fmt"
	"testing"

	"github.com/ZHANG-JIHUI/zephyr/tools/files"
)

func TestFileList(t *testing.T) {
	filePaths, err := files.FileList("./", ".go", false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(filePaths)
}
