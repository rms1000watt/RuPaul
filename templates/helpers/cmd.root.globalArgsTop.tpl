var (
    {{range $globalArg := .CommandLine.GlobalArgs}}{{$globalArg.Name | ToLower}} {{$globalArg.Type | ToLower}}
    {{end}}
)