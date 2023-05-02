package ade_linter

func Map[T any, U any](slice []T, mapFn func(T) U) []U {
	mappedSlice := make([]U, len(slice))

	for index, item := range slice {
		mappedSlice[index] = mapFn(item)
	}

	return mappedSlice
}

func Reverse[T any](slice []T) []T {
	newSlice := make([]T, len(slice))

	for i := 0; i < len(slice); i++ {
		newSlice[i] = slice[len(slice)-1-i]
	}

	return newSlice
}
