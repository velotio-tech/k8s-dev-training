package utils

func ConvertToPtrList[T any](arr []T) []*T {
	var list []*T
	for _, v := range arr {
		list = append(list, &v)
	}

	return list
}
