{{if .CommandLine.Command.API}}
func ServerHandler() http.Handler {
	mux := http.NewServeMux()

	{{range $path := .API.Paths}}mux.HandleFunc("{{$path.Pattern}}", {{$path.Name | Title}}Handler)
	{{end}}

	return mux
}

{{range $path := .API.Paths}}func {{$path.Name | Title}}Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting {{$path.Name | Title}}Handler...")

	// Assume JSON Serialization for now
	{{$path.Name}}Input := {{$path.Name | Title}}Input{}
	if err := json.NewDecoder(r.Body).Decode(&{{$path.Name}}Input); err != nil {
		fmt.Println("Failed decoding input:", err)
		http.Error(w, "Input Error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Finished {{$path.Name | Title}}Handler!")
}
{{end}}
{{end}}