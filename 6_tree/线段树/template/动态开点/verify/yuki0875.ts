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

const [N, Q] = input().split(' ').map(Number)
const a = input().split(' ').map(Number)

const INF = 2e15
const seg = new SegmentTreeDynamic(0, a.length, () => INF, Math.min)
a.forEach((v, i) => seg.set(i, v))

for (let i = 0; i < Q; i++) {
  let [c, l, r] = input().split(' ').map(Number)
  if (c === 1) {
    l--
    r--
    const vl = seg.get(l)
    seg.set(l, seg.get(r))
    seg.set(r, vl)
  } else {
    l--

    const mn = seg.query(l, r)
    const a1 = seg.maxRight(l, n => n > mn)
    let a2 = seg.minLeft(r, n => n > mn)
    a2--
    if (a1 !== a2) {
      throw new Error('a1 !== a2')
    }
    console.log(a1 + 1)
  }
}
