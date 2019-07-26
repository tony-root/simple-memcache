package domain

type KeyEngine interface {
	Delete(keys []string) int
	Expire(key string, seconds int) (int, error)
}

func (e *engine) Delete(keys []string) int {
	numDeleted := 0
	for _, key := range keys {
		value := e.storage[key]
		if value != nil {
			numDeleted++
		}
		e.storage[key] = nil
	}
	return numDeleted
}

func (e *engine) Expire(key string, seconds int) (int, error) {
	panic("not implemented")
}
