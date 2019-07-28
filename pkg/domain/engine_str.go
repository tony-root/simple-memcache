package domain

type StringEngine interface {
	SSet(key string, value string)
	SGet(key string) (string, error)
}

func (e *engine) SSet(key string, value string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.setString(key, value)
}

func (e *engine) SGet(key string) (string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.getString(key)
}
