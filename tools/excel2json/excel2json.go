package excel2json

import (
	"fmt"
	"github.com/ZHANG-JIHUI/zephyr/tools/files"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"os"
	"path"
)

const (
	lineSheetName = iota
	lineFieldName
	lineFieldType
	lineFieldDesc
)

const (
	SheetFieldTypeInt    = "int"
	SheetFieldTypeFloat  = "float"
	SheetFieldTypeString = "string"
)

const (
	sheetHeaderLines = 4
)

type (
	meta struct {
		Key string
		Idx int
		Typ string
	}
	line []any
)

func Excel2Json(inDir, outDir string) error {
	excels, err := files.FileList(inDir, ".xlsx", true)
	if err != nil {
		return err
	}
	for _, excel := range excels {
		xlsx, err := excelize.OpenFile(excel)
		if err != nil {
			return err
		}
		sheets := xlsx.GetSheetList()
		for idx, sheet := range sheets {
			rows, err := xlsx.GetRows(sheet)
			if err != nil {
				return err
			}
			if len(rows) < 5 {
				return errors.Errorf("excel format error, sheet:%s", sheet)
			}
			length := len(rows[lineFieldName])
			metas := make([]*meta, 0, length)
			lines := make([]line, 0, len(rows)-sheetHeaderLines)
			sheetName := xlsx.GetSheetName(idx)
			for idx, row := range rows {
				switch idx {
				case lineSheetName:
					sheetName = row[lineSheetName]
				case lineFieldName:
					for idx, colName := range row {
						fmt.Println(idx, colName, len(metas))
						metas = append(metas, &meta{Key: colName, Idx: idx})
					}
				case lineFieldType:
					for idx, typ := range row {
						metas[idx].Typ = typ
					}
				case lineFieldDesc:
				default:
					data := make(line, length)
					for k := 0; k < length; k++ {
						if k < len(row) {
							data[k] = row[k]
						}
					}
					lines = append(lines, data)
				}
			}
			outfile := fmt.Sprintf("%s.json", sheetName)
			if err = output(outDir, outfile, toJson(lines, metas)); err != nil {
				return err
			}
		}
	}
	return nil
}

func output(outDir, outfile string, str string) error {
	if err := os.MkdirAll(outDir, 0777); err != nil {
		return err
	}
	fp, err := os.OpenFile(path.Join(outDir, outfile), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer fp.Close()
	if _, err = fp.WriteString(str); err != nil {
		return err
	}
	return nil
}

func toJson(lines []line, metas []*meta) string {
	ret := "["
	for _, row := range lines {
		ret += "\n\t{"
		for idx, meta := range metas {
			ret += fmt.Sprintf("\n\t\t\"%s\":", meta.Key)
			switch meta.Typ {
			case SheetFieldTypeInt, SheetFieldTypeFloat:
				if row[idx] == nil || row[idx] == "" {
					ret += "0"
				} else {
					ret += fmt.Sprintf("%s", row[idx])
				}
			case SheetFieldTypeString:
				if row[idx] == nil {
					ret += "\"\""
				} else {
					ret += fmt.Sprintf("\"%s\"", row[idx])
				}
			default:
				return errors.Errorf("unknown field type:%s", meta.Typ).Error()
			}
			ret += ","
		}
		ret = ret[:len(ret)-1]

		ret += "\n\t},"
	}
	ret = ret[:len(ret)-1]

	ret += "\n]"
	return ret
}
