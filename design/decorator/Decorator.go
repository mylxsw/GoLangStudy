package decorator

type Decorator struct {
	component Component
}

func (dec *Decorator) setComponent(component Component) {
	dec.component = component
}

func (dec *Decorator) Operation() {
	dec.component.Operation()
}
