const { Worker } = require("worker_threads");
const fs = require('fs')

function main() {
  const numThreads = 4;
  const threads = [];

  for (let i = 0; i < numThreads; i++) {

    const worker = new Worker(__dirname + "/worker.js", {
      workerData: {
        numberChars: 5,
        limit: 1000000
      },
    });

    worker.on('error', err => {
      console.error(err);
    });

    worker.on("message", (msg) => {
      const { type, message } = msg
      if (type === 'log') {
        console.log(message)
      } else if (type === 'data') {
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