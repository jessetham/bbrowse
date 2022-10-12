package main

type adapter interface {
	description() string
	toString([]byte) (string, error)
}

func newAdapterList() []adapter {
	return []adapter{
		castToString{},
	}
}

type castToString struct{}

func (c castToString) description() string {
	return "Cast value to string"
}

func (c castToString) toString(b []byte) (string, error) {
	return string(b), nil
}
