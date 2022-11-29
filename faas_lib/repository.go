package FaaSLib

type Repository[T any] struct {
}

func (Repository[T]) Create(objToInsert T) error {
	return nil
}

func (Repository[T]) Delete(id int) error {
	return nil
}

func (Repository[T]) Get(id int) (*T, error) {
	return new(T), nil
}
