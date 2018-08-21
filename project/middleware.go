package magic

// Middleware structure
type Middleware struct {
	run func(Context) (int, error)
}

// NewMiddleware function
func NewMiddleware(handler func(context Context) (int, error)) Middleware {
	return Middleware{
		run: handler,
	}
}
