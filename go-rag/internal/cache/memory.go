package cache

var store = map[string]string{}

func Get(key string) (string, bool) {
	val, ok := store[key]
	return val, ok
}

func Set(key, value string) {
	store[key] = value
}
