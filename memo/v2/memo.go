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

type request struct {
	key      string
	response chan<- result
}

type Memo struct {
	requests chan request
}

func New(f cacheFunc) *Memo {
	memo := Memo{
		requests: make(chan request),
	}

	go memo.server(f)

	return &memo
}

func (m *Memo) Get(key string) (interface{}, error) {
	response := make(chan result)
	m.requests <- request{
		key:      key,
		response: response,
	}
	res := <-response

	return res.value, res.err
}

func (m *Memo) Close() {
	close(m.requests)
}

func (m *Memo) server(f cacheFunc) {
	cache := make(map[string]*entry)
	for req := range m.requests {
		e := cache[req.key]
		if e == nil {
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key)
		}
		go e.deliver(req.response)
	}
}

func (e *entry) call(f cacheFunc, key string) {
	e.res.value, e.res.err = f(key)

	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	<-e.ready

	response <- e.res
}
