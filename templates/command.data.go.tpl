package {{.CommandLine.Command.Name}}

{{range $path := .API.Paths}}
type {{$path.Name | Title}}Input struct {
    {{range $input := $path.Inputs}}{{$input.Name}} *{{$input.Type}} `json:"{{$input.Name | ToSnakeCase}},omitempty" transform:"{{GenTransformStr $input}}" validate:"{{GenValidationStr $input}}"`
    {{end}}
}

type {{$path.Name | Title}}Output struct {
    {{range $output := $path.Outputs}}{{$output.Name}} {{$output.Type}} `json:"{{$output.Name | ToSnakeCase}},omitempty"`
    {{end}}
}
{{end}}

{{range $path := .API.Paths}}func get{{$path.Name | Title}}Output({{$path.Name | ToLower}}Input *{{$path.Name | Title}}Input) ({{$path.Name | ToLower}}Output {{$path.Name | Title}}Output) {
	if {{$path.Name | ToLower}}Input == nil {
		return
	}
	
	{{range $output := $path.Outputs}}{{if OutputInInputs $output.Name $path.Inputs}}{{$output.Name | ToCamelCase}} := {{EmptyValue $output.Type}}
	if {{$path.Name | ToLower}}Input.{{$output.Name | Title}} != nil {
		{{$output.Name | ToCamelCase}} = *{{$path.Name | ToLower}}Input.{{$output.Name | Title}}
	}{{end}}
	
	{{end}}

	{{$path.Name | ToLower}}Output = {{$path.Name | Title}}Output{
		{{range $output := $path.Outputs}}{{if OutputInInputs $output.Name $path.Inputs}}{{$output.Name | Title}}: {{$output.Name | ToCamelCase}},
		{{end}}{{end}}
	}
	return
}
{{end}}
