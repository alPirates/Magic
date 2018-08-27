package magic

import (
	"errors"
	"net/http"
	"strings"
)

// Magic structure
type Magic struct {
	server *http.Server
	router *Router
}

// MaxBytes
var (
	MaxBytes             = int64(10000000)
	errStaticRoute       = errors.New("can't add route to static route")
	errStaticRouteParams = errors.New("can't add route to path which contains params")
	magic                *Magic
)

// NewMagic function
func NewMagic(port string) *Magic {
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
// Create new or get old route
// You can use middleware
func (magic *Magic) CreateRoute(path string, middlewares ...Middleware) *Route {
	router := magic.router.mainRoute.createRoute(path)
	router.middlewares = middlewares
	return router
}

// GET function
// Add get handler to route
func (magic *Magic) GET(path string, handler func(context *Context) error) {
	magic.router.mainRoute.add(path, "GET", handler)
}

// POST function
// Add post handler to route
func (magic *Magic) POST(path string, handler func(context *Context) error) {
	magic.router.mainRoute.add(path, "POST", handler)
}

// PUT function
// Add put handler to route
func (magic *Magic) PUT(path string, handler func(context *Context) error) {
	magic.router.mainRoute.add(path, "PUT", handler)
}

// DELETE function
// Add delete handler to route
func (magic *Magic) DELETE(path string, handler func(context *Context) error) {
	magic.router.mainRoute.add(path, "DELETE", handler)
}

// FILE function
// Add get handler for file to route
func (magic *Magic) FILE(path, filename string) {
	magic.router.mainRoute.add(path, "GET", func(context *Context) error {
		http.ServeFile(context.Writer, context.Request, filename)
		return nil
	})
}

// STATIC function
// Add static route
// Full path must't contain params like "/a/:id/static"
func (magic *Magic) STATIC(path, filePathName string) {
	route := magic.router.mainRoute
	if strings.Contains(route.fullPath+path, ":") {
		panic(errStaticRouteParams.Error() + ": " + route.fullPath + path)
	}
	magic.router.mainRoute.add(path, "STATIC", func(context *Context) error {
		fileName := strings.SplitN(context.Request.URL.Path, route.fullPath+path, 2)[1]
		http.ServeFile(context.Writer, context.Request, filePathName+fileName)
		return nil
	})
}

// SetMaxBytes function
// set max bytes which you can upload
func (magic *Magic) SetMaxBytes(maxBytes int64) {
	MaxBytes = maxBytes
}

// ListenAndServe function
// Start server
func (magic *Magic) ListenAndServe() {
	magic.server.ListenAndServe()
}
