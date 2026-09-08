package helper

import (
	"fmt"
	"regexp"
	"strings"
)

func GenerateEmail(name string, id int) string {
	email := strings.ToLower(name)
	email = strings.ReplaceAll(email, " ", "")
	reg := regexp.MustCompile(`[^a-z0-9]`)
	email = reg.ReplaceAllString(email, "")
	return fmt.Sprintf("%s%d@gmail.com", email, id)
}
