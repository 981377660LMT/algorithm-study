// !2213. 由单个字符重复的最长子字符串
// https://leetcode.cn/problems/longest-substring-of-one-repeating-character/solution/typescript-ke-duo-li-shu-by-981377660lmt-n77n/
// 1. 用珂朵莉树维护区间字符类型，有序数组维护每个连续段的长度。
// 2. 因为每次单点修改最多影响到左、中、右三个段，为了避免分类讨论，
//    每次修改前，遍历左段起点到右段终点的每个区间，删除这些区间的长度，
//    单点修改后，再遍历左段起点到右段终点的每个区间，把这些区间长度加回来。
// 3. 每次查询的单个字符重复的最长子字符串，就是有序数组中的最大值

import { enumerateGroup } from '../../../0_数组/数组api/groupby'
import { ErasableHeap } from '../../../8_heap/ErasableHeap'
import { ODT } from '../ODT-fastset'

function longestRepeating(s: string, queryCharacters: string, queryIndices: number[]): number[] {
  const n = s.length
  const q = queryIndices.length
  const res = Array(q).fill(1)

  const odt = new ODT(n, -1) // !珂朵莉树维护区间字符类型
  const lens = new ErasableHeap<number>((a, b) => a > b) // !有序容器维护每个连续段的长度
  enumerateGroup(s, (start, end) => {
    const value = s[start].charCodeAt(0) - 97
    odt.set(start, end, value)
    lens.push(end - start)
  })

  for (let i = 0; i < q; i++) {
    const target = queryCharacters.charCodeAt(i) - 97
    const pos = queryIndices[i]

    // !每次更新最多影响左中右三段区间
    // !先删除这三段区间的长度，修改后，再添加这三段区间的长度
    // 这种做法无需分类讨论
    const [start, end] = odt.get(pos)!
    const leftSeg = odt.get(start - 1)
    const rightSeg = odt.get(end)
    const first = leftSeg ? leftSeg[0] : 0
    const last = rightSeg ? rightSeg[1] : n
    odt.enumerateRange(first, last, (start, end, value) => {
      if (value !== -1) lens.remove(end - start)
    })

    odt.set(pos, pos + 1, target)

    odt.enumerateRange(first, last, (start, end, value) => {
      if (value !== -1) lens.push(end - start)
    })

    res[i] = lens.peek()
  }

  return res
}

// s = "babacc", queryCharacters = "bcb", queryIndices = [1,3,3]
// console.log(longestRepeating('babacc', 'bcb', [1, 3, 3]))

const res = longestRepeating(
  'mrbkgpioaeypvvvwnlegkjkhxgilqlzwmnusspcrqiaapkzljfodokdosufidsxfbygmnaxhsvmejdmcpqhbghtkoyzwgzgt',
  'csfiuruhfmxsdeiftbjaopdxndjfalmubseikqotnrisayzrlwgnsmqqavetaaapsifyjcernvxbpgbmnffuwaaruy',
  [
    14, 43, 39, 65, 4, 15, 80, 55, 24, 51, 91, 41, 29, 48, 41, 74, 4, 49, 28, 1, 28, 75, 57, 72, 61,
    0, 45, 43, 19, 6, 28, 87, 47, 27, 85, 73, 34, 26, 84, 47, 55, 34, 34, 48, 79, 23, 3, 11, 44, 2,
    61, 68, 44, 92, 51, 13, 32, 86, 41, 23, 59, 73, 6, 12, 79, 35, 5, 24, 6, 5, 32, 52, 75, 76, 80,
    4, 83, 41, 77, 57, 52, 88, 86, 21, 6, 48, 9, 61, 50, 48
  ]
)
console.log(res)
