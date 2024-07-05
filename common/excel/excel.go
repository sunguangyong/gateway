package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"reflect"
	"strconv"
)

// ExportByStruct excel导出(数据源为Struct)
func ExportByStruct(c io.Writer, titleList []string, data []interface{}, sheetName string) error {
	f := excelize.NewFile()

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Family: "宋体",
			Size:   11,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	contentStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "宋体",
			Size:   11,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	if err != nil {
		return err
	}

	f.SetSheetName("Sheet1", sheetName)

	_ = f.SetSheetRow(sheetName, "A1", &titleList)

	for _, colName := range letter(len(titleList)) {
		_ = f.SetCellStyle(sheetName, colName+"1", colName+"1", headerStyle)
	}
	_ = f.SetRowHeight(sheetName, 1, 30)

	lastColumn := letter(len(titleList))[len(titleList)-1]
	_ = f.SetColWidth(sheetName, "A", lastColumn, 30)

	rowNum := 1
	for _, v := range data {
		t := reflect.TypeOf(v)
		value := reflect.ValueOf(v)
		row := make([]interface{}, 0)
		for l := 0; l < t.NumField(); l++ {
			val := value.Field(l).Interface()
			row = append(row, val)
		}
		rowNum++
		err := f.SetSheetRow(sheetName, "A"+strconv.Itoa(rowNum), &row)
		_ = f.SetCellStyle(sheetName, fmt.Sprintf("A%d", rowNum), fmt.Sprintf("%s%d", lastColumn, rowNum), contentStyle)
		if err != nil {
			return err
		}
	}
	return f.Write(c)
}

// 遍历a-z
func letter(length int) []string {
	var str []string
	for i := 0; i < length; i++ {
		str = append(str, string(rune('A'+i)))
	}
	return str
}
