// snap 备忘录模式
// [map,map,map,...,map,map]

import { bisectLeft } from '../9_排序和搜索/二分/7_二分搜索寻找最左插入位置'

// 字典数组 + 二分查找
type SnapId = number
type Value = number

class SnapshotArray {
  private actions: Map<SnapId, Value>[]
  private snapId: number

  /**
   *
   * @param length 初始化一个与指定长度相等的 类数组 的数据结构。初始时，每个元素都等于 0
   */
  constructor(length: number) {
    this.actions = Array.from({ length }, () => new Map<SnapId, Value>())
    this.snapId = 0
  }

  /**
   *
   * @param index  会将指定索引 index 处的元素设置为 val。
   * @param val
   * data[index].keys() 为 index 变更的所有记录
   */
  set(index: number, val: number): void {
    this.actions[index].set(this.snapId, val)
  }

  /**
   * @description 获取该数组的快照，并返回快照的编号 snap_id（快照号是调用 snap() 的总次数减去 1）
   * @summary
   * 由于快照功能 snap() 的调用次数可能很多，
   * 所以我们如果采用每次快照都整体保存一次数组的方法，
   * 无论在时间复杂度还是空间复杂度上，都是行不通的
   * 优化的方法是，只保存每次快照变化的部分 (action)
   */
  snap(): number {
    this.snapId++
    return this.snapId - 1
  }

  /**
   *
   * @param index
   * @param snapId 根据指定的 snap_id 选择快照，并返回该快照指定索引 index 的值
   * 如果快照恰好存在, 直接返回
   * 不存在则进行二分搜索, 查找快照前最后一次修改
   */
  get(index: number, snapId: number): number | null {
    const map = this.actions[index]
    if (map.has(snapId)) return map.get(snapId)!
    const snapIds = [...this.actions[index].keys()]
    const pos = bisectLeft(snapIds, snapId)
    return map.get(snapIds[pos - 1]) || null
  }
}

const snapshotArr = new SnapshotArray(3) // 初始化一个长度为 3 的快照数组
snapshotArr.set(0, 5) // 令 array[0] = 5
snapshotArr.snap() // 获取快照，返回 snap_id = 0
snapshotArr.set(0, 6)
console.log(snapshotArr.get(0, 0)) // 获取 snap_id = 0 的快照中 array[0] 的值，返回 5

export default void 0
