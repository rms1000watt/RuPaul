package {{.CommandLine.Command.Name}}

import (
    "fmt"
)

func {{.CommandLine.Command.Name | Title}}(cfg Config) {
    fmt.Println("Config:", cfg)

    fmt.Println("{{.CommandLine.Command.Name}} called..")
}
