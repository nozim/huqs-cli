package cmd

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

//const entryPointBody = `const onTransfer = require("./main");
//
//async function readStdin() {
//    return new Promise((resolve, reject) => {
//        let data = "";
//        process.stdin.setEncoding("utf8");
//
//        process.stdin.on("data", chunk => {
//            data += chunk;
//        });
//
//        process.stdin.on("end", () => {
//            resolve(data);
//        });
//
//        process.stdin.on("error", reject);
//    });
//}
//
//(async () => {
//    try {
//        const input = await readStdin();
//        const parsed = JSON.parse(input); // expecting { event: ..., context: ... }
//
//        const result = onTransfer(parsed, {});
//
//        if (result instanceof Promise) {
//            const awaited = await result;
//            console.log(JSON.stringify(awaited, null, 2));
//        } else {
//            console.log(JSON.stringify(result, null, 2));
//        }
//    } catch (err) {
//        console.error("Error:", err.message);
//        process.exit(1);
//    }
//})();`

const callbackBody = `// tevent format is following
/*
{
  timestamp: Date.now(), // e.g. 1694567890123
  fromAddress: "0x1234567890abcdef1234567890abcdef12345678",
  fromOwner: { value: "Alice" }, // optional
  toAddress: "0xabcdef1234567890abcdef1234567890abcdef12",
  toOwner: { value: "Bob" },     // optional
  amount: "1000000000000000000",  // string for big.Int
  tokenAddress: "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
  symbol: "USDT",
  chain: "Ethereum",
  network: "mainnet",
  txHash: "0xa1b2c3d4e5f67890123456789abcdef0123456789abcdef0123456789abcdef",
  decimals: 18,
  position: 42n
}
*/

var stablecoins = ['USDC','USDT', 'DAI']

async function onTransfer(tevent) {
    var amt  = BigInt(tevent.amount)/BigInt(1e6) // ignore fractional
    if(stablecoins.includes(tevent.symbol) && amt > 1e6) {
        // take action on whale transfer across chains
        console.log("chain, txHash, amount")
        console.log(tevent.chain,"| ", tevent.txHash," | ", amt, " | ", tevent.symbol);
    }
}

module.exports = onTransfer;`

//
//const dockerBody = `# Use the latest Node.js LTS image from Docker Hub
//FROM node:latest
//
//# Set working directory inside the container
//WORKDIR /usr/src/app
//
//# Copy package.json if you have dependencies (optional)
//# COPY package*.json ./
//# RUN npm install --production
//
//# Copy all project files into the container
//COPY . .
//
//# Run your main.js file with Node.js
//CMD ["node", "entrypoint.js"]`

type Config struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a new Huq node",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		folder := strings.ReplaceAll(name, "/", "-")

		// Create folder
		if err := os.MkdirAll(folder, os.ModePerm); err != nil {
			fmt.Println("Error creating folder:", err)
			return
		}

		// Create YAML
		config := Config{
			Name:    name,
			Version: "0.0.1",
		}
		yamlData, _ := yaml.Marshal(&config)
		yamlPath := filepath.Join(folder, "huq.yaml")
		if err := os.WriteFile(yamlPath, yamlData, 0644); err != nil {
			fmt.Println("Error writing YAML:", err)
			return
		}

		// Create empty Dockerfile
		//dockerfilePath := filepath.Join(folder, "Dockerfile")
		//if err := os.WriteFile(dockerfilePath, []byte(dockerBody), 0644); err != nil {
		//	fmt.Println("Error writing Dockerfile:", err)
		//	return
		//}

		//entrypointPath := filepath.Join(folder, "entrypoint.js")
		//if err := os.WriteFile(entrypointPath, []byte(entryPointBody), 0644); err != nil {
		//	fmt.Println("Error writing Dockerfile:", err)
		//	return
		//}

		main := filepath.Join(folder, "main.js")
		if err := os.WriteFile(main, []byte(callbackBody), 0644); err != nil {
			fmt.Println("Error writing Dockerfile:", err)
			return
		}

		fmt.Println("Initialized project:", name)
	},
}
