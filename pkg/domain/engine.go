package domain

type Engine interface {
	SetString(key string, value string)
	GetString(key string) (string, error)
}

type engine struct {
	storage map[string]interface{}
}

func NewEngine() *engine {
	return &engine{
		storage: map[string]interface{}{},
	}
}

func (e *engine) SetString(key string, value string) {
	e.storage[key] = value
}

func (e *engine) GetString(key string) (string, error) {
	value := e.storage[key]
	if value == nil {
		return "", nil
	}

	strValue, ok := value.(string)
	if !ok {
		return "", Errorf(CodeWrongType, "%s is not of string type", key)
	}
	return strValue, nil
}
