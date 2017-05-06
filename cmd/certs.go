package cmd

import (
	"github.com/rms1000watt/rupaul/generate"
	"github.com/spf13/cobra"
)

var (
	certsPath   = ""
	letsEncrypt = false
)

var certsCmd = &cobra.Command{
	Use:   "certs",
	Short: "Generate Certs for HTTPS communication",
	Long: `Generate Certs:
1. Self signed
2. Let's Encrypt (TODO)`,
	Run: runCerts,
}

func init() {
	generateCmd.AddCommand(certsCmd)
	certsCmd.Flags().StringVar(&certsPath, "certs-path", "./certs", "Certs path that contains openssl.cnf")
	certsCmd.Flags().BoolVar(&letsEncrypt, "lets-encrypt", false, "Generate lets-encrypt certs")
}

func runCerts(cmd *cobra.Command, args []string) {
	generate.Certs(certsPath, letsEncrypt)
}
