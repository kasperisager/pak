package asset

type (
	Flags []*flag

	flag struct {
		key   string
		value interface{}
	}
)

func (f Flags) Has(key string) bool {
	for _, flag := range f {
		if key == flag.key {
			return true
		}
	}

	return false
}

func (f Flags) Get(key string) interface{} {
	for _, flag := range f {
		if key == flag.key {
			return flag.value
		}
	}

	return nil
}

func (f Flags) Set(key string, value interface{}) Flags {
	for _, flag := range f {
		if key == flag.key {
			flag.value = value
			return f
		}
	}

	return append(f, &flag{key, value})
}

func (f Flags) Delete(key string) Flags {
	for i, flag := range f {
		if key == flag.key {
			last := len(f) - 1

			f[i] = f[last]
			f[last] = nil

			return f[:last]
		}
	}

	return f
}
