// 町が N 個ある。町 i から町 Ai に移動することを K 回繰り返す。
// 町 1 から始めた時、最終的にどの町にたどり着くか？
// N<=2e5 K<=1e18

import { useInput } from '../ts入力'

const { input } = useInput(
  `
  6 727202214173249351
  6 5 2 5 3 2
`
)

let [a, b] = input().split(' ')
const n = Number(a)
let k = BigInt(b)

const towns = input()
  .split(' ')
  .map(v => Number(v) - 1)

let MAXJ = 0n
while (1n << MAXJ < k) MAXJ++
const dp = Array.from({ length: Number(MAXJ) + 1 }, () => new Int32Array(n + 1).fill(0))

for (let i = 0; i < n; i++) {
  dp[0][i] = towns[i]
}

for (let j = 0; j + 1 <= MAXJ; j++) {
  for (let i = 0; i < n; i++) {
    dp[j + 1][i] = dp[j][dp[j][i]]
  }
}

let res = 0
for (let bit = 0n; bit <= MAXJ; bit += 1n) {
  if ((k >> bit) & 1n) {
    res = dp[Number(bit)][res]
  }
}

console.log(res + 1)
