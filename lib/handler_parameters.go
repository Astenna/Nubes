package lib

type HandlerParameters struct {
	// Indicates the ID of assocaiated object instance
	Id string
	// Parameter of the orginal function
	// from which the handler is generated
	Parameter interface{}
}
