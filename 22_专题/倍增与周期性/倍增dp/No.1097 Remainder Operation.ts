// https://yukicoder.me/problems/no/1097
// 给定一个数组和q次查询
// 初始时res为0，每次查询会执行k次操作:
// 将res加上nums[res%n]的值
// 求每次查询后res的值
// (n,q<=2e5,k<=1e12)

import * as fs from 'fs'
import { resolve } from 'path'
import { Doubling } from './Doubling'

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

const D = new Doubling(
  n,
  1e12 + 10,
  () => 0,
  (a, b) => a + b
)

const nums = input().split(' ').map(Number)
for (let i = 0; i < n; i++) {
  D.add(i, (i + nums[i]) % n, nums[i])
}

const q = Number(input())
for (let i = 0; i < q; i++) {
  const step = Number(input())
  const [_, value] = D.jump(0, step)
  console.log(value)
}
