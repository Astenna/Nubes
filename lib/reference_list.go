package lib

type ReferenceList[T Nobject] struct {
	Id []string
}

func NewReferenceList[T Nobject](ids []string) *ReferenceList[T] {
	if ids != nil {
		newObj := &ReferenceList[T]{
			Id: ids,
		}
		return newObj
	}
	return new(ReferenceList[T])
}

// func (r ReferenceList[T]) Get() (*T, error) {
// 	return Get[T](r.Id)
// }
