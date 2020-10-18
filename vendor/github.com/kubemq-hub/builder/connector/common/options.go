package common

type DefaultOptions map[string][]string

func NewDefaultOptions() DefaultOptions {
	return map[string][]string{}
}

func (do DefaultOptions) Add(key string, value []string) DefaultOptions {
	do[key] = value
	return do
}
