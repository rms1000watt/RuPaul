package cmd

import (
	"fmt"

	"io/ioutil"

	"github.com/rms1000watt/rupaul/generate"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	configFilePath = ""
	generateCmd    = &cobra.Command{
		Use:   "generate",
		Short: "Generates code from a `rupaul.yml` file",
		Long:  `Generates code from a "rupaul.yml" file`,
		Run:   runGenerate,
	}
)

func runGenerate(cmd *cobra.Command, args []string) {
	genCfg := generate.Config{}

	fileBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("Failed reading config:", err)
		return
	}

	if err := yaml.Unmarshal(fileBytes, &genCfg); err != nil {
		fmt.Println("Failed unmarshalling Yaml:", err)
		return
	}

	generate.Generate(genCfg)
}

func init() {
	RootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&configFilePath, "config-file", "f", "./rupaul.yml", "Config File Path of the RuPaul YAML")
}
