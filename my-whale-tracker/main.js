// tevent format is following
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

var stablecoins = ['USDC', 'USDT', 'DAI']

var transferVolumes = {}

async function onTransfer(tevent) {
    var amt = BigInt(tevent.amount) / BigInt(1e6) // ignore fractional
    if (transferVolumes[tevent.chain + " " + tevent.tokenAddress] === undefined) {
        transferVolumes[tevent.chain + " " + tevent.tokenAddress] = {
            lastTx: tevent.txHash,
            volume: BigInt(0),
        }
    }

    transferVolumes[tevent.chain + " " + tevent.tokenAddress].volume += amt;
    transferVolumes[tevent.chain + " " + tevent.tokenAddress].lastTx = tevent.txHash;
}


setInterval(() => {
    console.log("-------------------------------------------")
    for (const key in transferVolumes) {
        console.log(`${key}: ${transferVolumes[key].volume} (last tx: ${transferVolumes[key].lastTx})\n`);
    }
    console.log("-------------------------------------------")
}, 10_000); // 10,000 ms = 10 seconds


module.exports = onTransfer;