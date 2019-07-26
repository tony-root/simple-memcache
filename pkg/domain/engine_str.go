package domain

type StringEngine interface {
	SSet(key string, value string)
	SGet(key string) (string, error)
}

func (e *engine) SSet(key string, value string) {
	e.setString(key, value)
}

func (e *engine) SGet(key string) (string, error) {
	return e.getString(key)
}
