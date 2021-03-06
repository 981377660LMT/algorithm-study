// https://atcoder.jp/contests/abc242/tasks/abc242_g

import { useInput } from '../../../20_杂题/atc競プロ/ts入力'
import { useMoAlgo } from './useMoAlgo'

// 静态查询区间 `元素频率//2` 的最大值
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
  },
}

const moAlgo = useMoAlgo(nums, manager)

const q = Number(input())
for (let i = 0; i < q; i++) {
  let [left, right] = input().split(' ').map(Number)
  left--, right--
  moAlgo.addQuery(left, right)
}

const res = moAlgo.work()
res.forEach(r => console.log(r))
