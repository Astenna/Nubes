package lib

type ReferenceNavigationList[Many, One Nobject] struct {
	ids []string
}

func (r ReferenceNavigationList[Many, One]) GetIds() ([]string, error) {
	// make call using index
	//oneTypeName := (*new(One)).GetTypeName()
	//out, err := GetByIndex[Many](oneTypeName)
	//_, _ = out, err
	return nil, nil
}

func (r ReferenceNavigationList[Many, One]) Get() ([]Many, error) {
	// make call using index

	return nil, nil
}
