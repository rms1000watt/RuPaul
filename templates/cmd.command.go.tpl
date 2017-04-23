package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	
)

{{template "cmd.root.commandArgsTop.tpl" .}}

// {{.Name}}Cmd represents the {{.Name}} command
var {{.Name}}Cmd = &cobra.Command{
	Use:   "{{.Name}}",
	Short: "{{.ShortDescription}}",
	Long: `{{.LongDescription}}`,
	Run: Run{{.Name | Title}},
}

func Run{{.Name | Title}}(cmd *cobra.Command, args []string) {
	// Get config arguments and pass it to the function itself
	{{.Name}}Cfg := {{.Name}}.Config{
		{{template "cmd.root.commandArgsMiddle.tpl" .}}
	}

	{{.Name}}.{{.Name | Title}}({{.Name}}Cfg)

	fmt.Println("Remove me later...{{.Name}}")
}

func init() {
	RootCmd.AddCommand({{.Name}}Cmd)

	{{template "cmd.root.commandArgsBottom.tpl" .}}

	SetFlagsFromEnv({{.Name}}Cmd)
}
