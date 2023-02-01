// 然后输入 m 个操作：
// 操作 1 形如 1 x y k，表示把 a 的区间 [x,x+k-1] 的元素拷贝到 b 的区间 [y,y+k-1] 上（输入保证下标不越界）。
// 操作 2 形如 2 x，输出 b[x]。
// 区间赋值 单点查询

/* eslint-disable @typescript-eslint/no-var-requires */
const { readFileSync } = require('fs')

const iter = readlines()
const input = (): string => iter.next().value
function* readlines(path = 0) {
  const lines = readFileSync(path)
    .toString()
    .trim()
    .split(/\r\n|\r|\n/)

  yield* lines
}

const n = Number(input())
const nums1 = new Int32Array([0, ...input().split(' ').map(Number)])
const nums2 = new Int32Array(n + 1).fill(-1)
const m = Number(input())
for (let _ = 0; _ < m; _++) {
  const [opt, ...rest] = input().split(' ').map(Number)
  if (opt === 1) {
    // !把 A 序列中从下标 x 位置开始的连续 k 个元素粘贴到 B 序列中从下标 y 开始的连续 k 个位置上。
    // 输入数据可能会出现粘贴后 k 个元素超出 B 序列原有长度的情况，超出部分可忽略
    const [xLen, xStart, yStart] = rest
    copy(nums1, nums2, xStart, xLen, yStart)
  } else {
    // !表示询问B序列下标 x 处的值是多少
    const qi = rest[0]
    console.log(nums2[qi])
  }
}

function copy(nums1: Int32Array, nums2: Int32Array, i1: number, k: number, i2: number): void {
  const len = Math.min(k, nums1.length - i1, nums2.length - i2)
  nums2.set(nums1.subarray(i1, i1 + len), i2)
}

export {}
