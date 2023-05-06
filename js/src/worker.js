const { parentPort, workerData, threadId } = require("worker_threads");
const ethers = require('ethers')

function areLastNCharsSame(str, n) {
    const lastNChars = str.slice(-n);
    const lastChar = lastNChars.charAt(0).toLowerCase();
    return lastNChars.split('').every(char => char === lastChar);
}

function areFirstNCharsSame(str, n) {
    const slicePrefix = 2
    const firstNChars = str.slice(slicePrefix, n + slicePrefix);
    const firstChar = firstNChars.charAt(0).toLowerCase();
    return firstNChars.split('').every(char => char === firstChar);
}

function hasMoreThanNRepeatedCharsInARow(str, N) {
    const regex = new RegExp(`(.)\\1{${N - 1},}`, 'g');
    return regex.test(str);
}

function scoreAddress(address, numberChars) {
    let score = 0
    const hasRepeatedChars = hasMoreThanNRepeatedCharsInARow(address, numberChars)
    if (hasRepeatedChars) {
        const firstCharsSame = areFirstNCharsSame(address, numberChars)
        const lastCharsSame = areLastNCharsSame(address, numberChars)
        if (firstCharsSame && lastCharsSame) {
            score += 100
        } else if (firstCharsSame) {
            score += 5
        } else if (lastCharsSame) {
            score += 10
        } else {
            score += 1 // uncomment if need chars in the middle
        }

    }
    return score

}

function findAddresses(numberChars, limit) {
    const logEvery = 10000
    let foundAddresses = 0
    for (let i = 0; i < limit; i++) {
        if (i % logEvery === 0) {
            parentPort.postMessage({type: 'processed', message: {processed: i}})
        }
        const privateKey = ethers.hexlify(ethers.randomBytes(32));
        const address = ethers.computeAddress(privateKey).toLowerCase();
        const score = scoreAddress(address, numberChars)
        if (score) {
            foundAddresses++
            const data = { score, address, privateKey }
            parentPort.postMessage({type: 'data', message: data})
        }

    }
}

findAddresses(workerData.numberChars, workerData.limit)