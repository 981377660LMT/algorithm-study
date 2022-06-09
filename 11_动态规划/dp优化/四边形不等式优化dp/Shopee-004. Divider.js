const { readFileSync } = require('fs')
const iter = readlines()
const input = () => iter.next().value
function* readlines(path = 0) {
  const lines = readFileSync(path)
    .toString()
    .trim()
    .split(/\r\n|\r|\n/)

  yield* lines
}

const [n, k] = input().split(' ').map(Number)
const nums = input().split(' ').map(Number)
// const [n, k] = [4, 2]
// const nums = [1, 3, 2, 4]
const preSum = new Uint32Array(n + 1)
for (let i = 1; i <= n; i++) preSum[i] = preSum[i - 1] + nums[i - 1]

const dp = Array.from({ length: n + 5 }, () => Array(k + 5).fill(Infinity))
// 记录每个分组对应的当前最优的遍历起始点
const prePos = new Int16Array(k + 5).map((_, i) => i - 1)

for (let i = 1; i <= n; i++) {
  const maxG = Math.min(k, i)
  for (let g = 1; g <= maxG; g++) {
    if (g === 1) {
      dp[i][g] = preSum[i] * i
    } else {
      for (let preI = prePos[g]; preI < i; preI++) {
        const cand = dp[preI][g - 1] + (preSum[i] - preSum[preI]) * (i - preI)
        if (cand < dp[i][g]) {
          dp[i][g] = cand
          prePos[g] = preI
        }
      }
    }
  }
}

console.log(dp[n][k])
