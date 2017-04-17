package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// {{.Command.Name}}Cmd represents the {{.Command.Name}} command
var {{.Command.Name}}Cmd = &cobra.Command{
	Use:   "{{.Command.Name}}",
	Short: "{{.Command.ShortDescription}}",
	Long: `{{.Command.LongDescription}}`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("{{.Command.Name}} called")
	},
}

func init() {
	RootCmd.AddCommand({{.Command.Name}}Cmd)
}
