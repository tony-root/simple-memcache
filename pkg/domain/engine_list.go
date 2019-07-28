package domain

type ListEngine interface {
	LLeftPop(key string) (string, error)
	LRightPop(key string) (string, error)
	LRightPush(key string, values []string) (int, error)
	LLeftPush(key string, values []string) (int, error)
	LRange(key string, start int, end int) ([]string, error)
}

func (e *engine) LLeftPop(key string) (string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	listString, err := e.getList(key)
	if err != nil {
		return "", err
	}

	if len(listString) == 0 {
		return "", nil
	}

	leftmost := listString[0]

	e.setList(key, listString[1:])

	return leftmost, nil
}

func (e *engine) LRightPop(key string) (string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	listString, err := e.getList(key)
	if err != nil {
		return "", err
	}

	listLen := len(listString)
	if listLen == 0 {
		return "", nil
	}

	rightmost := listString[listLen-1]

	e.setList(key, listString[:listLen-1])

	return rightmost, nil
}

func (e *engine) LLeftPush(key string, values []string) (int, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	listString, err := e.getOrCreateList(key)
	if err != nil {
		return -1, err
	}

	listString = append(reverseSlice(values), listString...)

	e.setList(key, listString)

	listLen := len(listString)

	return listLen, nil
}

func (e *engine) LRightPush(key string, values []string) (int, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	listString, err := e.getOrCreateList(key)
	if err != nil {
		return -1, err
	}

	listString = append(listString, values...)

	e.setList(key, listString)

	listLen := len(listString)

	return listLen, nil
}

func (e *engine) LRange(key string, start int, end int) ([]string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	listString, err := e.getList(key)
	if err != nil {
		return nil, err
	}

	listLen := len(listString)

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

	return listString[start:end], nil
}

func reverseSlice(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
