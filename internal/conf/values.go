package conf

type Values map[string][]string

func (v Values) Add(key, value string) {
	v[key] = append(v[key], value)
}

func (v Values) Set(key, value string) {
	v[key] = []string{value}
}

func (v Values) Del(key, value string) {
	delete(v, key)
}

func (v Values) Get(key string) string {
	vals := v[key]
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}
