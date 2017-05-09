package {{.CommandLine.Command.Name}}

{{if .API.Middlewares.CORS}}{{template "middleware-cors-vars.tpl" .}}{{end}}

{{if .API.Middlewares}}
// Usage: HandleMiddlewares(PersonHandlerGET, MiddlewareNoCache, MiddlewareCORS)(w, r)
func HandleMiddlewares(handlerFunc http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) (h http.HandlerFunc) {
	for _, mw := range middlewares {
		handlerFunc = mw(handlerFunc)
	}
	return handlerFunc
}
{{end}}

{{if .API.Middlewares}}{{template "middleware-no-cache.tpl" .}}{{end}}
{{if .API.Middlewares}}{{template "middleware-cors-func.tpl" .}}{{end}}
{{if .API.Middlewares}}{{template "middleware-logging.tpl" .}}{{end}}
