// TODO: Put gnarley logic in here for parsing based on type
var (
    {{range $globalArg := .CommandLine.GlobalArgs}}{{$globalArg.Name}} = {{$globalArg.Default}}
    {{end}})