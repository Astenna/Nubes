package lib

// Reference type to define composite types
// in while defining types and its methods
// in the project where lambdas definitions
// are to be generated
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
	return nil
}

func (r Reference[T]) Get() (*T, error) {
	return Get[T](r.Id)
}
