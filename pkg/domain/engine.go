package domain

type engine struct {
	storage map[string]interface{}
}

func NewEngine() *engine {
	return &engine{storage: map[string]interface{}{}}
}
