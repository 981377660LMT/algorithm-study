/**
 * 哈希堆的思想来源
 */
class RandomizedSet {
  private readonly indexMap: Map<number, number> = new Map()
  private readonly nums: number[] = []
  private size = 0

  insert(val: number): boolean {
    if (this.indexMap.has(val)) return false
    this.indexMap.set(val, this.size)
    this.nums.push(val)
    this.size++
    return true
  }

  /**
   *
   * @param val 我们需要交换需要删除的数和数组末尾，
   * 并约定数组末尾的 n 项是被删除过的。（其中 n 为删除次数）
   */
  remove(val: number): boolean {
    if (!this.indexMap.has(val)) return false
    const removeIndex = this.indexMap.get(val)!
    const lastVal = this.nums[this.nums.length - 1]
    // 更新map
    this.indexMap.set(lastVal, removeIndex)
    this.indexMap.delete(val)
    // 更新arr
    this.nums[removeIndex] = lastVal
    this.nums.pop()
    this.size--
    return true
  }

  /**
   * 随机返回现有集合中的一项（测试用例保证调用此方法时集合中至少存在一个元素）。每个元素应该有 `相同的概率` 被返回。
   * 数组 + 哈希表
   */
  getRandom(): number {
    return this.nums[~~(Math.random() * this.size)]
  }
}

export { RandomizedSet }
