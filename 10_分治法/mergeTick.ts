/* eslint-disable no-inner-declarations */
// https://www.luogu.com.cn/blog/_post/269778
// 操作1：给区间 [l,r) 中的每个数都加上 x.要求O(n)的时间复杂度.
// !操作2：查询数组排序后的结果.要求O(n)的时间复杂度.
// 解法：将sorted按照是否在更改区间内分类，分别压入两个数组，然后对这两个数组进行归并

type SortedItem = {
  value: number
  index: number
}

class MergeTrick {
  private readonly _nums: number[]
  private readonly _sortedItems: SortedItem[]
  private _sortedNums: number[]
  private dirty = false

  constructor(nums: ArrayLike<number>) {
    const copy = Array(nums.length)
    for (let i = 0; i < nums.length; i++) copy[i] = nums[i]
    const sortedItems: SortedItem[] = Array(nums.length)
    for (let i = 0; i < nums.length; i++) sortedItems[i] = { value: nums[i], index: i }
    sortedItems.sort((a, b) => a.value - b.value)
    const sortedNums = Array(nums.length)
    for (let i = 0; i < nums.length; i++) sortedNums[i] = sortedItems[i].value
    this._nums = copy
    this._sortedItems = sortedItems
    this._sortedNums = sortedNums
  }

  add(start: number, end: number, delta: number): void {
    this.dirty = true
    const n = this._nums.length
    const modified: SortedItem[] = Array(end - start)
    const unmodified: SortedItem[] = Array(n - (end - start))
    for (let i = 0, ptr1 = 0, ptr2 = 0; i < n; i++) {
      const item = this._sortedItems[i]
      if (item.index >= start && item.index < end) {
        item.value += delta
        modified[ptr1++] = item
        this._nums[item.index] += delta
      } else {
        unmodified[ptr2++] = item
      }
    }

    let i1 = 0
    let i2 = 0
    let k = 0
    while (i1 < modified.length && i2 < unmodified.length) {
      if (modified[i1].value < unmodified[i2].value) {
        this._sortedItems[k] = modified[i1]
        i1++
      } else {
        this._sortedItems[k] = unmodified[i2]
        i2++
      }
      k++
    }

    while (i1 < modified.length) {
      this._sortedItems[k] = modified[i1]
      i1++
      k++
    }

    while (i2 < unmodified.length) {
      this._sortedItems[k] = unmodified[i2]
      i2++
      k++
    }
  }

  /**
   * 返回原始数组.
   */
  get nums(): number[] {
    return this._nums
  }

  /**
   * 返回排序后的数组.
   */
  get sortedNums(): number[] {
    if (!this.dirty) return this._sortedNums
    this.dirty = false
    const res = Array(this._nums.length)
    for (let i = 0; i < res.length; i++) res[i] = this._sortedItems[i].value
    this._sortedNums = res
    return res
  }
}

export { MergeTrick }

if (require.main === module) {
  class _MergeTrickNaive {
    private _sorted: number[]

    constructor(nums: number[]) {
      this._sorted = nums.sort((a, b) => a - b)
    }

    add(start: number, end: number, delta: number): void {
      for (let i = start; i < end; i++) {
        this._sorted[i] += delta
      }
      this._sorted.sort((a, b) => a - b)
    }

    sorted(): number[] {
      return this._sorted.slice()
    }
  }

  // // time
  // const bigArr = Array(2e3)
  // for (let i = 0; i < bigArr.length; i++) {
  //   bigArr[i] = Math.floor(Math.random() * 1e9)
  // }
  // const mt2 = new MergeTrick(bigArr)

  // const mt3 = new _MergeTrickNaive(bigArr)

  // console.time('merge trick')
  // for (let i = 0; i < 1e5; i++) {
  //   mt2.add(0, 800, 1)
  //   mt2.sortedNums
  // }
  // console.timeEnd('merge trick')

  // console.time('naive')
  // for (let i = 0; i < 1e5; i++) {
  //   mt3.add(0, 800, 1)
  //   mt3.sorted()
  // }
  // console.timeEnd('naive')

  check()
  function check(): void {
    const arr = Array(4)
    for (let i = 0; i < arr.length; i++) {
      arr[i] = Math.floor(Math.random() * 4)
    }
    console.log(arr)
    const mt = new MergeTrick(arr)
    const naive = new _MergeTrickNaive(arr)
    for (let i = 0; i < 4; i++) {
      let start = Math.floor(Math.random() * arr.length)
      let end = Math.floor(Math.random() * arr.length)
      if (start > end) [start, end] = [end, start]
      const delta = Math.floor(Math.random() * 4)
      mt.add(start, end, delta)
      naive.add(start, end, delta)
      console.log('add', start, end, delta)
      const res1 = mt.sortedNums
      const res2 = naive.sorted()

      for (let i = 0; i < res1.length; i++) {
        if (res1[i] !== res2[i]) {
          console.log('not match')
          console.log(res1)
          console.log(res2)
          return
        }
      }
    }
  }
}
