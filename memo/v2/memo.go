package v2

type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{}
}

type cacheFunc func(key string) (interface{}, error)
