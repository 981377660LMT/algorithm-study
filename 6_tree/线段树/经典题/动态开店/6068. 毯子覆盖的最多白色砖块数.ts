import { SegmentTreeDynamicLazy } from '../../template/动态开点/SegmentTreeDynamicLazy'

// 2209. 用地毯覆盖后的最少白色砖块
// https://leetcode.cn/problems/maximum-white-tiles-covered-by-a-carpet/
function maximumWhiteTiles(tiles: number[][], carpetLen: number): number {
  let min = 1e9
  let max = 0
  tiles.forEach(([left, right]) => {
    min = Math.min(min, left)
    max = Math.max(max, right)
  })

  // RangeAssignRangeSum
  const tree = new SegmentTreeDynamicLazy(min, max + 10, {
    e() {
      return 0
    },
    id() {
      return -1
    },
    op(e1, e2) {
      return e1 + e2
    },
    mapping(lazy, data, size) {
      return lazy === -1 ? data : lazy * size
    },
    composition(f, g) {
      return f === -1 ? g : f
    }
  })

  tiles.forEach(([left, right]) => {
    tree.updateRange(left, right + 1, 1)
  })

  let res = 0
  for (const [left] of tiles) res = Math.max(res, tree.query(left, left + carpetLen))
  return res
}

export {}
