package utils

type array []interface{}

func InArray(arr []interface{}, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}

	return false
}
