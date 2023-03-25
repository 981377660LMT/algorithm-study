// https://atcoder.jp/contests/abc294/tasks/abc294_g

import * as fs from 'fs'
import { resolve } from 'path'
import { BITArray2 } from '../BIT'
import { Tree } from '../Tree'

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

const { input } = useInput()

const n = Number(input())
const edges: [number, number, number][] = []
const hld = new Tree(n)
for (let i = 0; i < n - 1; i++) {
  const [u, v, w] = input().split(' ').map(Number)
  edges.push([u - 1, v - 1, w])
  hld.addEdge(u - 1, v - 1, w)
}
hld.build(0)

const bit = new BITArray2(n)
edges.forEach(([u, v, w]) => {
  const eid = hld.eid(u, v)
  bit.add(eid, eid + 1, w)
})
const q = Number(input())
for (let i = 0; i < q; i++) {
  const [op, u, v] = input().split(' ').map(Number)
  if (op === 1) {
    const e = edges[u - 1]
    const eid = hld.eid(e[0], e[1])
    const preW = bit.query(eid, eid + 1)
    const diff = v - preW
    bit.add(eid, eid + 1, diff)
  } else {
    let res = 0
    hld.enumeratePathDecomposition(u - 1, v - 1, false, (start, end) => {
      res += bit.query(start, end)
    })
    console.log(res)
  }
}
