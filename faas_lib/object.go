package faas_lib

var Separator = "::"

type Object interface {
	GetTypeName() string
}
