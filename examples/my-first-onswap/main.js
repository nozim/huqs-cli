// tevent format is following
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

var stablecoins = ['USDC', 'USDT', 'DAI']

async function onSwap(sevent) {
    if ((sevent.amount0Net == undefined) ||
        (sevent.amount1Net === undefined)) {
        return;
    }

    var amt0 = BigInt(sevent.amount0Net) / BigInt(1e18); // ignore fractional
    var amt1 = BigInt(sevent.amount1Net) / BigInt(1e18) // ignore fractional

    if ((stablecoins.includes(sevent.token0Symbol) && (amt0 > 10_000)) ||
        (stablecoins.includes(sevent.token1Symbol) && (amt1 > 10_000))) {
        console.log("whale trade alert!!!!");
        console.log(sevent);
    }
}

module.exports = onSwap;
