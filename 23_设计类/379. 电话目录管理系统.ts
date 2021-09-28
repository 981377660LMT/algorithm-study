import { RandomizedSet } from './特殊数据结构/380. 无重复元素O(1) 时间插入、删除和获取随机元素 copy'

class RandomSet extends RandomizedSet {
  has(key: number) {
    return this.map.has(key)
  }

  override getRandom(): number {
    if (!this.size) return -1
    const rand = super.getRandom()
    this.remove(rand)
    return rand
  }
}

class PhoneDirectory {
  private randomSet: RandomSet
  constructor(maxNumbers: number) {
    const nums = Array.from({ length: maxNumbers }, (_, i) => i)
    this.randomSet = new RandomSet()
    nums.forEach(num => this.randomSet.insert(num))
  }

  // 随机分配给用户一个未被使用的电话号码，获取失败请返回 -1
  // JS的集合无法随机pop元素 python集合可以
  // 故此处使用RandomizedSet
  get(): number {
    return this.randomSet.getRandom()
  }

  // 检查指定的电话号码是否被使用
  check(number: number): boolean {
    return !this.randomSet.has(number)
  }

  // 释放掉一个电话号码，使其能够重新被分配
  release(number: number): void {
    this.randomSet.remove(number)
  }
}
