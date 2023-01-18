package lib

import "fmt"

type ReferenceList[T Nobject] struct {
	Ids []string
}

func NewReferenceList[T Nobject](ids []string) ReferenceList[T] {
	if ids != nil {
		newRefList := &ReferenceList[T]{
			Ids: ids,
		}
		return *newRefList
	}
	return *new(ReferenceList[T])
}

func (r ReferenceList[T]) Get() ([]T, error) {
	result := make([]T, len(r.Ids))

	for index, id := range r.Ids {
		instance, err := Get[T](id)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve object with id: %d. Error: %w", index, err)
		}
		result[index] = *instance
	}

	return result, nil
}

func (r ReferenceList[T]) GetByIndex(index int) (*T, error) {
	if len(r.Ids)-1 < index || index < 0 {
		return nil, fmt.Errorf("provided index: %d is out of bounds of the list", index)
	}
	instance, err := Get[T](r.Ids[index])
	if err != nil {
		return nil, fmt.Errorf("could not retrieve object with id: %d. Error: %w", index, err)
	}

	return instance, nil
}
