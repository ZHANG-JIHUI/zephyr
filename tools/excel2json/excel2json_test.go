package excel2json_test

import (
	"github.com/ZHANG-JIHUI/zephyr/tools/excel2json"
	"testing"
)

func TestExcel2Json(t *testing.T) {
	if err := excel2json.Excel2Json("./excel", "./json"); err != nil {
		panic(err)
	}
}
