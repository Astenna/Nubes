package lib

type HandlerParameters struct {
	// Indicates the ID of associated object instance
	Id string
	// Indicates the ID of associated object instance
	TypeName string
	// Parameter of the orginal function
	// from which the handler is generated
	Parameter interface{}
}
