package {{.CommandLine.Command.Name}}

import (
    "fmt"
)

func {{.CommandLine.Command.Name | Title}}(cfg Config) {
    fmt.Println("Config:", cfg)

    {{template "command.command.apiMiddle.tpl" .}}

    fmt.Println("{{.CommandLine.Command.Name}} called..")
}

{{template "command.command.apiBottom.tpl" .}}
