// # 求第N项模1e9+7的值
// # a1=1 a2=1 a3=2 an=a(n-1)+a(n-2)+a(n-3)
// # N<=1e18

import { useInput } from '../../../20_杂题/atc競プロ/ts入力'
import { matqpow } from './matqpow'

const { input } = useInput()
const N = Number(input())

if (N <= 3) {
  console.log(N === 3 ? 2 : 1)
  process.exit(0)
}

const MOD = 1e9 + 7
const T = [
  [1, 1, 1],
  [1, 0, 0],
  [0, 1, 0],
]
const resT = matqpow(T, N - 3, MOD)
const [a3, a2, a1] = [2, 1, 1]
const res = (a3 * resT[0][0] + a2 * resT[0][1] + a1 * resT[0][2]) % MOD
console.log(res)
