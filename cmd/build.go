package cmd

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build Docker image",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := filepath.Join(".", "huq.yaml")
		data, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Println("Could not read huq.yaml:", err)
			return
		}

		var config Config
		if err := yaml.Unmarshal(data, &config); err != nil {
			fmt.Println("Error parsing YAML:", err)
			return
		}

		imageTag := fmt.Sprintf("%s:%s", config.Name, config.Version)
		fmt.Println("Building Docker image:", imageTag)

		cmdBuild := exec.Command("docker", "build", "-t", imageTag, ".")
		cmdBuild.Stdout = os.Stdout
		cmdBuild.Stderr = os.Stderr

		if err := cmdBuild.Run(); err != nil {
			fmt.Println("Docker build failed:", err)
			return
		}

		fmt.Println("Docker build complete!")
	},
}
