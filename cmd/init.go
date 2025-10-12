package cmd

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const callbackBody = `// tevent format is following
/*
{
  timestamp: 1694567890123,
  fromAddress: "0x1234567890abcdef1234567890abcdef12345678",
  fromOwner: '0xSolanaFromOwner', // exists only for solana. spl token sender owner address 
  toAddress: "0xabcdef1234567890abcdef1234567890abcdef12",
  toOwner: '0xSolanaToOwner',  // exists only for solana. spl token receiver owner address 
  amount: "1000000000000000000", 
  tokenAddress: "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef", // erc20 or trc20 or spl-token mint address
  symbol: "USDT",
  chain: "ethereum",
  network: "mainnet",
  txHash: "0xa1b2c3d4e5f67890123456789abcdef0123456789abcdef0123456789abcdef",
  decimals: 18,
  position: 42 // abstract for block or slot 
}
*/

/*
Init your bot here
const TelegramBot = require('node-telegram-bot-api');
const botToken = '// your token here';
const bot = new TelegramBot(botToken, {polling: true});
*/

/*
 OR 
 const { Axiom } = require('@axiomhq/js');
 or use built in node js fetch api. 
 
 More built in standard libraries are coming soon, stay tuned ;)
*/

var stablecoins = ['USDC','USDT', 'DAI']

async function onTransfer(tevent) {
    var amt  = BigInt(tevent.amount)/BigInt(1e6) // ignore fractional
    if(stablecoins.includes(tevent.symbol) && amt > 1e6) {
        // take action on whale transfer across chains
        console.log("chain, txHash, amount")
        console.log(tevent.chain,"| ", tevent.txHash," | ", amt, " | ", tevent.symbol);
 		
		// OR 
		// send notification to your telegram bot 
		// bot.send .../

		// OR 
		// send slack message to your chat 
		// axiom ...


		// OR 
		// store it to your database via http 
		// axiom ...
	

		// OR 
		// call your own backend to handle processing  
		// axiom ...
 
    }
}

module.exports = onTransfer;`

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

		main := filepath.Join(folder, "main.js")
		if err := os.WriteFile(main, []byte(callbackBody), 0644); err != nil {
			fmt.Println("Error writing Dockerfile:", err)
			return
		}

		fmt.Println("Initialized project:", name)
	},
}
