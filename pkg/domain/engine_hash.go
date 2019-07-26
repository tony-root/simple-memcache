package domain

type HashEngine interface {
	HSet(key, field, value string) (int, error)
	HMultiSet(key string, entries map[string]string) error
	HGet(key string, field string) (string, error)
	HMultiGet(key string, fields []string) ([]string, error)
	HGetAll(key string) ([]string, error)
	HDelete(key string, fields []string) (int, error)
}

func (e *engine) HSet(key, field, value string) (int, error) {
	stringMap, err := castToMap(key, e.getOrCreateMap(key))
	if err != nil {
		return -1, err
	}

	previous := stringMap[field]

	stringMap[field] = value

	e.storage[key] = stringMap

	if previous == "" {
		return 1, nil
	}

	return 0, nil
}

func (e *engine) HMultiSet(key string, entries map[string]string) error {
	stringMap, err := castToMap(key, e.getOrCreateMap(key))
	if err != nil {
		return err
	}

	for k, v := range entries {
		stringMap[k] = v
	}

	e.storage[key] = stringMap

	return nil
}

func (e *engine) HGet(key string, field string) (string, error) {
	stringMap, err := castToMap(key, e.getExistingMap(key))
	if err != nil {
		return "", err
	}

	if stringMap == nil {
		return "", nil
	}

	return stringMap[field], nil
}

func (e *engine) HMultiGet(key string, fields []string) ([]string, error) {
	stringMap, err := castToMap(key, e.getExistingMap(key))
	if err != nil {
		return nil, err
	}

	result := make([]string, len(fields))

	if stringMap == nil {
		return result, nil
	}

	for i, field := range fields {
		result[i] = stringMap[field]
	}

	return result, nil
}

func (e *engine) HGetAll(key string) ([]string, error) {
	stringMap, err := castToMap(key, e.getExistingMap(key))
	if err != nil {
		return nil, err
	}

	if stringMap == nil {
		return make([]string, 0, 0), nil
	}

	result := make([]string, 0, 2*len(stringMap))
	for k, v := range stringMap {
		result = append(result, k, v)
	}

	return result, nil
}

func (e *engine) HDelete(key string, fields []string) (int, error) {
	stringMap, err := castToMap(key, e.getExistingMap(key))
	if err != nil {
		return -1, err
	}

	deleted := 0

	if stringMap == nil {
		return deleted, nil
	}

	for _, field := range fields {
		if stringMap[field] != "" {
			deleted++
		}
		stringMap[field] = ""
	}

	e.storage[key] = stringMap

	return deleted, nil
}

func (e *engine) getOrCreateMap(key string) interface{} {
	entry := e.storage[key]
	if entry == nil {
		entry = map[string]string{}
	}
	return entry
}

func (e *engine) getExistingMap(key string) interface{} {
	return e.storage[key]
}

func castToMap(key string, mapInterface interface{}) (map[string]string, error) {
	if mapInterface == nil {
		return nil, nil
	}

	currentMap, ok := mapInterface.(map[string]string)
	if !ok {
		return nil, Errorf(CodeWrongType, "%s is not of map type", key)
	}
	return currentMap, nil
}
