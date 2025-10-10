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

var senders = {}
var grouped = {};

var chains = {};

async function onTransfer(tevent) {
    if (!stablecoins.includes(tevent.symbol)) {
        return;
    }
    var key = tevent.chain + " " + tevent.fromAddress + " " + tevent.symbol;

    var amt = BigInt(tevent.amount) / BigInt(1e6) // ignore fractional
    if (senders[key] === undefined) {
        senders[key] = {
            lastTx: tevent.txHash,
            volume: BigInt(0),
        }
    }

    senders[key].volume += amt;
    senders[key].lastTx = tevent.txHash;

}

setInterval(() => {
    const sorted = Object.entries(senders)
        .sort((a, b) => {
            if (a[1].volume > b[1].volume) return -1;
            if (a[1].volume < b[1].volume) return 1;
            return 0;
        })
        .slice(0, 10);

    // Clear console before printing new data
    console.clear();

    const now = new Date().toISOString();
    console.log(`\n=== Top 10 senders by volume @ ${now} ===\n`);

    sorted.forEach(([key, info], index) => {
        console.log(
            `${(index + 1).toString().padStart(2, " ")}. ${key} → ${info.volume.toString()} (lastTx: ${info.lastTx})`
        );
    });

    if (sorted.length === 0) {
        console.log("No senders yet.");
    }

    console.log("\n----------------------------------------\n");
}, 20000);

module.exports = onTransfer;