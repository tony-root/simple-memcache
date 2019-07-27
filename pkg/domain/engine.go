package domain

import "time"

type expirationUnixNano int64

func unixNano(time time.Time) expirationUnixNano {
	return expirationUnixNano(time.UnixNano())
}

func (ex expirationUnixNano) seconds() int {
	return int(ex / 1000000000)
}

type engine struct {
	storage map[string]interface{}
	expires map[string]expirationUnixNano
}

func NewEngine() *engine {
	return &engine{
		storage: map[string]interface{}{},
		expires: map[string]expirationUnixNano{},
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

	if expiresAt < unixNano(time.Now()) {
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
	e.expires[key] = unixNano(time.Now().Add(time.Duration(seconds) * time.Second))
}

func (e *engine) getTtl(key string) int {
	value := e.getKeyCheckExpire(key)
	if value == nil {
		return -2
	}

	if _, ok := e.expires[key]; !ok {
		return -1
	}

	return (e.expires[key] - unixNano(time.Now())).seconds()
}

func (e *engine) setString(key, value string) {
	e.setKeyClearExpire(key, value)
}

func (e *engine) getString(key string) (string, error) {
	value := e.getKeyCheckExpire(key)
	if value == nil {
		return "", nil
	}

	strValue, ok := value.(string)
	if !ok {
		return "", Errorf(CodeWrongType, "%s is not of string type", key)
	}

	return strValue, nil
}

func (e *engine) getOrCreateList(key string) ([]string, error) {
	listInterface := e.getKeyCheckExpire(key)
	if listInterface == nil {
		return []string{}, nil
	}

	return castToList(key, listInterface)
}

func (e *engine) getList(key string) ([]string, error) {
	listInterface := e.getKeyCheckExpire(key)
	if listInterface == nil {
		return nil, nil
	}

	return castToList(key, listInterface)
}

func castToList(key string, listInterface interface{}) ([]string, error) {
	if listInterface == nil {
		return nil, nil
	}

	listString, ok := listInterface.([]string)
	if !ok {
		return nil, Errorf(CodeWrongType, "%s is not of list type", key)
	}

	return listString, nil
}

func (e *engine) setList(key string, l []string) {
	e.storage[key] = l
}

func (e *engine) getOrCreateMap(key string) (map[string]string, error) {
	mapInterface := e.getKeyCheckExpire(key)
	if mapInterface == nil {
		return map[string]string{}, nil
	}

	return castToMap(key, mapInterface)
}

func (e *engine) getMap(key string) (map[string]string, error) {
	mapInterface := e.getKeyCheckExpire(key)
	if mapInterface == nil {
		return nil, nil
	}

	return castToMap(key, mapInterface)
}

func (e *engine) setMap(key string, m map[string]string) {
	e.storage[key] = m
}

func castToMap(key string, mapInterface interface{}) (map[string]string, error) {
	if mapInterface == nil {
		return nil, nil
	}

	mapString, ok := mapInterface.(map[string]string)
	if !ok {
		return nil, Errorf(CodeWrongType, "%s is not of map type", key)
	}

	return mapString, nil
}
