// swap event format is following
/*

{
  chain: 'binance',
  network: 'mainnet',
  blockNumber: Long { low: 66220409, high: 0, unsigned: true },
  blockHash: '0xca11090ba5db4a289e3ff7692ff3be4734d348714702fce869d59ae9a0a29134',
  txHash: '0x7938b5d64350c4aa1804c370407492a35aed32c52e00119a30e69d0a06650327',
  txIndex: 168,
  logIndex: 863,
  protocol: 'uniswap-v3',
  sender: '0x7cDa585e917FECB3a33C6d5A8F8A15dd956694dc',
  recipient: '0x8e76EBb1c71939982c9Ac267C0eb25F4aA739535',
  poolAddress: '0x47a90a2d92a8367a91efa1906bfc8c1e05bf10c4',
  token0Address: '0x55d398326f99059fF775485246999027B3197955',
  token0Symbol: 'USDT',
  token0Decimals: 18,
  token1Address: '0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c',
  token1Symbol: 'WBNB',
  token1Decimals: 18,
  amount0Net: '19012597680434120000000',
  amount1Net: '-16764862675733777699',
  sqrtPriceX96: '2352487631285775979655153994',
  liquidity: '2303367113602387546279419',
  tick: '-70341',
  position: Long { low: 863, high: 0, unsigned: true }
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
    if ((sevent.amount0Net === undefined) ||
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
