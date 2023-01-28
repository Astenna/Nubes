package types

type Discount struct {
	Id            string
	Percentage    string
	isInitialized bool
}

func (Discount) GetTypeName() string {
	return "Discount"
}
func (receiver *Discount) Init() {
	receiver.isInitialized = true
}
