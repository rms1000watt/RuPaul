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

	switch r.Method {
	{{range $method := $path.Methods}}case http.{{GetHTTPMethod $method.Name}}:
		{{$path.Name | Title}}Handler{{$method.Name | ToUpper}}(w, r)
	{{end}}default:
		fmt.Println("Method not allowed:", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

	fmt.Println("Finished {{$path.Name | Title}}Handler!")
}
{{end}}

{{range $path := .API.Paths}}
{{range $method := $path.Methods}}
func {{$path.Name | Title}}Handler{{$method.Name | ToUpper}}(w http.ResponseWriter, r *http.Request) {
	// Assume JSON Serialization for now
	input := &{{$path.Name | Title}}Input{{$method.Name | ToUpper}}{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		fmt.Println("Failed decoding input:", err)
		http.Error(w, ErrorJSON("Input Error"), http.StatusInternalServerError)
		return
	}

	msg, err := Validate(input)
	if err != nil {
		fmt.Println("Failed validating input:", err)
		http.Error(w, ErrorJSON("Input Error"), http.StatusInternalServerError)
		return
	}

	if msg != "" {
		fmt.Println("Failed validation:", msg)
		http.Error(w, ErrorJSON("Invalid Input"), http.StatusBadRequest)
		return
	}	

	if err := Transform(input); err != nil {
		fmt.Println("Failed transforming input:", err)
		http.Error(w, ErrorJSON("Transform Error"), http.StatusInternalServerError)
		return
	}

	// Developer make updates here...

	output := get{{$path.Name | Title}}Output{{$method.Name | ToUpper}}(input)

	jsonBytes, err := json.Marshal(output)
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

	
}
{{end}}
{{end}}
{{end}}