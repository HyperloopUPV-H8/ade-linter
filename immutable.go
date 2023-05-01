package ade_linter

func NewMap[K comparable, V any](oldMap map[K]V) map[K]V {
	newMap := make(map[K]V, len(oldMap))

	for key, value := range oldMap {
		newMap[key] = value
	}

	return newMap
}
