// https://yukicoder.me/problems/no/789
import * as fs from 'fs'
import { resolve } from 'path'
import { SegmentTreeDynamic } from '../SegmentTreeDynamicSparse'

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

const seg = new SegmentTreeDynamic<number>(
  0,
  1e9 + 10,
  () => 0,
  (a, b) => a + b
)
const n = Number(input())
let res = 0
for (let i = 0; i < n; i++) {
  const arr = input().split(' ').map(Number)
  if (arr[0] === 0) {
    // seg.update(arr[1], arr[2])
    const pre = seg.get(arr[1])
    seg.set(arr[1], arr[2] + pre)
  } else {
    res += seg.query(arr[1], arr[2] + 1)
  }
}
console.log(res)
