// https://atcoder.jp/contests/abc295/tasks/abc295_g
// G - Minimum Reachable City-有向树加边后抵达的最小编号
// 给定一颗特殊的外向树,第i条边连接p和i+1且p<i+1.(沿着树边走,编号递增)
// 现在给定q个操作
// !1 u v: 在u和v之间加一条边,保证连边之前可以在最开始的树上从v到达u
// !2 u: 询问从u出发,能到达的最小编号的点是多少.

import * as fs from 'fs'
import { resolve } from 'path'

import { UnionFindRangeOnTree } from './RangeUnionFindOnTree'

function solve(n: number, parents: number[], queries: [number, number, number][]): number[] {
  const uf = new UnionFindRangeOnTree(n, parents)
  const groupMin = Array<number>(n)
  for (let i = 0; i < n; i++) groupMin[i] = i

  const res: number[] = []
  for (let i = 0; i < queries.length; i++) {
    const { 0: op, 1: child, 2: ancestor } = queries[i]
    if (op === 1) {
      uf.unionRange(ancestor, child, (mergeTo, mergeFrom) => {
        groupMin[mergeTo] = Math.min(groupMin[mergeTo], groupMin[mergeFrom])
      })
    } else {
      const root = uf.find(child)
      res.push(groupMin[root])
    }
  }

  return res
}

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

const { input } = useInput('a.txt')

const n = Number(input())
const parents = input()
  .split(' ')
  .map(Number)
  .map(v => v - 1)
parents.unshift(-1)

const q = Number(input())
const queries = Array(q)
for (let i = 0; i < q; i++) {
  const [op, ...rest] = input().split(' ').map(Number)
  if (op === 1) {
    const [u, v] = rest
    queries[i] = [op, u - 1, v - 1]
  } else {
    const [u] = rest
    queries[i] = [op, u - 1, 0]
  }
}

const res = solve(n, parents, queries)
res.forEach(v => console.log(v + 1))
