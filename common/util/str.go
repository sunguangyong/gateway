package util

import (
	"fmt"

	"github.com/axgle/mahonia"
	"github.com/google/uuid"
)

func GetUuid() string {
	return uuid.New().String()
}

func ConvertToString(src string, srcCode string, tagCode string) string {

	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func GetInteriorHost(host string, port int) string {
	return fmt.Sprintf(`http://%s:%d`, host, port)
}
