type NumSet = Set<number>
type Value = number

class RandomizedCollection {
  private map: Map<Value, NumSet>
  private arr: number[]
  private size: number
  constructor() {
    this.map = new Map()
    this.arr = []
    this.size = 0
  }
  /**
   *
   * @param val
   * @returns 返回 true 表示集合不包含 val
   */
  insert(val: number): boolean {
    // console.log(this.arr, this.map)
    !this.map.has(val) && this.map.set(val, new Set())
    this.map.get(val)!.add(this.size)
    this.arr.push(val)
    this.size++
    return this.map.get(val)!.size === 1
  }

  /**
   *
   * @param val 我们需要交换需要删除的数和数组末尾，
   * 并约定数组末尾的 n 项是被删除过的。（其中 n 为删除次数）
   */
  remove(val: number): boolean {
    console.log(this.arr, this.map)
    if (!this.map.has(val) || this.map.get(val)?.size === 0) return false
    const len = this.size
    const removeIndex = this.map.get(val)!.keys().next().value
    const lastVal = this.arr[len - 1]
    // 更新arr
    this.arr[removeIndex] = lastVal
    this.arr.pop()
    // 更新map:更新lastVal和val对应的map
    this.map.get(lastVal)?.delete(len - 1)
    lastVal !== val && this.map.get(lastVal)?.add(removeIndex)
    this.map.get(val)!.delete(removeIndex)
    this.map.get(val)!.size === 0 && this.map.delete(val)
    this.size--
    return true
  }

  /**
   * 随机返回现有集合中的一项（测试用例保证调用此方法时集合中至少存在一个元素）。每个元素应该有 `相同的概率` 被返回。
   * 数组 + 哈希表
   */
  getRandom(): number {
    return this.arr[~~(Math.random() * this.size)]
  }
}

const obj = new RandomizedCollection()
obj.insert(0)
obj.insert(1)
obj.remove(0)
obj.insert(2)
obj.remove(1)
console.log(obj)
console.log(obj.getRandom())
export {}
