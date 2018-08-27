package magic

import (
	"net/http"
	"strings"
)

// Route structure
type Route struct {
	handlerPOST   func(*Context) error
	handlerGET    func(*Context) error
	handlerDELETE func(*Context) error
	handlerPUT    func(*Context) error
	middlewares   []Middleware
	branches      map[string]*Route
	param         *Route
	path          string
	fullPath      string
	isStatic      bool
}

// NewRoute function
// Generate new route
func NewRoute(path string) *Route {
	route := &Route{
		path: path,
	}
	route.branches = make(map[string]*Route)
	return route
}

// CreateRoute function
// Create new or get old route
// You can use middleware
func (route *Route) CreateRoute(path string, middlewares ...Middleware) *Route {
	router := route.createRoute(path)
	router.middlewares = middlewares
	return router
}

// GET function
// Add get handler to route
func (route *Route) GET(path string, handler func(context *Context) error) {
	route.add(path, "GET", handler)
}

// POST function
// Add post handler to route
func (route *Route) POST(path string, handler func(context *Context) error) {
	route.add(path, "POST", handler)
}

// PUT function
// Add put handler to route
func (route *Route) PUT(path string, handler func(context *Context) error) {
	route.add(path, "PUT", handler)
}

// DELETE function
// Add delete handler to route
func (route *Route) DELETE(path string, handler func(context *Context) error) {
	route.add(path, "DELETE", handler)
}

// FILE function
// Add get handler for file to route
func (route *Route) FILE(path, fileName string) {
	route.add(path, "GET", func(context *Context) error {
		http.ServeFile(context.Writer, context.Request, fileName)
		return nil
	})
}

// STATIC function
// Add static route
// Full path must't contain params like "/a/:id/static"
func (route *Route) STATIC(path, filePathName string) {
	if strings.Contains(route.fullPath+path, ":") {
		panic(errStaticRouteParams.Error() + ": " + route.fullPath + path)
	}
	route.add(path, "STATIC", func(context *Context) error {
		fileName := strings.SplitN(context.Request.URL.Path, route.fullPath+path, 2)[1]
		http.ServeFile(context.Writer, context.Request, filePathName+fileName)
		return nil
	})
}

func (route *Route) add(path, method string, handler func(*Context) error) {
	nowRoute := route
	branches := strings.Split(path, "/")
	len := len(branches)
	for i := 1; i < len; i++ {
		if nowRoute.isStatic {
			panic(errStaticRoute.Error() + ": " + route.path + path)
		}
		branch := branches[i]
		if branch != "" && branch[0] == ':' {
			nextRoute := nowRoute.param
			if nextRoute == nil {
				nextRoute = NewRoute(strings.Split(branch, ":")[1])
				nowRoute.param = nextRoute
			}
			nowRoute = nextRoute
		} else {
			nextRoute := nowRoute.branches[branch]
			if nextRoute == nil {
				nextRoute = NewRoute(branch)
				nowRoute.branches[branch] = nextRoute
			}
			nowRoute = nextRoute
		}
	}
	nowRoute.fullPath = route.fullPath + path
	nowRoute.isStatic = false
	switch method {
	case "POST":
		nowRoute.handlerPOST = handler
		break
	case "PUT":
		nowRoute.handlerPUT = handler
		break
	case "GET":
		nowRoute.handlerGET = handler
		break
	case "DELETE":
		nowRoute.handlerDELETE = handler
		break
	case "STATIC":
		nowRoute.handlerGET = handler
		nowRoute.isStatic = true
		break
	}
}

func (route *Route) createRoute(path string) *Route {
	nowRoute := route
	branches := strings.Split(path, "/")
	len := len(branches)
	for i := 1; i < len; i++ {
		if nowRoute.isStatic {
			panic(errStaticRoute)
		}
		branch := branches[i]
		if branch[0] == ':' {
			nextRoute := nowRoute.param
			if nextRoute == nil {
				nextRoute = NewRoute(strings.Split(branch, ":")[1])
				nowRoute.param = nextRoute
			}
			nowRoute = nextRoute
		} else {
			nextRoute := nowRoute.branches[branch]
			if nextRoute == nil {
				nextRoute = NewRoute(branch)
				nowRoute.branches[branch] = nextRoute
			}
			nowRoute = nextRoute
		}
	}
	nowRoute.fullPath = route.fullPath + path
	return nowRoute
}

func (route *Route) find(path, method string) (func(*Context) error, []Middleware, map[string]string) {
	nowRoute := route
	params := make(map[string]string)
	middlewares := []Middleware{}
	middlewares = append(middlewares, route.middlewares...)
	branches := strings.Split(path, "/")
	len := len(branches)
	for i := 1; i < len; i++ {
		if nowRoute.isStatic {
			break
		}
		branch := branches[i]
		nextRoute := nowRoute.branches[branch]
		if nextRoute == nil {
			nextRoute = nowRoute.param
			if nextRoute == nil {
				return nil, nil, nil
			}
			params[nextRoute.path] = branch
		}
		nowRoute = nextRoute
		middlewares = append(middlewares, nowRoute.middlewares...)
	}

	var result func(*Context) error
	switch method {
	case "POST":
		result = nowRoute.handlerPOST
		break
	case "PUT":
		result = nowRoute.handlerPUT
		break
	case "GET":
		result = nowRoute.handlerGET
		break
	case "DELETE":
		result = nowRoute.handlerDELETE
		break
	default:
		result = nil
	}
	return result, middlewares, params
}
