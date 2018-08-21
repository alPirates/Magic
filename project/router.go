package magic

import (
	"net/http"
)

// Router structure
type Router struct {
	mainRoute *Route
}

// NewRouter function
func NewRouter() *Router {
	var router *Router
	router = &Router{
		mainRoute: NewRoute(""),
	}
	return router
}

// Handle interface
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, middlewares, params := router.mainRoute.find(r.URL.Path, r.Method)
	// queryParams := r.

	context := getContext(w, r, params)

	if handler != nil {
		startHandler(context, middlewares, handler)
	} else {
		http.NotFound(w, r)
	}
}

func startHandler(context Context, middlewares []Middleware, handler func(context Context) (int, error)) {

	for _, middleware := range middlewares {
		code, err := middleware.run(context)
		if err != nil {
			http.Error(context.Writer, err.Error(), code)
		}
	}

	code, err := handler(context)
	if err != nil {
		http.Error(context.Writer, err.Error(), code)
	}

}
