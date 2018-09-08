package magic

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
)

var (
	wait = &sync.WaitGroup{}
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
	errMethod            = errors.New("invalid method")
	magic                *Magic
)

// NewMagic function
func NewMagic(port string) *Magic {
	wait.Add(1)
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
func (magic *Magic) FILE(path, fileName string) {
	magic.router.mainRoute.FILE(path, fileName)
}

// STATIC function
// Add static route
// Full path must't contain params like "/a/:id/static"
func (magic *Magic) STATIC(path, filePathName string) {
	magic.router.mainRoute.STATIC(path, filePathName)
}

// CUSTOM function
// Can be STATIC or FILE
// In static you can get filename in storage | context.Storage["fileName"]
func (magic *Magic) CUSTOM(path, method string, handler func(context *Context) error) {
	magic.router.mainRoute.CUSTOM(path, method, handler)
}

// SetMaxBytes function
// set max bytes which you can upload
func (magic *Magic) SetMaxBytes(maxBytes int64) {
	MaxBytes = maxBytes
}

// Close function
// Close all routes (server)
func (magic *Magic) Close() {
	magic.server.Close()
}

// Restart function
// Restart magic (server)
func (magic *Magic) Restart() {
	wait.Add(1)
	magic.server.Close()
	go magic.start(wait)
}

// ListenAndServe function
// Start server
func (magic *Magic) ListenAndServe() {
	fmt.Println("")
	fmt.Println(
		` -----------------------------------------------------`)
	fmt.Println("")
	fmt.Println(" "+
		`    __   __     __     _______    ________   _______`+"\n",
		`   /  | /  |   /  \   |   ____|  |__    __| |   ____|`+"\n",
		`  /   |/   |  / /\ \  |  | ____     |  |    |  |`+"\n",
		` /   __    | |  \/  | |  |__|  |  __|  |__  |  |____`+"\n",
		`/___/  |___| |__/\__| |________| |________| |_______|`+"\n",
	)
	fmt.Println("")
	fmt.Println(
		` -----------------------------------------------------`)
	fmt.Println("")
	go magic.start(wait)
	wait.Wait()
}

func (magic *Magic) start(wait *sync.WaitGroup) {
	magic.server.ListenAndServe()
	wait.Done()
}
