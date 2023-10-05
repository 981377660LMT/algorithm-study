// 第 i 个掉落的方块（positions[i] = (left, side_length)）是正方形
// 1 <= positions.length <= 1000.
// 1 <= positions[i][0] <= 10^8.
// 1 <= positions[i][1] <= 10^6.
// https://leetcode.cn/problems/falling-squares/

import { createRangeUpdateRangeMax } from '../../template/atcoder_segtree/SegmentTreeUtils'
import { SegmentTreeDynamicLazy } from '../../template/动态开点/SegmentTreeDynamicLazy'

const INF = 2e15

function fallingSquares(positions: number[][]): number[] {
  const res = Array<number>(positions.length).fill(0)
  const tree = new SegmentTreeDynamicLazy(0, 2e8 + 10, {
    e() {
      return 0
    },
    id() {
      return -INF
    },
    op(x, y) {
      return Math.max(x, y)
    },
    mapping(f, x) {
      return f === -INF ? x : Math.max(f, x)
    },
    composition(f, g) {
      return f === -INF ? g : Math.max(f, g)
    }
  })

  positions.forEach(([left, size], i) => {
    const right = left + size - 1
    const preHeihgt = tree.query(left, right + 1)
    tree.updateRange(left, right + 1, preHeihgt + size)
    res[i] = tree.queryAll()
  })

  return res
}

function fallingSquares2(positions: number[][]): number[] {
  const res = Array<number>(positions.length).fill(0)
  const allNums = new Set<number>()
  positions.forEach(([left, size]) => {
    allNums.add(left)
    allNums.add(left + size - 1)
  })

  const [rank, count] = discretize([...allNums])
  const tree = createRangeUpdateRangeMax(count)

  positions.forEach(([left, size], i) => {
    const right = left + size - 1
    const rankLeft = rank(left)
    const rankRight = rank(right + 1)
    const preHeihgt = tree.query(rankLeft, rankRight)
    tree.update(rankLeft, rankRight, preHeihgt + size)
    res[i] = tree.queryAll()
  })

  return res
}

/**
 * (松)离散化.
 * @returns
 * rank: 给定一个数,返回它的排名`(0-count)`.
 * count: 离散化(去重)后的元素个数.
 */
function discretize(nums: number[]): [rank: (num: number) => number, count: number] {
  const allNums = [...new Set(nums)].sort((a, b) => a - b)
  const rank = (num: number) => {
    let left = 0
    let right = allNums.length - 1
    while (left <= right) {
      const mid = (left + right) >>> 1
      if (allNums[mid] >= num) {
        right = mid - 1
      } else {
        left = mid + 1
      }
    }
    return left
  }
  return [rank, allNums.length]
}

export {}

if (require.main === module) {
  console.log(
    fallingSquares([
      [1, 2],
      [2, 3],
      [6, 1]
    ])
  )
  console.log(
    fallingSquares([
      [100, 100],
      [200, 100]
    ])
  )
}
