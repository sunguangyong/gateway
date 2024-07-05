package convert

import "fmt"

func Float64ToStr(f float64) string {
	return fmt.Sprintf("%.8f", f)
}
