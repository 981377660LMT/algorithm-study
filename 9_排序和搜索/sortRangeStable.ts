/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */

// sortRangeStable/rangeSortStable

// 1. 使用高度优化的排序算法
// !不稳定的部分排序pdqsort:
//  https://cs.opensource.google/go/go/+/refs/tags/go1.20.5:src/sort/zsortfunc.go;l=61
// !稳定的部分排序:
//  https://cs.opensource.google/go/go/+/refs/tags/go1.20.5:src/sort/zsortfunc.go;l=335
//
// 2. 使用类型数组加速
// !存id，然后用类型数组subarray对子数组排序，fill更新区间
// 3. 为什么不slice再排序再赋值回去呢?

/**
 * 稳定的部分排序.
 * @see {@link https://cs.opensource.google/go/go/+/refs/tags/go1.20.5:src/sort/slice.go;l=35}
 *      {@link https://cs.opensource.google/go/go/+/refs/tags/go1.20.5:src/sort/zsortfunc.go;l=335}
 */
function sortRangeStable<V>(
  arr: V[],
  compareFn: (a: V, b: V) => number,
  start = 0,
  end = arr.length
): void {
  if (start < 0) start = 0
  if (end > arr.length) end = arr.length
  if (start >= end) return

  stableSort(start, end)

  /**
   * 1. 分块排序，每块大小为 20(待调节)，每个块内部使用插入排序.
   * 2. 循环合并相邻的两个 block，每次循环 blockSize 扩大一倍，直到 blockSize > n 为止.
   *
   * TODO 调节分块大小(自适应)
   */
  function stableSort(start: number, end: number): void {
    let blockSize = 20 // must be > 0
    const n = end - start
    let a = start
    let b = a + blockSize
    while (b <= n) {
      insertionSort(a, b)
      a = b
      b += blockSize
    }
    insertionSort(a, end)

    while (blockSize < n) {
      let a = start
      let b = a + blockSize * 2
      while (b <= n) {
        symMerge(a, a + blockSize, b)
        a = b
        b += blockSize * 2
      }
      const mid = a + blockSize
      if (mid < end) {
        symMerge(a, mid, end)
      }
      blockSize *= 2
    }
  }

  /**
   * 对切片 arr[start:end] 进行插入排序.
   * 适用于短数组.
   */
  function insertionSort(start: number, end: number): void {
    for (let i = start + 1; i < end; i++) {
      for (let j = i; j > start && compareFn(arr[j - 1], arr[j]) > 0; j--) {
        swap(j - 1, j)
      }
    }
  }

  /**
   * 稳定合并两个有序数组 `arr[a:m]` 和 `arr[m:b]`.
   * @see {@link https://link.springer.com/chapter/10.1007/978-3-540-30140-0_63}
   */
  function symMerge(a: number, m: number, b: number): void {
    if (m - a === 1) {
      let i = m
      let j = b
      while (i < j) {
        const h = (i + j) >>> 1
        if (compareFn(arr[h], arr[a]) < 0) {
          i = h + 1
        } else {
          j = h
        }
      }

      for (let k = a; k < i - 1; k++) {
        swap(k, k + 1)
      }
      return
    }

    if (b - m === 1) {
      let i = a
      let j = m
      while (i < j) {
        const h = (i + j) >>> 1
        if (compareFn(arr[m], arr[h]) < 0) {
          j = h
        } else {
          i = h + 1
        }
      }

      for (let k = m; k > i; k--) {
        swap(k, k - 1)
      }
      return
    }

    const mid = (a + b) >>> 1
    const n = mid + m
    let start = a
    let r = m
    if (m > mid) {
      start = n - b
      r = mid
    }
    const p = n - 1

    while (start < r) {
      const c = (start + r) >>> 1
      if (compareFn(arr[p - c], arr[c]) < 0) {
        r = c
      } else {
        start = c + 1
      }
    }

    const end = n - start
    if (start < m && m < end) {
      rotate(start, m, end)
    }
    if (a < start && start < mid) {
      symMerge(a, start, mid)
    }
    if (mid < end && end < b) {
      symMerge(mid, end, b)
    }
  }

  function swap(i: number, j: number): void {
    const tmp = arr[i]
    arr[i] = arr[j]
    arr[j] = tmp
  }

  function swapRange(a: number, b: number, n: number): void {
    for (let i = 0; i < n; i++) {
      swap(a + i, b + i)
    }
  }

  /**
   * 块的循环旋转操作，使得指定范围内的元素按照特定顺序进行重新排列。
   */
  function rotate(a: number, m: number, b: number): void {
    let i = m - a
    let j = b - m
    while (i !== j) {
      if (i > j) {
        swapRange(m - i, m, j)
        i -= j
      } else {
        swapRange(m - i, m + j - i, i)
        j -= i
      }
    }
    swapRange(m - i, m, i) // i === j
  }
}

export { sortRangeStable }

if (require.main === module) {
  demo()
  testPerform()
  stableDemo()

  function demo(): void {
    const arr = [1, 2, 5, 4, 6, -1, 7, 1, 1]
    sortRangeStable(arr, (a, b) => a - b, 2, 6)
    console.log(arr)
  }

  function testPerform(): void {
    const n = 1e6
    const arr = Array(n)
    for (let i = 0; i < n; i++) {
      arr[i] = (Math.random() * 1e9) | 0
    }

    const copy1 = arr.slice()
    console.time('sortRangeStable')
    sortRangeStable(copy1, (a, b) => a - b)
    console.timeEnd('sortRangeStable')

    const copy2 = arr.slice()
    console.time('sort')
    copy2.sort((a, b) => a - b)
    console.timeEnd('sort')
  }

  function stableDemo(): void {
    type V = { key: number; value: string }
    const arr: V[] = Array(10)
    for (let i = 0; i < arr.length; i++) {
      arr[i] = { key: i % 3, value: `${i}` }
    }
    console.log(arr)

    sortRangeStable(arr, (a, b) => a.key - b.key)
    console.log(arr)
  }
}
