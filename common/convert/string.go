package convert

import "strconv"

func StringToInterface(input string) (interface{}, error) {
	// 尝试将输入字符串转换为 int 类型
	intValue, err := strconv.Atoi(input)
	if err == nil {
		return intValue, nil
	}

	// 尝试将输入字符串转换为 float64 类型
	floatValue, err := strconv.ParseFloat(input, 64)
	if err == nil {
		return floatValue, nil
	}

	// 尝试将输入字符串转换为 bool 类型
	boolValue, err := strconv.ParseBool(input)
	if err == nil {
		return boolValue, nil
	}

	// 如果以上都不成功,则返回原始字符串
	return input, nil
}
