package lib

// Reference type to define composite types
// in while defining types and its methods
// in the project where lambdas definitions
// are to be generated
type FaasReference[T Nobject] struct {
	Id string
}

func NewFaasReference[T Nobject](id string) *FaasReference[T] {
	if id != "" {
		newObj := &FaasReference[T]{
			Id: id,
		}
		return newObj
	}
	return nil
}

func (r FaasReference[T]) Get() (*T, error) {
	return Get[T](r.Id)
}

func (r *FaasReference[T]) Set(i string) {
	r.Id = i
}
