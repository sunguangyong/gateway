package util

import (
	"bufio"
	"io"
	"strings"
)

func ReadIniFile(r io.Reader) (data map[string]string) {
	// 创建一个 Scanner 对象，用于逐行读取文件内容
	data = make(map[string]string)
	scanner := bufio.NewScanner(r)

	// 逐行读取文件内容并输出
	for scanner.Scan() {
		//fmt.Println("hhhhhhhhhhhhh \n",)
		info := ConvertToString(scanner.Text(), "gbk", "utf-8")

		infoList := strings.Split(info, "=")
		if len(infoList) == 2 {
			data[infoList[0]] = infoList[1]
		}

	}
	return data
}
