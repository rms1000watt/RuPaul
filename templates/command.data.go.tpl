package {{.CommandLine.Command.Name}}

{{range $path := .API.Paths}}
type {{$path.Name | Title}}Input struct {
    {{range $input := $path.Inputs}}{{$input.Name}} {{$input.Type}} `json:"{{$input.Name | ToSnakeCase}},omitempty" transform:"{{GenTransformStr $input}}" validate:"{{GenValidationStr $input}}"`
    {{end}}
}

type {{$path.Name | Title}}Output struct {
    {{range $output := $path.Outputs}}{{$output.Name}} {{$output.Type}} `json:"{{$output.Name | ToSnakeCase}},omitempty"`
    {{end}}
}
{{end}}
