package query

type QueryCache struct {
	cache map[string]string
}

func NewQueryCache() *QueryCache {
	return &QueryCache{cache: make(map[string]string)}
}

func (qc *QueryCache) Get(query string) (string, bool) {
	result, found := qc.cache[query]
	return result, found
}

func (qc *QueryCache) Set(query, result string) {
	qc.cache[query] = result
}
