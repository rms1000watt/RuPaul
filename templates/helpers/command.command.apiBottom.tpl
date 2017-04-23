{{if .CommandLine.Command.API}}
func ServerHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/pizza", PizzaHandler)
	mux.HandleFunc("/", RootHandler)

	return mux
}

func PizzaHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pizza"))
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
{{end}}