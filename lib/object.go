package lib

type Nobject interface {
	GetTypeName() string
}

type nobjectInit interface {
	Init()
}
