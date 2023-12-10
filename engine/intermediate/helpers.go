package intermediate

import (
	"os"
)

func interfaceToInt(value interface{}) int32 {
	val, ok := value.(int32)
	if !ok {
		return int32(0)
	}
	return val
}

func getClientSelector() string {
	sel := os.Getenv("SELECTOR")
	if sel == "" {
		sel = "static"
	}
	return sel
}
