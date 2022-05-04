import { bisectLeft } from '../../9_排序和搜索/二分/7_二分搜索寻找最左插入位置'
import { bisectRight } from '../../9_排序和搜索/二分/7_二分搜索寻找最插右入位置'

// 询问次数q<=1e4
// 数组长度<=2e4
class MajorityChecker {
  private readonly indexMap: Map<number, number[]> // 每个元素出现的位置

  constructor(arr: number[]) {
    this.indexMap = new Map()
    for (let i = 0; i < arr.length; i++) {
      const key = arr[i]
      !this.indexMap.has(key) && this.indexMap.set(key, [])
      this.indexMap.get(key)!.push(i)
    }
  }

  // 返回在 arr[left], arr[left+1], ..., arr[right] 中至少出现阈值次数 threshold 的元素
  // threshold 始终比子序列长度的一半还要大
  query(left: number, right: number, threshold: number): number {
    for (const [num, indexes] of this.indexMap.entries()) {
      if (indexes.length < threshold) continue
      const leftIndex = bisectLeft(indexes, left)
      const rightIndex = bisectRight(indexes, right)
      if (rightIndex - leftIndex >= threshold) return num
    }

    return -1
  }
}

export {}

// 预处理每个元素出现的位置是优化查询的常用手法
