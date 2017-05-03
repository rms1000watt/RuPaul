{{if .CommandLine.Command.API}}
func ServerHandler() http.Handler {
	mux := http.NewServeMux()

	{{range $path := .API.Paths}}mux.HandleFunc("{{$path.Pattern}}", {{$path.Name | Title}}Handler)
	{{end}}

	return mux
}

{{range $path := .API.Paths}}
func {{$path.Name | Title}}Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting {{$path.Name | Title}}Handler...")

	// Assume JSON Serialization for now
	{{$path.Name}}Input := &{{$path.Name | Title}}Input{}
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

	if err := Transform({{$path.Name}}Input); err != nil {
		fmt.Println("Failed transforming input:", err)
		http.Error(w, ErrorJSON("Transform Error"), http.StatusInternalServerError)
		return
	}

	// Developer make updates here...

	{{$path.Name | ToLower}}Output := get{{$path.Name | Title}}Output({{$path.Name | ToLower}}Input)

	jsonBytes, err := json.Marshal({{$path.Name | ToLower}}Output)
	if err != nil {
		fmt.Println("Failed marshalling to JSON:", err)
		http.Error(w, ErrorJSON("JSON Marshal Error"), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(jsonBytes); err != nil {
		fmt.Println("Failed writing to response writer:", err)
		http.Error(w, ErrorJSON("Failed writing to output"), http.StatusInternalServerError)
		return
	}

	fmt.Println("Finished {{$path.Name | Title}}Handler!")
}
{{end}}
{{end}}