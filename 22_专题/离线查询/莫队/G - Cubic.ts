/* eslint-disable no-shadow */
// !离线查询区间乘积是否为立方数
// n,q<=2e5
// nums[i]<=1e6

// !统计每个素数在区间乘积里的个数 注意到<=1e6时每个数的素因子不超过7个
// !总时间复杂度为 O(n*sqrt(q)*log(nums[i]))
import * as fs from 'fs'
import { usePrime } from '../../../19_数学/因数筛/prime'
import { useMoAlgo, WindowManager } from './useMoAlgo'

const { input } = useInput()
const [n, q] = input().split(' ').map(Number)
const nums = input().split(' ').map(Number)
const P = usePrime(1e6)
const primeFactors = Array(1e6 + 1).fill(null)
for (let i = 1; i <= 1e6; i++) {
  primeFactors[i] = P.getPrimeFactors(i)
}

let bad = 0
const counts = new Uint32Array(1e6 + 1) // 每种素因子个数
const colors = new Uint8Array(1e6 + 1) // 每种素因子是否贡献过(染色)

const check = (factor: number): void => {
  if (counts[factor] % 3 !== 0) {
    if (colors[factor] === 0) {
      colors[factor] = 1
      bad++
    }
  } else if (colors[factor] === 1) {
    colors[factor] = 0
    bad--
  }
}
const windowManager: WindowManager<boolean> = {
  add(index) {
    const value = nums[index]
    for (const [factor, count] of primeFactors[value]) {
      counts[factor] += count
      check(factor)
    }
  },
  remove(index) {
    const value = nums[index]
    for (const [factor, count] of primeFactors[value]) {
      counts[factor] -= count
      check(factor)
    }
  },
  query() {
    return bad === 0
  }
}
const moAlgo = useMoAlgo(n, q, windowManager)

for (let _ = 0; _ < q; _++) {
  let [left, right] = input().split(' ').map(Number)
  left--, right--
  moAlgo.addQuery(left, right)
}

console.log(
  moAlgo
    .work()
    .map(flag => (flag ? 'Yes' : 'No'))
    .join('\n')
)

function useInput(debugCase?: string) {
  const data = debugCase === void 0 ? fs.readFileSync(process.stdin.fd, 'utf8') : debugCase
  const dataIter = _makeIter(data)

  function input(): string {
    return dataIter.next().value.trim()
  }

  function* _makeIter(str: string): Generator<string, string, undefined> {
    yield* str.trim().split(/\r\n|\r|\n/)
    return ''
  }

  return {
    input
  }
}
