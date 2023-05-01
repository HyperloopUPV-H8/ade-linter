package ade_linter

func EveryMap[K comparable, V any](m map[K]V, predicate func(key K, value V) bool) bool {
	for key, value := range m {
		if !predicate(key, value) {
			return false
		}
	}

	return true
}

func Every(tests []Test) bool {
	for _, test := range tests {
		if !test.Run() {
			return false
		}
	}

	return true
}
