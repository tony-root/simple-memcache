package domain

type StringResult struct {
	Value string
}

type BulkStringResult struct {
	Value string
}

type IntResult struct {
	Value int
}

type ArrayResult struct {
	Value []string
}

type ErrorResult struct {
	Value error
}

type NilResult struct{}

var OkResult = StringResult{Value: "OK"}
