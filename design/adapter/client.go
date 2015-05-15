package adapter

func ClientFunc(target Target) {
	target.Request()
}
