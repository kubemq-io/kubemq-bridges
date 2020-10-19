package survey

import (
	"fmt"
	"strings"
)

func ValidateNoneEmptyString(val interface{}) error {
	str, _ := val.(string)
	if str == "" {
		return fmt.Errorf("value cannot be empty")
	}
	return nil
}
func ValidateNoneSpace(val interface{}) error {
	str, _ := val.(string)
	if strings.Contains(str, " ") {
		return fmt.Errorf("value cannot contains spaces")
	}
	return nil
}
