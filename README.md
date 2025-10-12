# Huqs CLI

Build real-time blockchain applications with JavaScript.

## What is Huqs?

Huqs lets you write JavaScript functions that react to blockchain events across all chains - Ethereum, Arbitrum, Base,
Polygon, Tron, Solana, Bitcoin.

No nodes. No infrastructure. Just code.

## Installation

Download binary Get the latest release from GitHub releases

or build from repo:

```bash
git clone https://github.com/nozim/huqs-cli.git
go build -o huq .
sudo cp ./huq /usr/local/bin
```

## Quick Start

#### 1. Create a new huq:

```bash
huq init whale-tracker
cd whale-tracker
```

#### 2. Edit main.js:

```javascript
var stablecoins = ['USDC', 'USDT', 'DAI']

async function onTransfer(tevent) {
    var amt = BigInt(tevent.amount) / BigInt(1e6) // ignore fractional
    if (stablecoins.includes(tevent.symbol) && amt > 1e6) {
        // take action on whale transfer across chains
        console.log("chain, txHash, amount")
        console.log(tevent.chain, "| ", tevent.txHash, " | ", amt, " | ", tevent.symbol);
    }
}

module.exports = onTransfer;
```

#### 4. Test it locally:

```bash
 huq play   
2025/10/12 20:55:17 {"imageTag":"simple-w:0.0.1","status":"huq registered"}                                                                                             
Press Ctrl+C to stop                      
STARTING                                                                                                                                                                
RUNNING                                   
chain, txHash, amount                                                                                                                                                   

arbitrum |  0x432b32ab44731783b749202fa93fd1d437f44412f002eeb13600f41eed736699  |  164966n  |  USDC
                                                                                                                                                                        
chain, txHash, amount
arbitrum |  0x432b32ab44731783b749202fa93fd1d437f44412f002eeb13600f41eed736699  |  164966n  |  USDC
                                                                                                                                                                        
chain, txHash, amount
                                          
arbitrum |  0x432b32ab44731783b749202fa93fd1d437f44412f002eeb13600f41eed736699  |  119975n  |  USDC
                                          
chain, txHash, amount
arbitrum |  0x432b32ab44731783b749202fa93fd1d437f44412f002eeb13600f41eed736699  |  119975n  |  USDC
```

#### 5. Stop it

```bash 
....
chain, txHash, amount
tron |  e5226966fb57f3750999e5716c75834b1d3ffa8b0fdf6eac9dc78d0db9c679e6  |  200000n  |  USDT

chain, txHash, amount
tron |  e5226966fb57f3750999e5716c75834b1d3ffa8b0fdf6eac9dc78d0db9c679e6  |  200000n  |  USDT

^CSTOPPING
EXITED
```

#### 5. Deploy it:

```bash
 huq deploy 
2025/10/12 20:59:21 {"imageTag":"simple-w:0.0.1","status":"huq registered"}
2025/10/12 20:59:21 {fb8fe8df-e874-4677-94e3-d164101151c9 simple-w 0.0.1 RUNNING}
````

#### 6. Done.

Your huq runs forever, processing blockchain events in real-time.

### Check Examples folder for more:

#### Token volume tracker by chain

Aggregates volume by chain and each address and outputs every 10 seconds.

#### Top sender (interval)

Aggregates volume by chain, token and sender and picks top 10 and sends report to telegram bot subscribers.

#### Monitor DEX activity (Coming soon)

```javascript
async function onSwap(swap) {
    if (swap.tokenIn > 50000) {
        console.log(`${swap.dex}: ${swap.tokenIn} → ${swap.tokenOut}`);
    }
}

module.exports = onSwap;
````

### Build a Telegram bot

```javascript
const TelegramBot = require('node-telegram-bot-api');
const bot = new TelegramBot(/*token*/);

async function onTransfer(transfer) {
    if (transfer.value_usd > 1000000) {
        bot.sendMessage(CHAT_ID, `Whale: $${transfer.amount}`);
    }
}

```

### Coming

#### Fully Uniswap Swap v2 and V3 compatible events live

```javascript
/*
{
  position: 5645,
  txHash: 
  chain: 'ethereum',
  dex: 'uniswap-v3',
  tokenIn: 'USDC',
  tokenOut: 'WETH',
  amountIn: '1000000',
  amountOut: '500',
  value_usd: 1000000,
  txHash: '0x...',
  timestamp: 1234567890
}
*/

async function onSwap(swap) {
    // handle swap here 
}

module.exports = onSwap;
```

#### Token suppply events

```javascript

async function onBurn(info) { // and onMint

}

module.exports = onBurn;
```

#### Block production event

```javascript

async function onBlock(info) {

}

module.exports = onBlock;
```

#### Tx production event

```javascript

async function onTransaction(info) { // and onMint

}

module.exports = onTransaction;
```

and many more.