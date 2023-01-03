package lib

// Reference type to define composite types
// in client_lib and client projects
// that use lambdas geneared by Nubes
type ClientReference[T Nobject] struct {
	instance *T `dynamodbav:"-"`
	Id       string
}

func NewClientReference[T Nobject](id string) *ClientReference[T] {
	if id != "" {
		newObj := &ClientReference[T]{
			Id: id,
		}
		return newObj
	}
	return nil
}

func (r ClientReference[T]) Get() *T {
	if r.instance == nil {
		r.instance = new(T)
	}

	return r.instance
}

func (r *ClientReference[T]) Set(i string) {
	r.Id = i
}
