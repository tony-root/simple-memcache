package domain

type StringEngine interface {
	SSet(key string, value string)
	SGet(key string) (string, error)
}

func (e *engine) SSet(key string, value string) {
	e.storage[key] = value
}

func (e *engine) SGet(key string) (string, error) {
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
