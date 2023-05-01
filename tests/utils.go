package tests

type MockTest[T any, U any] struct {
	Prompt T
	Want   U
}

func AssertMocks[T any](mocks []MockTest[T, bool], predicate func(item T) bool, errorFn func(mock MockTest[T, bool], got bool)) {
	for _, mock := range mocks {
		got := predicate(mock.Prompt)
		if got != mock.Want {
			errorFn(mock, got)
		}
	}
}
