package magic

import "net/http"

// Context structure
type Context struct {
	Writer      http.ResponseWriter
	Request     *http.Request
	Params      map[string]string
	QueryParams map[string]string
}

func getContext(writer http.ResponseWriter, request *http.Request, params map[string]string) Context {
	return Context{
		Writer:  writer,
		Request: request,
		Params:  params,
	}
}
