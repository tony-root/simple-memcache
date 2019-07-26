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
	mapString, err := e.getOrCreateMap(key)
	if err != nil {
		return -1, err
	}

	previous := mapString[field]

	mapString[field] = value

	if previous == "" {
		return 1, nil
	}

	e.setMap(key, mapString)

	return 0, nil
}

func (e *engine) HMultiSet(key string, entries map[string]string) error {
	mapString, err := e.getOrCreateMap(key)
	if err != nil {
		return err
	}

	for k, v := range entries {
		mapString[k] = v
	}

	e.setMap(key, mapString)

	return nil
}

func (e *engine) HGet(key string, field string) (string, error) {
	mapString, err := e.getMap(key)
	if err != nil {
		return "", err
	}

	if mapString == nil {
		return "", nil
	}

	return mapString[field], nil
}

func (e *engine) HMultiGet(key string, fields []string) ([]string, error) {
	mapString, err := e.getMap(key)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(fields))

	if mapString == nil {
		return result, nil
	}

	for i, field := range fields {
		result[i] = mapString[field]
	}

	return result, nil
}

func (e *engine) HGetAll(key string) ([]string, error) {
	mapString, err := e.getMap(key)
	if err != nil {
		return nil, err
	}

	if mapString == nil {
		return make([]string, 0, 0), nil
	}

	result := make([]string, 0, 2*len(mapString))
	for k, v := range mapString {
		result = append(result, k, v)
	}

	return result, nil
}

func (e *engine) HDelete(key string, fields []string) (int, error) {
	mapString, err := e.getMap(key)
	if err != nil {
		return -1, err
	}

	deleted := 0

	if mapString == nil {
		return deleted, nil
	}

	for _, field := range fields {
		if mapString[field] != "" {
			deleted++
		}
		mapString[field] = ""
	}

	e.setMap(key, mapString)

	return deleted, nil
}
