import { SqrtDecomposition } from '../SqrtDecomposition'

// 给你两个下标从 0 开始的数组 nums1 和 nums2 ，和一个二维数组 queries 表示一些操作。总共有 3 种类型的操作：

// 操作类型 1 为 queries[i] = [1, l, r] 。你需要将 nums1 从下标 l 到下标 r 的所有 0 反转成 1 或将 1 反转成 0 。l 和 r 下标都从 0 开始。
// 操作类型 2 为 queries[i] = [2, p, 0] 。对于 0 <= i < n 中的所有下标，令 nums2[i] = nums2[i] + nums1[i] * p 。
// 操作类型 3 为 queries[i] = [3, 0, 0] 。求 nums2 中所有元素的和。
// 请你返回一个数组，包含所有第三种操作类型的答案。
// 区间反转+查询区间1的个数

function handleQuery(nums1: number[], nums2: number[], queries: number[][]): number[] {
  const n = nums1.length

  const sqrt = new SqrtDecomposition<number, 0 | 1>(n, (_, start, end) => {
    const curNums = nums1.slice(start, end)
    const len = end - start
    let ones = 0
    let lazyFlip = 0

    const created = () => {
      updated()
    }
    const updated = () => {
      ones = curNums.reduce((pre, cur) => pre + cur, 0)
    }
    const updateAll = (flip: 0 | 1) => {
      lazyFlip ^= flip
    }
    const updatePart = (start: number, end: number, flip: 0 | 1) => {
      for (let i = start; i < end; i++) {
        curNums[i] ^= flip
      }
    }
    const queryAll = () => (lazyFlip === 1 ? len - ones : ones)
    const queryPart = (start: number, end: number) => {
      let res = 0
      for (let i = start; i < end; i++) {
        curNums[i] ^= lazyFlip
        res += curNums[i] ? 1 : 0
      }
      return res
    }

    return {
      created,
      updated,
      updateAll,
      updatePart,
      queryAll,
      queryPart
    }
  })

  const res: number[] = []
  let allSum = nums2.reduce((a, b) => a + b, 0)
  queries.forEach(([op, a, b]) => {
    if (op === 1) {
      sqrt.update(a, b + 1, 1)
    } else if (op === 2) {
      let ones = 0
      sqrt.query(0, n, blockRes => {
        ones += blockRes
      })
      allSum += ones * a
    } else {
      res.push(allSum)
    }
  })
  return res
}
