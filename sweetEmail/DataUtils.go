package sweetEmail

import "strings"

func isEmpty(str string) bool {
	if len(str) == 0 {
		return true
	}
	upperCaseString := strings.ToUpper(str)
	if upperCaseString == "NULL" {
		return true
	}
	return false
}

func isNotEmpty(str string) bool {
	return !isEmpty(str)
}
