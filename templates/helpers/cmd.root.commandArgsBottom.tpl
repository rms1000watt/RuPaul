{{range $arg := .Args}}{{$.Name}}Cmd.Flags().{{$arg.Type | Title}}Var(&{{$arg.Name | ToLower}}, "{{$arg.Name | ToLower}}", {{HandleQuotes $arg.Default $arg.Type}} ,"{{$arg.Description}}")
{{end}}