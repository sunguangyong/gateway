package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// 读取服务器文件
func RemoteReadFile(path string) {
	//"http://82.156.56.237/data/test.txt"
	resp, err := http.Get(path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("File contents:", string(body))
}

func OpenFile(path string) {
	// 打开文件 "/Users/sunguangyong/Desktop/demo.ini"
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 创建一个 Scanner 对象，用于逐行读取文件内容
	scanner := bufio.NewScanner(file)

	// 逐行读取文件内容并输出
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// 检查是否有错误发生
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}
