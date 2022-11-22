package FaaSLib

type Reference[T any] struct {
	instance *T
	Id       string
}

func (r Reference[T]) Get() T {
	if r.instance == nil {
		r.instance = new(T)
	}

	return *r.instance
}
