// !离线查询区间乘积是否为立方数
// n,q<=2e5
// nums[i]<=1e6

import * as fs from 'fs'
import { useMoAlgo } from './useMoAlgo'

const { input } = useInput()
const [n, q] = input().split(' ').map(Number)
const nums = input().split(' ').map(Number)

const moAlgo = useMoAlgo<number, boolean>({
  add(value) {
    console.log(value)
  },
  remove(value) {
    console.log(value)
  },
  query() {
    return false
  }
})(nums)

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
