// https://leetcode.cn/problems/queue-reconstruction-by-height/
// 406. 根据身高重建队列-线段树树上二分 O(nlogn)
// people[i] = [hi, ki] 表示第 i 个人的身高为 hi ，前面 正好 有 ki 个身高大于或等于 hi 的人。
// 返回的队列应该格式化为数组 queue ，其中 queue[j] = [hj, kj] 是队列中第 j 个人的属性（queue[0] 是排在队列前面的人）。

import { SegmentTree01 } from './SegmentTree01'

function reconstructQueue(people: number[][]): number[][] {
  const n = people.length
  people.sort((a, b) => a[0] - b[0] || -(a[1] - b[1]))

  const tree = new SegmentTree01(new Uint8Array(n).fill(1))
  const res = Array.from<unknown, [height: number, preCount: number]>({ length: n }, () => [0, 0])
  people.forEach(([height, preCount]) => {
    const pos = tree.kth(1, preCount + 1)
    res[pos] = [height, preCount]
    tree.flip(pos, pos + 1)
  })

  return res
}

if (require.main === module) {
  console.log(
    reconstructQueue([
      [7, 0],
      [4, 4],
      [7, 1],
      [5, 0],
      [6, 1],
      [5, 2]
    ])
  )
}

export {}
