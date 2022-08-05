// https://atcoder.jp/contests/abc242/tasks/abc242_g
// 静态查询区间 `元素频率//2` 的最大值
import * as fs from 'fs'
import { useMoAlgo } from './useMoAlgo'

function useInput(debugCase?: string) {
  const data = debugCase == void 0 ? fs.readFileSync(process.stdin.fd, 'utf8') : debugCase
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

const { input } = useInput()

const n = Number(input())
const nums = input().split(' ').map(Number)

// 维护区间的api
const manager = {
  pair: 0,
  counter: new Map<number, number>(),
  add(num: number) {
    this.pair -= Math.floor((this.counter.get(num) || 0) / 2)
    this.counter.set(num, (this.counter.get(num) || 0) + 1)
    this.pair += Math.floor((this.counter.get(num) || 0) / 2)
  },
  remove(num: number) {
    this.pair -= Math.floor((this.counter.get(num) || 0) / 2)
    this.counter.set(num, (this.counter.get(num) || 0) - 1)
    this.pair += Math.floor((this.counter.get(num) || 0) / 2)
  },
  query() {
    return this.pair
  }
}

const moAlgo = useMoAlgo(manager)(nums)

const q = Number(input())
for (let i = 0; i < q; i++) {
  let [left, right] = input().split(' ').map(Number)
  left--, right--
  moAlgo.addQuery(left, right)
}

const res = moAlgo.work()
res.forEach(r => console.log(r))
