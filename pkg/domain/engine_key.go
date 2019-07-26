package domain

type KeyEngine interface {
	Delete(keys []string) int
	Expire(key string, seconds int) (int, error)
	Ttl(key string) int
}

func (e *engine) Delete(keys []string) int {
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

func (e *engine) Expire(key string, seconds int) (int, error) {
	value := e.getKeyCheckExpire(key)
	if value == nil {
		return 0, nil
	}

	e.saveKeyExpire(key, seconds)

	return 1, nil
}

func (e *engine) Ttl(key string) int {
	return e.getTtl(key)
}
