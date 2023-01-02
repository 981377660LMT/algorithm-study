// https://atcoder.jp/contests/abc242/tasks/abc242_g
// !同じ色の服を着た 2 人からなるペアは最大何組作れるか答えよ。
// 静态查询区间 `元素的count //2` 的和

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
const q = Number(input())

// 维护区间的api
const manager = {
  pair: 0,
  counter: new Uint32Array(1e5 + 10),
  add(index: number) {
    const num = nums[index]
    this.pair -= Math.floor(this.counter[num] / 2)
    this.counter[num]++
    this.pair += Math.floor(this.counter[num] / 2)
  },
  remove(index: number) {
    const num = nums[index]
    this.pair -= Math.floor(this.counter[num] / 2)
    this.counter[num]--
    this.pair += Math.floor(this.counter[num] / 2)
  },
  query() {
    return this.pair
  }
}

const moAlgo = useMoAlgo(n, q, manager)

for (let i = 0; i < q; i++) {
  let [left, right] = input().split(' ').map(Number)
  left--, right--
  moAlgo.addQuery(left, right)
}

const res = moAlgo.work()
res.forEach(r => console.log(r))
