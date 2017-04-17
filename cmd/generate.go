package cmd

import (
	"fmt"

	"io/ioutil"

	"github.com/rms1000watt/rygen/generate"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	cfg         = &Config{}
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generates code from a `rygen.yml` file",
		Long:  `Generates code from a "rygen.yml" file`,
		Run: func(cmd *cobra.Command, args []string) {
			genCfg := generate.Config{}

			fileBytes, err := ioutil.ReadFile(cfg.ConfigFilePath)
			if err != nil {
				fmt.Println("Failed reading config:", err)
				return
			}

			if err := yaml.Unmarshal(fileBytes, &genCfg); err != nil {
				fmt.Println("Failed unmarshalling Yaml:", err)
				return
			}

			generate.Generate(genCfg)
		},
	}
)

func init() {
	RootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&cfg.ConfigFilePath, "config-file", "f", "./rygen.yml", "Config File Path of the RyGen YAML")
}
