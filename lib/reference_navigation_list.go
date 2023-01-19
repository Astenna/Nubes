package lib

type ReferenceNavigationList[T Nobject] struct {
	ownerId            string
	ownerTypeName      string
	referringFieldName string
}

func NewReferenceNavigationList[T Nobject](ownerId, ownerTypeName, referringFieldName string) *ReferenceNavigationList[T] {
	result := new(ReferenceNavigationList[T])
	result.ownerId = ownerId
	result.ownerTypeName = ownerTypeName
	result.referringFieldName = referringFieldName
	return result
}

func (r ReferenceNavigationList[T]) GetIds() ([]string, error) {
	out, err := GetByIndex[T](r.ownerId, r.referringFieldName, r.ownerTypeName)
	return out, err
}

func (r ReferenceNavigationList[T]) Get() ([]T, error) {
	// make call using index

	return nil, nil
}
