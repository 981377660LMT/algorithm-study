import { bisectLeft } from './二分api/7_二分搜索寻找最左插入位置'
import { bisectRight } from './二分api/7_二分搜索寻找最插右入位置'

class MajorityChecker {
  private location: Map<number, number[]> // 每个元素出现的位置

  constructor(arr: number[]) {
    this.location = new Map()
    for (let i = 0; i < arr.length; i++) {
      const key = arr[i]
      !this.location.has(key) && this.location.set(key, [])
      this.location.get(key)!.push(i)
    }
  }

  // 返回在 arr[left], arr[left+1], ..., arr[right] 中至少出现阈值次数 threshold 的元素
  // threshold 始终比子序列长度的一半还要大
  query(left: number, right: number, threshold: number): number {
    for (const [num, positions] of this.location.entries()) {
      if (positions.length < threshold) continue
      const leftIndex = bisectLeft(positions, left)
      const rightIndex = bisectRight(positions, right)
      if (rightIndex - leftIndex >= threshold) return num
    }
    return -1
  }
}

export {}

// 预处理每个元素出现的位置是优化查询的常用手法
