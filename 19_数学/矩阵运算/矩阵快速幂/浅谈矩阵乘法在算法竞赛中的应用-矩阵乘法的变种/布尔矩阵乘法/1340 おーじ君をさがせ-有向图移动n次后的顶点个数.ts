import * as fs from 'fs'
import { resolve } from 'path'
import { BooleanSquareMatrixDense } from './BooleanSquareMatrix-dense'

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

// https://yukicoder.me/problems/no/1340
// 给定一个n个点m条边的有向图，求t步后可能所在的顶点个数(每一步必须移动到一个相邻点).
// 类似于`无向图中两点间是否存在长度为k的路径`.
// n<=100 m<=1e4 t<=1e18

const { input } = useInput('./a.txt')
const [n, m, t] = input().split(' ').map(Number)

const mat = new BooleanSquareMatrixDense(n)
for (let i = 0; i < m; i++) {
  const [a, b] = input().split(' ').map(Number)
  mat.set(a, b, true)
}

mat.ipow(t)

let res = 0
for (let i = 0; i < n; i++) res += +mat.get(0, i)
console.log(res)
