package domain

type ListEngine interface {
	LLeftPop(key string) (string, error)
	LRightPop(key string) (string, error)
	LRightPush(key string, values []string) (int, error)
	LLeftPush(key string, values []string) (int, error)
	LRange(key string, start int, end int) ([]string, error)
}

func (e *engine) LLeftPop(key string) (string, error) {
	value := e.storage[key]
	if value == nil {
		return "", nil
	}

	listValue, err := castToList(key, value)
	if err != nil {
		return "", err
	}

	if len(listValue) == 0 {
		return "", nil
	}

	leftmost := listValue[0]

	e.storage[key] = listValue[1:]

	return leftmost, nil
}

func (e *engine) LRightPop(key string) (string, error) {
	value := e.storage[key]
	if value == nil {
		return "", nil
	}

	listValue, err := castToList(key, value)
	if err != nil {
		return "", err
	}

	listLen := len(listValue)
	if listLen == 0 {
		return "", nil
	}

	rightmost := listValue[listLen-1]

	e.storage[key] = listValue[:listLen-1]

	return rightmost, nil
}

func (e *engine) LLeftPush(key string, values []string) (int, error) {
	entry := e.getOrCreateList(key)

	listValue, err := castToList(key, entry)
	if err != nil {
		return -1, err
	}

	listValue = append(reverseSlice(values), listValue...)

	e.storage[key] = listValue

	listLen := len(listValue)

	return listLen, nil
}

func (e *engine) LRightPush(key string, values []string) (int, error) {
	entry := e.getOrCreateList(key)

	listValue, err := castToList(key, entry)
	if err != nil {
		return -1, err
	}

	listValue = append(listValue, values...)

	e.storage[key] = listValue

	listLen := len(listValue)

	return listLen, nil
}

func (e *engine) LRange(key string, start int, end int) ([]string, error) {
	entry := e.getOrCreateList(key)

	listValue, err := castToList(key, entry)
	if err != nil {
		return nil, err
	}

	listLen := len(listValue)

	if start > listLen-1 {
		return []string{}, nil
	}

	if start < 0 {
		start = listLen + start
	}

	// The rightmost item is included, so it's [start, end]
	// Negative end count backwards from list end
	// Beyond list end numbers count as list end
	if end >= listLen {
		end = listLen
	} else if end < 0 {
		end = listLen + end + 1
	} else {
		end++
	}

	if start > end {
		return []string{}, nil
	}

	return listValue[start:end], nil
}

func (e *engine) getOrCreateList(key string) interface{} {
	currentEntry := e.storage[key]
	if currentEntry == nil {
		currentEntry = []string{}
	}
	return currentEntry
}

func castToList(key string, listInterface interface{}) ([]string, error) {
	currentList, ok := listInterface.([]string)
	if !ok {
		return nil, Errorf(CodeWrongType, "%s is not of list type", key)
	}
	return currentList, nil
}

func reverseSlice(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
