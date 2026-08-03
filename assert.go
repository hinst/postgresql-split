package main

func AssertError(e error) {
	if e != nil {
		panic(e)
	}
}

func AssertResultError[T any](result T, e error) T {
	AssertError(e)
	return result
}
