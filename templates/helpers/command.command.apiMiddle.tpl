{{if .CommandLine.Command.API}}
addr := fmt.Sprintf("%s:%d", cfg.host, cfg.port)

fmt.Println("Starting server at:", addr)
log.Fatal(http.ListenAndServe(addr, ServerHandler()))
{{end}}