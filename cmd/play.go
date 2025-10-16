package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/coder/websocket"
	pb "github.com/nozim/huqs-cli/proto"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

//const serverURL = "http://huqs.heimdahl.xyz"

const serverURL = "http://localhost:9099"

type RegisterRequest struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Language string `json:"language"`
	Code     string `json:"code"`
}

type StartRequest struct {
	UserID  string `json:"user_id"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Payload string `json:"payload"`
}

type DeployRequest struct {
	UserID  string `json:"user_id"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Payload string `json:"payload"`
}

type StartResponse struct {
	InvocationID string `json:"invocation_id"`
	Name         string `json:"name"`
	Version      string `json:"version"`
	Status       string `json:"status"`
}

type DeployResponse struct {
	InvocationID string `json:"invocation_id"`
	Name         string `json:"name"`
	Version      string `json:"version"`
	Status       string `json:"status"`
}

type OutputMessage struct {
	HuqName      string
	HuqVersion   string
	InvocationID string
	Source       string
	Line         string
}

type StopRequest struct {
	UserID       string `json:"user_id"`
	InvocationID string `json:"invocation_id"`
}

type StopResponse struct {
	InvocationID string `json:"invocation_id"`
	Name         string `json:"name"`
	Version      string `json:"version"`
	Status       string `json:"status"`
}

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Play the current huq as a Docker container",
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
		//assert.Equal(t, "{\"imageTag\":\"nodesecond/transfer:0.0.1\",\"status\":\"huq registered\"}", string(b))
		var callRequest = StartRequest{
			UserID:  "testuserid",
			Name:    cfg.Name,
			Version: cfg.Version,
		}

		c, err := json.Marshal(callRequest)
		if err != nil {
			log.Fatalf("failed to call huq %s", err)
		}

		startRequest, _ := http.NewRequest("POST", serverURL+"/start", bytes.NewReader(c))
		startRequest.Header.Set("Content-Type", "application/json")
		startRequest.Header.Set("Authorization", token)

		//w := httptest.NewRecorder()
		resp2, err := http.DefaultClient.Do(startRequest)
		if err != nil {
			log.Fatalf("failed to call huq %s", err)
		}
		var startResp StartResponse

		err = json.NewDecoder(resp2.Body).Decode(&startResp)
		if err != nil {
			log.Fatalf("failed to decode call huq  response %s", err)
		}

		url := strings.ReplaceAll(serverURL, "http://", "ws://")
		conn, _, err := websocket.Dial(context.Background(), url+"/ws?auth_token="+token, nil)
		if err != nil {
			log.Fatalf("failed to setup connection to huq  %s", err)
		}

		// Create a channel to receive OS signals
		stop := make(chan os.Signal, 1)

		// Notify channel on interrupt (Ctrl+C) or SIGTERM
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		fmt.Println("Press Ctrl+C to stop")

		go func() {
			for {
				_, b, err := conn.Read(context.Background())
				if err != nil {
					log.Printf("failed to read ws message %s", err)
					return
				}
				var out pb.OutputMessage
				err = proto.Unmarshal(b, &out)
				if err != nil {
					log.Printf("failed to decode ws message %s", err)
					continue
				}

				fmt.Printf("%v\n", out.Line)

			}
		}()

		<-stop

		stopRequest := StopRequest{
			InvocationID: startResp.InvocationID,
			UserID:       token,
		}
		//
		stp, _ := json.Marshal(stopRequest)
		stopReq, _ := http.NewRequest("POST", serverURL+"/stop", bytes.NewReader(stp))
		startRequest.Header.Set("Content-Type", "application/json")
		startRequest.Header.Set("Authorization", token)

		resp4, err := http.DefaultClient.Do(stopReq)
		if err != nil {
			log.Fatalf("failed to invoke huq stop  %s", err)
		}
		var stopResp StopResponse

		err = json.NewDecoder(resp4.Body).Decode(&stopResp)
		if err != nil {
			log.Fatalf("failed to stop huq  %s", err)
		}

		log.Println(stopResp)
	},
}
