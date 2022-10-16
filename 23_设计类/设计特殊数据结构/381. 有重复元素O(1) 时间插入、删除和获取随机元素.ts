class RandomizedCollection {
  private readonly indexMap: Map<number, Set<number>> = new Map()
  private readonly nums: number[] = []
  private size = 0

  /**
   *
   * @param val
   * @returns 返回 true 表示集合不包含 val
   */
  insert(val: number): boolean {
    !this.indexMap.has(val) && this.indexMap.set(val, new Set())
    this.indexMap.get(val)!.add(this.size)
    this.nums.push(val)
    this.size++
    return this.indexMap.get(val)!.size === 1
  }

  /**
   *
   * @param val 我们需要交换需要删除的数和数组末尾，
   * 并约定数组末尾的 n 项是被删除过的。（其中 n 为删除次数）
   */
  remove(val: number): boolean {
    console.log(this.nums, this.indexMap)
    if (!this.indexMap.has(val) || this.indexMap.get(val)?.size === 0) return false

    const lastValIndex = this.size - 1
    const lastVal = this.nums[lastValIndex]
    const removeIndex = this.indexMap.get(val)!.keys().next().value

    // 待删除的元素删除removeIndex
    this.indexMap.get(val)?.delete(removeIndex)
    // 先加再删，处理一个元素的情况
    // 最后一个元素值加入removeIndex
    // 最后一个元素值移除之前最后一个元素值的位置
    this.indexMap.get(lastVal)?.add(removeIndex).delete(lastValIndex)

    // removeIndex处置为最后一个元素值
    this.nums[removeIndex] = lastVal
    // pop操作
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

// const obj = new RandomizedCollection()
// obj.insert(0)
// obj.insert(1)
// obj.remove(0)
// obj.insert(2)
// obj.remove(1)
// console.log(obj)
// console.log(obj.getRandom())
const obj = new RandomizedCollection()
obj.insert(0)
obj.remove(0)
obj.insert(-1)
obj.remove(0)
console.log(obj.getRandom())
console.log(obj.getRandom())
console.log(obj.getRandom())
console.log(obj.getRandom())
console.log(obj.getRandom())
console.log(obj)
export {}
