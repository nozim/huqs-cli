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

var transferVolumes = {}

const seenTransfers = new Map();

async function onTransfer(tevent) {
    if (tevent.symbol === 'unknown' || tevent.symbol === undefined) {
        return
    }

    const key = `${tevent.chain}-${tevent.txHash}-${tevent.tokenAddress}-${tevent.fromAddress}`;

    if (seenTransfers.has(key)) {
        // Duplicate, ignore
        return;
    }

    // Mark this key as seen
    seenTransfers.set(key, Date.now());

    var amt = BigInt(tevent.amount) / BigInt(1e6) // ignore fractional

    if (transferVolumes[tevent.chain] === undefined) {
        transferVolumes[tevent.chain] = new Map();
    }

    if (transferVolumes[tevent.chain][tevent.symbol] === undefined) {
        transferVolumes[tevent.chain][tevent.symbol] = {
            lastTx: tevent.txHash,
            transfers: 1,
            volume: BigInt(0),
        }
    }

    transferVolumes[tevent.chain][tevent.symbol].volume += amt;
    transferVolumes[tevent.chain][tevent.symbol].transfers += 1;
    transferVolumes[tevent.chain][tevent.symbol].lastTx = tevent.txHash;
}

const maxAge = 30 * 60 * 1000; // 30 minutes

setInterval(() => {
    const now = Date.now();
    // pruning old entries
    for (const [key, timestamp] of seenTransfers.entries()) {
        if (now - timestamp > maxAge) {
            seenTransfers.delete(key);
        }
    }

    console.log("-------------------------------------------")
    for (const chain in transferVolumes) {
        console.log(`Chain ${chain}:`);
        for (const token in transferVolumes[chain]) {

            console.log(`    ${token}:
            volume ${transferVolumes[chain][token].volume.toLocaleString()}\n
            transfers ${transferVolumes[chain][token].transfers}\n
            (last tx: ${transferVolumes[chain][token].lastTx})\n`);
        }
    }
    console.log("-------------------------------------------")
}, 10_000); // 10,000 ms = 10 seconds


module.exports = onTransfer;
