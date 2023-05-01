package ade_linter

func Map[T any, U any](slice []T, mapFn func(T) U) []U {
	mappedSlice := make([]U, len(slice))

	for index, item := range slice {
		mappedSlice[index] = mapFn(item)
	}

	return mappedSlice
}
