const onTransfer = require("./main");

async function readStdin() {
    return new Promise((resolve, reject) => {
        let data = "";
        process.stdin.setEncoding("utf8");

        process.stdin.on("data", chunk => {
            data += chunk;
        });

        process.stdin.on("end", () => {
            resolve(data);
        });

        process.stdin.on("error", reject);
    });
}

(async () => {
    try {
        const input = await readStdin();
        const parsed = JSON.parse(input); // expecting { event: ..., context: ... }

        const result = onTransfer(parsed, {});

        if (result instanceof Promise) {
            const awaited = await result;
            console.log(JSON.stringify(awaited, null, 2));
        } else {
            console.log(JSON.stringify(result, null, 2));
        }
    } catch (err) {
        console.error("Error:", err.message);
        process.exit(1);
    }
})();