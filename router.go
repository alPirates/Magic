package magic

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
)

// Router structure
type Router struct {
	mainRoute *Route
}

// NewRouter function
// Generate new Router
func NewRouter() *Router {
	var router *Router
	router = &Router{
		mainRoute: &Route{
			path:     "",
			fullPath: "",
		},
	}
	router.mainRoute.branches = make(map[string]*Route)
	return router
}

// Handle interface
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, middlewares, params := router.mainRoute.find(r.URL.Path, r.Method)
	context := getContext(w, r)
	if handler != nil {

		context.Params = params
		queryParams, err := url.ParseQuery(r.URL.RawQuery)
		if err == nil {
			context.QueryParams = map[string][]string(queryParams)
		}

		err = r.ParseForm()
		if err == nil {
			postParams := map[string][]string(r.PostForm)
			context.PostParams = postParams
		}

		err = r.ParseMultipartForm(MaxBytes)
		if err == nil {
			multipartParams := map[string][]string(r.MultipartForm.Value)
			context.MultipartParams = multipartParams
			files := map[string][]*multipart.FileHeader(r.MultipartForm.File)
			context.FileParams = files
		}

		headers := map[string][]string(r.Header)

		context.Headers = headers

		bytes, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		context.Body = string(bytes)

		startHandler(context, middlewares, handler)
	} else {
		context.SendErrorString("page not found")
	}
}

func startHandler(context *Context, middlewares []Middleware, handler func(context *Context) error) {

	for _, middleware := range middlewares {
		err := middleware.run(context)
		if err != nil {
			return
		}
	}

	err := handler(context)
	if err != nil {
		return
	}

}
