package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net/http"
	"os"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the current huq as a Docker container",
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

		code, err := os.ReadFile("main.js")
		if err != nil {
			fmt.Println("Error: could not read huq.yaml:", err)
			return
		}

		token := "testtoken"

		var rRequest = RegisterRequest{
			Name:     cfg.Name,
			Version:  cfg.Version,
			Language: "js",
			Code:     string(code),
		}

		bb, _ := json.Marshal(rRequest)
		registerRequest, _ := http.NewRequest("POST", serverURL+"/register", bytes.NewReader(bb))
		registerRequest.Header.Set("Content-Type", "application/json")
		registerRequest.Header.Set("Authorization", token)

		resp, err := http.DefaultClient.Do(registerRequest)
		if err != nil {
			log.Fatalf("failed to register huq %s", err)
		}

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("failed to register huq %s", err)
		}

		log.Println(string(b))
		var callRequest = DeployRequest{
			UserID:  "testuserid",
			Name:    cfg.Name,
			Version: cfg.Version,
		}

		c, err := json.Marshal(callRequest)
		if err != nil {
			log.Fatalf("failed to call huq %s", err)
		}

		startRequest, _ := http.NewRequest("POST", serverURL+"/deploy", bytes.NewReader(c))
		startRequest.Header.Set("Content-Type", "application/json")
		startRequest.Header.Set("Authorization", token)

		resp2, err := http.DefaultClient.Do(startRequest)
		if err != nil {
			log.Fatalf("failed to call huq %s", err)
		}

		var deployResp DeployResponse
		err = json.NewDecoder(resp2.Body).Decode(&deployResp)
		if err != nil {
			log.Fatalf("failed to decode call huq  response %s", err)
		}

		log.Println(deployResp)

	},
}
