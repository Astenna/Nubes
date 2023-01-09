package lib

type Reference[T Nobject] struct {
	Id string
}

func NewReference[T Nobject](id string) *Reference[T] {
	if id != "" {
		newObj := &Reference[T]{
			Id: id,
		}
		return newObj
	}
	return new(Reference[T])
}

func (r Reference[T]) Get() (*T, error) {
	return Get[T](r.Id)
}
