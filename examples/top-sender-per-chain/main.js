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

const TelegramBot = require('node-telegram-bot-api');
const botToken = '// add your token here';
const bot = new TelegramBot(botToken, {polling: true});

var subscribers = [];

// Anyone who messages the bot gets added
bot.on('message', (msg) => {
    const chatId = msg.chat.id;
    if (!subscribers.includes(chatId)) {
        subscribers.push(chatId);
        bot.sendMessage(chatId, 'Added! You will get updates every 1 min');
    }
});

const EXPLORER_URLS = {
    // 🪙 EVM-based
    ethereum: 'https://etherscan.io/tx/0x',
    arbitrum: 'https://arbiscan.io/tx/0x',
    polygon: 'https://polygonscan.com/tx/0x',
    binance: 'https://bscscan.com/tx/0x',
    base: 'https://basescan.org/tx/0x',

    // 🧱 Non-EVM chains
    solana: 'https://explorer.solana.com/tx/',
    bitcoin: 'https://mempool.space/tx/',
    tron: 'https://tronscan.org/#/transaction/',
};

function getExplorerLink(chain, tx) {
    if (!chain || !tx) return null;
    const txHash = typeof tx === 'string' ? tx : tx?.txHash;
    if (!txHash) return null;

    const base = EXPLORER_URLS[chain.toLowerCase()];
    return base ? `${base}${txHash}` : null;
}

const MAX_LENGTH = 4096;

function splitMessage(text) {
    const messages = [];
    for (let i = 0; i < text.length; i += MAX_LENGTH) {
        messages.push(text.slice(i, i + MAX_LENGTH));
    }
    return messages;
}


// leaderboard.js
class ChainTokenLeaderboard {
    constructor({topN = 1, maxEntries = 100} = {}) {
        this.data = new Map(); // chain -> symbol -> sender -> { total, lastTx }
        this.topN = topN;
        this.maxEntries = maxEntries;
    }

    /**
     * Escape MarkdownV2 special chars
     */
    escape(str = '') {
        return str.replace(/([_*\[\]()~`>#+\-=|{}.!])/g, '\\$1');
    }

    /**
     * Update leaderboard with a new transaction.
     * Keeps cumulative totals and stores last transaction details.
     */
    update(tx) {
        const {chain, symbol, fromAddress, amount} = tx;
        const amt = BigInt(amount);

        if (!this.data.has(chain)) this.data.set(chain, new Map());
        const chainMap = this.data.get(chain);

        if (!chainMap.has(symbol)) chainMap.set(symbol, new Map());
        const symbolMap = chainMap.get(symbol);

        const prev = symbolMap.get(fromAddress);
        const newTotal = (prev?.total ?? 0n) + amt;

        symbolMap.set(fromAddress, {
            total: newTotal,
            lastTx: tx, // store full last transaction object
        });

        if (symbolMap.size > this.maxEntries * 2) {
            this.prune(symbolMap);
        }
    }

    /**
     * Get top N senders for a specific chain + token
     */
    getTop(chain, symbol, n = this.topN) {
        const chainMap = this.data.get(chain);
        if (!chainMap) return [];

        const symbolMap = chainMap.get(symbol);
        if (!symbolMap) return [];

        return [...symbolMap.entries()]
            .sort((a, b) => (b[1].total > a[1].total ? 1 : -1))
            .slice(0, n)
            .map(([address, entry]) => ({
                address,
                total: entry.total,
                lastTx: entry.lastTx,
            }));
    }

    /**
     * Print top N senders per chain/token.
     */
    printAll(n = this.topN) {
        for (const [chain, chainMap] of this.data.entries()) {
            console.log(`\nChain: ${chain}`);
            for (const [symbol] of chainMap.entries()) {
                const top = this.getTop(chain, symbol, n);
                console.log(`  Token: ${symbol}`);

                if (top.length === 0) {
                    console.log('    (no senders yet)');
                    continue;
                }

                top.forEach(({address, total, lastTx}, i) => {
                    console.log(`    #${i + 1} ${address}: ${total.toString()} (last tx: ${getExplorerLink(chain, lastTx.txHash) ?? 'N/A'})`);
                });
            }
        }
    }

    printAllRich(n = this.topN) {
        const now = new Date().toISOString();
        const lines = [];

        for (const [chain, chainMap] of this.data.entries()) {
            lines.push(`${chain}`); // chain line

            for (const [symbol] of chainMap.entries()) {
                lines.push(`    ${symbol}`); // token line

                const top = this.getTop(chain, symbol, n);

                if (top.length === 0) {
                    lines.push(`        (no senders yet)`);
                    continue;
                }

                top.forEach(({address, total, lastTx}, i) => {
                    const rank = ['🥇', '🥈', '🥉'][i] || `#${i + 1}`;
                    const link = getExplorerLink(chain, lastTx?.txHash) ?? 'N/A';
                    const totalStr = total.toLocaleString();

                    lines.push(`        ${rank} Address: ${address}`);
                    lines.push(`           Total: ${totalStr}`);
                    lines.push(`           TX: ${link}`);
                    lines.push(`           Time: ${now}`);
                    lines.push(''); // blank line between senders
                });
            }
        }

        return lines.join('\n') || 'No data yet';
    }

    /**
     * Keep only top maxEntries for a symbol
     */
    prune(symbolMap) {
        const topEntries = [...symbolMap.entries()]
            .sort((a, b) => (b[1].total > a[1].total ? 1 : -1))
            .slice(0, this.maxEntries);

        symbolMap.clear();
        for (const [addr, entry] of topEntries) {
            symbolMap.set(addr, entry);
        }
    }

    /**
     * Return all chain-symbol pairs
     */
    getAllKeys() {
        const result = [];
        for (const [chain, chainMap] of this.data.entries()) {
            for (const symbol of chainMap.keys()) {
                result.push({chain, symbol});
            }
        }
        return result;
    }
}

const lb = new ChainTokenLeaderboard({topN: 3});

const seenTransfers = new Set();


async function onTransfer(tevent) {
    if (tevent.symbol === 'unknown' || tevent.symbol === undefined) {
        return;
    }

    const key = `${tevent.chain}-${tevent.txHash}-${tevent.tokenAddress}-${tevent.fromAddress}`;

    if (seenTransfers.has(key)) {
        // Duplicate, ignore
        return;
    }

    // Mark this key as seen
    seenTransfers.add(key);

    lb.update(tevent);
}

setInterval(() => {
    const msg = lb.printAllRich();
    const chunks = splitMessage(msg);

    subscribers.forEach(chatId => {
        // Usage
        for (const chunk of chunks) {
            bot.sendMessage(chatId, chunk, {parse_mode: 'HTML'});
        }
    });
}, 5 * 60 * 1_000); // every 5 mins

module.exports = onTransfer;
