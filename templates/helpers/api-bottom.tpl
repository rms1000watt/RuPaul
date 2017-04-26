{{if .CommandLine.Command.API}}
func ServerHandler() http.Handler {
	mux := http.NewServeMux()

	{{range $path := .API.Paths}}mux.HandleFunc("{{$path.Pattern}}", {{$path.Name | Title}}Handler)
	{{end}}

	return mux
}

{{range $path := .API.Paths}}// curl -X POST -d '{"first_name":"Chet","middle_name":"Darf","last_name":"Star"}' localhost:8080/{{$path.Name}}
func {{$path.Name | Title}}Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting {{$path.Name | Title}}Handler...")

	// Assume JSON Serialization for now
	{{$path.Name}}Input := {{$path.Name | Title}}Input{}
	if err := json.NewDecoder(r.Body).Decode(&{{$path.Name}}Input); err != nil {
		fmt.Println("Failed decoding input:", err)
		http.Error(w, ErrorJSON("Input Error"), http.StatusInternalServerError)
		return
	}

	ok, msg, err := Validate({{$path.Name}}Input)
	if err != nil {
		fmt.Println("Failed validating input:", err)
		http.Error(w, ErrorJSON("Input Error"), http.StatusInternalServerError)
		return
	}

	if !ok {
		fmt.Println("Failed validation:", msg)
		http.Error(w, ErrorJSON("Invalid Input"), http.StatusBadRequest)
		return
	}	

	// TODO: Transform function

	// Developer make updates here...

	// TODO: JSON Marshal output

	fmt.Println("Finished {{$path.Name | Title}}Handler!")
}
{{end}}
{{end}}