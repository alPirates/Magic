package magic

import (
	"strings"
)

// Route structure
type Route struct {
	handlerPOST   func(Context) (int, error)
	handlerGET    func(Context) (int, error)
	handlerDELETE func(Context) (int, error)
	handlerPUT    func(Context) (int, error)
	middlewares   []Middleware
	branches      map[string]*Route
	param         *Route
	path          string
}

// NewRoute function
func NewRoute(path string) *Route {
	route := &Route{
		path: path,
	}
	route.branches = make(map[string]*Route)
	return route
}

// CreateRoute function
func (route *Route) CreateRoute(path string, middlewares ...Middleware) *Route {
	router := route.createRoute(path)
	router.middlewares = middlewares
	return router
}

// GET function
func (route *Route) GET(path string, handler func(context Context) (int, error)) {
	route.add(path, "GET", handler)
}

// POST function
func (route *Route) POST(path string, handler func(context Context) (int, error)) {
	route.add(path, "POST", handler)
}

// PUT function
func (route *Route) PUT(path string, handler func(context Context) (int, error)) {
	route.add(path, "PUT", handler)
}

// DELETE function
func (route *Route) DELETE(path string, handler func(context Context) (int, error)) {
	route.add(path, "DELETE", handler)
}

func (route *Route) add(path, method string, handler func(Context) (int, error)) {
	nowRoute := route
	branches := strings.Split(path, "/")
	len := len(branches)
	for i := 1; i < len; i++ {
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
	}
}

func (route *Route) createRoute(path string) *Route {
	nowRoute := route
	branches := strings.Split(path, "/")
	len := len(branches)
	for i := 1; i < len; i++ {
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

	return nowRoute

}

func (route *Route) find(path, method string) (func(Context) (int, error), []Middleware, map[string]string) {
	nowRoute := route
	params := make(map[string]string)
	middlewares := []Middleware{}
	middlewares = append(middlewares, route.middlewares...)
	branches := strings.Split(path, "/")
	len := len(branches)
	for i := 1; i < len; i++ {
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

	var result func(Context) (int, error)
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
