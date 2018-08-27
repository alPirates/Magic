package magic

import "net/http"

// Magic structure
type Magic struct {
	server *http.Server
	router *Router
}

// MaxBytes
var (
	MaxBytes = int64(10000000)
)

// NewMagic function
func NewMagic(port string) *Magic {
	var magic *Magic
	magic = &Magic{
		server: &http.Server{
			Addr: ":" + port,
		},
		router: NewRouter(),
	}
	magic.server.Handler = magic.router
	return magic
}

// CreateRoute function
func (magic *Magic) CreateRoute(path string, middlewares ...Middleware) *Route {
	router := magic.router.mainRoute.createRoute(path)
	router.middlewares = middlewares
	return router
}

// GET function
func (magic *Magic) GET(path string, handler func(context *Context) error) {
	magic.router.mainRoute.add(path, "GET", handler)
}

// POST function
func (magic *Magic) POST(path string, handler func(context *Context) error) {
	magic.router.mainRoute.add(path, "POST", handler)
}

// PUT function
func (magic *Magic) PUT(path string, handler func(context *Context) error) {
	magic.router.mainRoute.add(path, "PUT", handler)
}

// DELETE function
func (magic *Magic) DELETE(path string, handler func(context *Context) error) {
	magic.router.mainRoute.add(path, "DELETE", handler)
}

// SetMaxBytes function
func (magic *Magic) SetMaxBytes(maxBytes int64) {
	MaxBytes = maxBytes
}

// ListenAndServe function
func (magic *Magic) ListenAndServe() {
	magic.server.ListenAndServe()
}
