package model

type List[T any] struct {
	Count int
	Items []T
}
