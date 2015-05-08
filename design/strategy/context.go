package strategy

type Context struct {
	strategy Strategy
}

func (context *Context) SetStrategy(strategy Strategy) {
	context.strategy = strategy
}

func (context *Context) Algorithm() {
	context.strategy.Algorithm()
}
