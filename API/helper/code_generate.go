package helper

import "fmt"

func GenerateCode(prefix string, number uint) string {
	return fmt.Sprintf("%s%06d", prefix, number)
}
