package domain

type KeyEngine interface {
	Delete(keys []string) int
	Expire(key string, seconds int) int
	Ttl(key string) int
}

func (e *engine) Delete(keys []string) int {
	e.mu.Lock()
	defer e.mu.Unlock()

	numDeleted := 0
	for _, key := range keys {
		value := e.getKeyCheckExpire(key)
		if value != nil {
			numDeleted++
		}

		e.setKeyClearExpire(key, nil)
	}
	return numDeleted
}

func (e *engine) Expire(key string, seconds int) int {
	e.mu.Lock()
	defer e.mu.Unlock()

	value := e.getKeyCheckExpire(key)
	if value == nil {
		return 0
	}

	e.saveKeyExpire(key, seconds)

	return 1
}

func (e *engine) Ttl(key string) int {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.getTtl(key)
}
