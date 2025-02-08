package model

type contextKey int

const (
	RequestIDKey contextKey = iota + 1
)

const (
	RequestIDHeader = "X-Request-ID"
)
