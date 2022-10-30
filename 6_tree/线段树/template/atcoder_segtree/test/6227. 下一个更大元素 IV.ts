import { useAtcoderLazySegmentTree } from '../AtcoderLazySegmentTree'

function secondGreaterElement(nums: number[]): number[] {
  const n = nums.length
  const maxTree = useAtcoderLazySegmentTree(nums, {
    dataUnit() {
      return 0
    },
    lazyUnit() {
      return 0
    },
    mergeChildren(data1, data2) {
      return Math.max(data1, data2)
    },
    updateData(parentLazy, childData) {
      return childData // 无需更新
    },
    updateLazy(parentLazy, childLazy) {
      return 0 // 无需更新
    }
  })

  const res = Array(n).fill(-1)
  for (let left = 0; left < n; left++) {
    // 右侧第二个严格大于 nums[left] 的元素
  }
  return res
}
