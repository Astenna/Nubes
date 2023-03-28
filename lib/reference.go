package lib

type Reference[T Nobject] string

func NewReference[T Nobject](id string) *Reference[T] {
	result := Reference[T](id)
	return &result
}

func (r Reference[T]) Id() string {
	return string(r)
}

func (r Reference[T]) GetLoaded() (*T, error) {
	return Load[T](string(r))
}

func (r Reference[T]) GetWithoutLoading() (*T, error) {
	object := new(T)
	err := GetObjectState(string(r), object)
	return object, err
}
