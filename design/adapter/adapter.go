package adapter

type Adapter struct {
	adaptee *Adaptee
}

func NewAdapter() *Adapter {
	adapter := &Adapter{}
	adapter.adaptee = &Adaptee{}

	return adapter
}

func (ada *Adapter) Request() {
	ada.adaptee.SpecificRequest()
}
