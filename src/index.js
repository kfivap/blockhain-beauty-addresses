const { Worker } = require("worker_threads");
const fs = require('fs')

function main() {
    const numThreads = 4;
    const threads = [];
    const startAt = Date.now()
    let processed = 0
    let foundTotal = 0

    for (let i = 0; i < numThreads; i++) {

        const worker = new Worker(__dirname + "/worker.js", {
            workerData: {
                numberChars: 7,
                limit: 10 ** 10
            },
        });

        worker.on('error', err => {
            console.error(err);
        });

        worker.on("message", (msg) => {
            const { type, message } = msg
            if (type === 'log') {
                console.log(message)
            } else if (type === 'processed') {
                const msRun = Date.now() - startAt
                processed += message.processed
                console.log(new Date(), `=== running: ${((msRun) / 1000).toFixed()}s. processed ${processed} keys, found ${foundTotal} total addresses ===`)
            } else if (type === 'data') {
                foundTotal++
                fs.appendFileSync(__dirname + '/beauty_addresses.json', JSON.stringify(message) + '\n')
            } else {
                console.error('upsupported msg', msg)
            }
        });

        threads.push(worker);
    }

    threads.forEach(worker => {
        worker.on('exit', () => {
            console.log('Worker has exited.');
        });
    });
}
main()