package http

type Header map[string][]string

func (h Header) Add(key string, values ...string) {
	_, ok := h[key]
	if !ok {
		h.Set(key, values...)
	} else {
		h[key] = append(h[key], values...)
	}
}

func (h Header) Set(key string, values ...string) {
	h[key] = values
}

func (h Header) Get(key string) []string {
	return h[key]
}

func (h Header) Del(key string) {
	delete(h, key)
}

func (h Header) Clear() {
	for key := range h {
		delete(h, key)
	}
}
