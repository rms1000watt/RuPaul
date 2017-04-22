var (
    {{range $arg := .Args}}{{$arg.Name | ToLower}} {{$arg.Type | ToLower}}
    {{end}}
)