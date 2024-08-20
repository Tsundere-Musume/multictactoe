package main

func any[T comparable](arr []T, predicate func(int, T) bool) bool {
	for idx, val := range arr {
		if predicate(idx, val) {
			return true
		}
	}
	return false
}
