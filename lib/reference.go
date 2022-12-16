package lib

type Reference[T Object] struct {
	instance *T `dynamodbav:"-"`
	Id       string
}

func NewReference[T Object](id string) *Reference[T] {
	if id != "" {
		newObj := &Reference[T]{
			Id: id,
		}
		return newObj
	}
	return nil
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
