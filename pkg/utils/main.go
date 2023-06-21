package utils

func Contains(strings *[]string, query string) bool {
	for _, item := range *strings {
		if item == query {
			return true
		}
	}

	return false
}
