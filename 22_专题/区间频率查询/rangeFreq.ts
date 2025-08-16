/**
 * 区间频率查询类.
 * 支持查询指定区间内某个值的出现次数、第一个位置和最后一个位置.
 */
class RangeFreq<T> {
  private readonly _valueToIndexes: Map<T, number[]>

  constructor(nums: T[]) {
    this._valueToIndexes = new Map()
    for (let i = 0; i < nums.length; i++) {
      const v = nums[i]
      if (!this._valueToIndexes.has(v)) this._valueToIndexes.set(v, [])
      this._valueToIndexes.get(v)!.push(i)
    }
  }

  /**
   * 查询区间 [start, end) 中值为 value 的元素个数.
   */
  query(start: number, end: number, value: T): number {
    if (start >= end) return 0
    const pos = this._valueToIndexes.get(value)
    if (!pos) return 0
    return this._bisectLeft(pos, end) - this._bisectLeft(pos, start)
  }

  /**
   * 查找区间 [start, end) 中值为 value 的第一个位置.
   * @returns 第一个位置的索引，如果不存在返回 -1.
   */
  findFirst(start: number, end: number, value: T): number {
    if (start >= end) return -1
    const pos = this._valueToIndexes.get(value)
    if (!pos) return -1
    const idx = this._bisectLeft(pos, start)
    if (idx < pos.length && pos[idx] < end) {
      return pos[idx]
    }
    return -1
  }

  /**
   * 查找区间 [start, end) 中值为 value 的最后一个位置.
   * @returns 最后一个位置的索引，如果不存在返回 -1.
   */
  findLast(start: number, end: number, value: T): number {
    if (start >= end) return -1
    const pos = this._valueToIndexes.get(value)
    if (!pos) return -1
    const idx = this._bisectLeft(pos, end)
    if (idx > 0 && pos[idx - 1] >= start) {
      return pos[idx - 1]
    }
    return -1
  }

  private _bisectLeft(nums: number[], target: number): number {
    let left = 0
    let right = nums.length
    while (left < right) {
      const mid = (left + right) >>> 1
      if (nums[mid] < target) {
        left = mid + 1
      } else {
        right = mid
      }
    }
    return left
  }
}

if (require.main === module) {
  // 测试函数
  function testRangeFreq(): void {
    // 测试基本功能
    const nums = [1, 2, 3, 2, 4, 2, 5]
    const rf = new RangeFreq(nums)

    // 测试query方法
    console.assert(rf.query(0, 7, 2) === 3, '元素2在整个数组中出现3次')
    console.assert(rf.query(1, 5, 2) === 2, '元素2在区间[1,5)中出现2次')
    console.assert(rf.query(0, 3, 2) === 1, '元素2在区间[0,3)中出现1次')
    console.assert(rf.query(0, 7, 6) === 0, '元素6不存在')

    // 测试findFirst方法
    console.assert(rf.findFirst(0, 7, 2) === 1, '元素2的第一个位置是1')
    console.assert(rf.findFirst(2, 7, 2) === 3, '在区间[2,7)中元素2的第一个位置是3')
    console.assert(rf.findFirst(4, 7, 2) === 5, '在区间[4,7)中元素2的第一个位置是5')
    console.assert(rf.findFirst(0, 7, 6) === -1, '元素6不存在')
    console.assert(rf.findFirst(6, 7, 2) === -1, '区间[6,7)中没有元素2')

    // 测试findLast方法
    console.assert(rf.findLast(0, 7, 2) === 5, '元素2的最后一个位置是5')
    console.assert(rf.findLast(0, 4, 2) === 3, '在区间[0,4)中元素2的最后一个位置是3')
    console.assert(rf.findLast(0, 2, 2) === 1, '在区间[0,2)中元素2的最后一个位置是1')
    console.assert(rf.findLast(0, 7, 6) === -1, '元素6不存在')
    console.assert(rf.findLast(0, 1, 2) === -1, '区间[0,1)中没有元素2')

    // 测试边界情况
    console.assert(rf.query(3, 3, 2) === 0, '空区间')
    console.assert(rf.findFirst(3, 3, 2) === -1, '空区间')
    console.assert(rf.findLast(3, 3, 2) === -1, '空区间')

    // 测试字符串类型
    const strNums = ['a', 'b', 'c', 'b', 'd', 'b']
    const strRf = new RangeFreq(strNums)
    console.assert(strRf.query(0, 6, 'b') === 3)
    console.assert(strRf.findFirst(0, 6, 'b') === 1)
    console.assert(strRf.findLast(0, 6, 'b') === 5)

    console.log('所有测试通过!')
  }

  testRangeFreq()
}

export { RangeFreq }
