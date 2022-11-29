package FaaSLib

type Reference[T any] struct {
	instance *T
	id       string
}

func NewReference(id string) *Reference[any] {
	return &Reference[any]{id: id}
}

func (r Reference[T]) Get() T {
	if r.instance == nil {
		r.instance = new(T)
	}

	return *r.instance
}

func (r Reference[T]) Set(i int) {
	// perform update on the object
}
