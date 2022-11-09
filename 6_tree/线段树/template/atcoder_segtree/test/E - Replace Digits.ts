/* eslint-disable no-console */

// https://atcoder.jp/contests/abl/tasks/abl_e

// 输入 n(≤2e5) 和 q(≤2e5)。
// 初始有一个长为 n 的字符串 s，
// !所有字符都是 1，s 的下标从 1 开始。
// 然后输入 q 个替换操作，每个操作输入 L,R (1≤L≤R≤n) 和 d (1≤d≤9)。
// !你需要把 s 的 [L,R] 内的所有字符替换为 d。
// !对每个操作，把替换后的 s 看成一个十进制数，输出这个数模 998244353 的结果。

// TLE
import * as fs from 'fs'
import { resolve } from 'path'
import { useAtcoderLazySegmentTree } from '../AtcoderLazySegmentTree'

function useInput(path?: string) {
  let data: string
  if (path) {
    data = fs.readFileSync(resolve(__dirname, path), 'utf8')
  } else {
    data = fs.readFileSync(process.stdin.fd, 'utf8')
  }

  const lines = data.split(/\r\n|\r|\n/)
  let lineId = 0
  const input = (): string => lines[lineId++]

  return {
    input
  }
}

const MOD = 998244353n
const { input } = useInput('a.txt')
const [n, q] = input().split(' ').map(Number)
const operations: Operation[] = []
for (let i = 0; i < q; i++) {
  const [left, right, target] = input().split(' ').map(Number)
  operations.push([left, right, BigInt(target)])
}

const res = replaceDigits(n, operations)
res.forEach(v => console.log(v))

//
//
type Operation = [left: number, right: number, target: bigint]
type Data = [sum: bigint, length: number]
type Lazy = bigint

function replaceDigits(n: number, operations: Operation[]): number[] {
  const pow10 = new BigUint64Array(n + 1).fill(1n)
  const pow10PreSum = new BigUint64Array(n + 1).fill(1n) // !第i个1111...是多少
  for (let i = 1; i <= n; i++) {
    pow10[i] = (pow10[i - 1] * 10n) % MOD
    pow10PreSum[i] = (pow10PreSum[i - 1] + pow10[i]) % MOD
  }

  const initNums = Array(n)
    .fill(null)
    .map<Data>(() => [1n, 1])
  const tree = useAtcoderLazySegmentTree<Data, Lazy>(initNums, {
    dataUnit: () => [0n, 0],
    lazyUnit: () => -1n,
    mergeChildren(data1, data2) {
      const [sum1, length1] = data1
      const [sum2, length2] = data2
      return [(sum1 * pow10[length2] + sum2) % MOD, length1 + length2]
    },
    updateData(parentLazy, childData) {
      if (parentLazy === -1n) return childData
      const length = childData[1]
      childData[0] = (parentLazy * pow10PreSum[length - 1]) % MOD
      return childData
    },
    updateLazy(parentLazy, childLazy) {
      return parentLazy >= 0n ? parentLazy : childLazy
    }
  })

  const res: number[] = []
  operations.forEach(([left, right, target]) => {
    tree.update(left - 1, right, target)
    res.push(Number(tree.queryAll()[0]))
  })

  return res
}
