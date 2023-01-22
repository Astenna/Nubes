package lib

type Reference[T Nobject] string

func NewReference[T Nobject](id string) *Reference[T] {
	result := Reference[T](id)
	return &result
}

func (r Reference[T]) Id() string {
	return string(r)
}

func (r Reference[T]) Get() (*T, error) {
	return Load[T](string(r))
}
