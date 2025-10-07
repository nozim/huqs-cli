package cmd

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Play the current huq node as a Docker container",
	Run: func(cmd *cobra.Command, args []string) {
		// 1️⃣ Read huq.yaml
		data, err := os.ReadFile("huq.yaml")
		if err != nil {
			fmt.Println("Error: could not read huq.yaml:", err)
			return
		}

		var cfg Config
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			fmt.Println("Error: invalid huq.yaml:", err)
			return
		}

		if cfg.Name == "" || cfg.Version == "" {
			fmt.Println("Error: huq.yaml missing name or version")
			return
		}

		imageTag := fmt.Sprintf("%s:%s", cfg.Name, cfg.Version)
		fmt.Println("Calling Docker container:", imageTag)

		// 2️⃣ Prepare docker run command
		// -i: interactive to pass stdin
		// --rm: remove container after exit
		dockerCmd := exec.Command("docker", "run", "-i", "--rm", imageTag)

		// 3️⃣ Connect stdin/stdout
		dockerCmd.Stdin = os.Stdin
		dockerCmd.Stdout = os.Stdout
		dockerCmd.Stderr = os.Stderr

		// 4️⃣ Execute container
		if err := dockerCmd.Run(); err != nil {
			fmt.Println("Error running Docker container:", err)
			return
		}
	},
}
