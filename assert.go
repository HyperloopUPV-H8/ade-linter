package ade_linter

func EveryMap[K comparable, V any](m map[K]V, predicate func(key K, value V) bool) bool {
	for key, value := range m {
		if !predicate(key, value) {
			return false
		}
	}

	return true
}

func CheckAll(tests []Test) bool {
	result := true

	for _, test := range tests {
		result = test.Run() && result
	}

	return result
}
