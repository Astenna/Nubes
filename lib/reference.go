package lib

type Reference[T Nobject] string

func (r Reference[T]) Id() string {
	return string(r)
}

func (r Reference[T]) Get() (*T, error) {
	return GetObjectState[T](string(r))
}
