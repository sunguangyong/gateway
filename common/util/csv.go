package util

import (
	"encoding/csv"
	"fmt"
	"io"
)

func ReadCsv(r io.Reader) (data []map[string]string) {
	//file, err := os.Open(path)
	data = []map[string]string{}
	titleMap := make(map[int]string)

	reader := csv.NewReader(r)
	rows, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return data
	}

	if len(rows) > 0 {
		for i, v := range rows[0] {
			titleMap[i] = v
		}
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		info := make(map[string]string)
		for k, v := range row {
			title, ok := titleMap[k]
			if ok {
				info[title] = ConvertToString(v, "gbk", "utf-8")
			}
		}
		data = append(data, info)
	}
	return data
}
