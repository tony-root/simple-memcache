package domain

import (
	"sync"
	"time"
)

type engine struct {
	storage map[string]interface{}
	expires map[string]time.Time
	mu      sync.Mutex
}

func NewEngine() *engine {
	return &engine{
		storage: map[string]interface{}{},
		expires: map[string]time.Time{},
	}
}

func (e *engine) getKeyCheckExpire(key string) interface{} {
	value := e.storage[key]
	if value == nil {
		return nil
	}

	expiresAt, ok := e.expires[key]
	if !ok {
		return value
	}

	if expiresAt.Before(time.Now()) {
		e.setKeyClearExpire(key, nil)
		return nil
	}

	return value
}

func (e *engine) setKeyClearExpire(key string, value interface{}) {
	e.storage[key] = value
	if _, ok := e.expires[key]; !ok {
		delete(e.expires, key)
	}
}

func (e *engine) saveKeyExpire(key string, seconds int) {
	e.expires[key] = time.Now().Add(time.Duration(seconds) * time.Second)
}

func (e *engine) getTtl(key string) (int, error) {
	value := e.getKeyCheckExpire(key)
	if value == nil {
		return -1, ErrTtlKeyNotFound
	}

	if _, ok := e.expires[key]; !ok {
		return -1, ErrNoTtlForKey
	}

	return int(e.expires[key].Sub(time.Now()).Seconds()), nil
}

func (e *engine) setString(key, value string) {
	e.setKeyClearExpire(key, value)
}

func (e *engine) getString(key string) (string, error) {
	value := e.getKeyCheckExpire(key)
	if value == nil {
		return "", nil
	}

	casted, ok := value.(string)
	if !ok {
		return "", errTypeMismatch(key, "string")
	}

	return casted, nil
}

func (e *engine) getOrCreateList(key string) ([]string, error) {
	value := e.getKeyCheckExpire(key)
	if value == nil {
		return []string{}, nil
	}

	return castToList(key, value)
}

func (e *engine) getList(key string) ([]string, error) {
	return castToList(key, e.getKeyCheckExpire(key))
}

func castToList(key string, value interface{}) ([]string, error) {
	if value == nil {
		return nil, nil
	}

	casted, ok := value.([]string)
	if !ok {
		return nil, errTypeMismatch(key, "list")
	}

	return casted, nil
}

func (e *engine) setList(key string, l []string) {
	e.storage[key] = l
}

func (e *engine) getOrCreateMap(key string) (map[string]string, error) {
	value := e.getKeyCheckExpire(key)
	if value == nil {
		return map[string]string{}, nil
	}

	return castToMap(key, value)
}

func (e *engine) getMap(key string) (map[string]string, error) {
	return castToMap(key, e.getKeyCheckExpire(key))
}

func (e *engine) setMap(key string, m map[string]string) {
	e.storage[key] = m
}

func castToMap(key string, value interface{}) (map[string]string, error) {
	if value == nil {
		return nil, nil
	}

	casted, ok := value.(map[string]string)
	if !ok {
		return nil, errTypeMismatch(key, "map")
	}

	return casted, nil
}
