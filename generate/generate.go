package generate

import (
	"fmt"
	"os"
	"path/filepath"
)

func Generate(cfg Config) {
	fmt.Println("Config:", cfg)

	if err := os.MkdirAll(filepath.Join("out", "cmd"), os.ModePerm); err != nil {
		fmt.Println("Failed mkdir out/cmd:", err)
		return
	}
}
