/* eslint-disable no-inner-declarations */

// https://atcoder.jp/contests/arc028/tasks/arc028_4
// D - 注文の多い高橋商店
// 有n种商品，第i种商品有counts[i]个。给定购买的商品数need。
// 给定q次查询，每次查询给定i和x，求第i种商品恰好拿走x个的情况下，在n种物品中一共拿走need个物品的方案数。
// n,need,counts[i]<=2e3
// q<=5e5
// O(n*need*logn)

import * as fs from 'fs'

import { BoundedKnapsackRemovable } from '../BoundedKnapsackRemovable'

const MOD = 1e9 + 7

function 注文の多い高橋商店(counts: number[], need: number, queries: [pos: number, take: number][]): number[] {
  type Entry = { take: number; qid: number }
  const n = counts.length
  const q = queries.length
  const queryGroups: Entry[][] = Array(n)
  for (let i = 0; i < n; i++) queryGroups[i] = []
  for (let i = 0; i < q; i++) {
    const { 0: pos, 1: take } = queries[i]
    if (take > need || take > counts[pos]) continue
    queryGroups[pos].push({ take, qid: i })
  }

  const res = Array<number>(q).fill(0)
  const dp = new BoundedKnapsackRemovable(need, MOD)
  for (let i = 0; i < n; i++) {
    dp.add(1, counts[i])
  }
  for (let i = 0; i < n; i++) {
    const tmp = dp.copy()
    tmp.remove(1, counts[i])
    const group = queryGroups[i]
    group.forEach(({ take, qid }) => {
      res[qid] = tmp.query(need - take)
    })
  }

  return res
}

if (require.main === module) {
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
  const [n, need, q] = input().split(' ').map(Number)
  const counts = input().split(' ').map(Number)
  const queries = Array<[number, number]>(q)
  for (let i = 0; i < q; i++) {
    queries[i] = input().split(' ').map(Number) as [number, number]
    queries[i][0]--
  }

  const res = 注文の多い高橋商店(counts, need, queries)
  res.forEach(v => console.log(v))
}
