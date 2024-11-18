package utils

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func Unique(s []string) []string {
	unique := make(map[string]bool, len(s))
	for _, v := range s {
		unique[v] = true
	}

	result := make([]string, 1)

	for v := range unique {
		result = append(result, v)
	}

	return result
}

func DifferenceList(a, b []string) []string {
	elements := make(map[string]struct{})

	for _, elem := range b {
		elements[elem] = struct{}{}
	}

	result := []string{}
	for _, elem := range a {
		if _, found := elements[elem]; !found {
			result = append(result, elem)
		}
	}

	return result
}
