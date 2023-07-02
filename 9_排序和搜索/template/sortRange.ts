/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */

// sortRange/rangeSort

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
 * 不稳定的部分排序.
 * @see {@link https://cs.opensource.google/go/go/+/refs/tags/go1.20.5:src/sort/slice.go;l=21}
 *      {@link https://cs.opensource.google/go/go/+/refs/tags/go1.20.5:src/sort/zsortfunc.go;l=61}
 */
function sortRange<V>(
  arr: V[],
  compareFn: (a: V, b: V) => number,
  start = 0,
  end = arr.length
): void {
  if (start < 0) start = 0
  if (end > arr.length) end = arr.length
  if (start >= end) return

  // hint for pdqsort when choosing the pivot
  const UNKNOWN_HINT = 0
  const INCREASING_HINT = 1
  const DECREASING_HINT = 2

  const len = arr.length
  const limit = 32 - Math.clz32(len)
  pdqSort(start, end, limit)

  function pdqSort(a: number, b: number, limit: number): void {
    const MAX_INSERTION = 12
    let wasBalanced = true // whether the last partitioning was reasonably balanced
    let wasPartitioned = false // whether the slice was already partitioned

    while (true) {
      const length = b - a
      if (length <= MAX_INSERTION) {
        insertionSortFunc(a, b)
        return
      }

      // Fall back to heapsort if too many bad choices were made.
      if (!limit) {
        heapSortFunc(a, b)
        return
      }

      // If the last partitioning was imbalanced, we need to breaking patterns.
      if (!wasBalanced) {
        breakPatterns(a, b)
        limit--
      }

      const pivotPair = choosePivot(a, b)
      let pivot = pivotPair[0]
      let hint = pivotPair[1]
      if (hint === DECREASING_HINT) {
        reverseRange(a, b)
        // The chosen pivot was pivot-a elements after the start of the array.
        // After reversing it is pivot-a elements before the end of the array.
        // The idea came from Rust's implementation.
        pivot = b - 1 - (pivot - a)
        hint = INCREASING_HINT
      }

      // The slice is likely already sorted.
      if (wasBalanced && wasPartitioned && hint === INCREASING_HINT) {
        if (partialInsertionSort(a, b)) {
          return
        }
      }

      // Probably the slice contains many duplicate elements, partition the slice into
      // elements equal to and elements greater than the pivot.
      if (a > 0 && !(compareFn(arr[a - 1], arr[pivot]) < 0)) {
        const mid = partitionEqual(a, b, pivot)
        a = mid
        continue
      }

      const partitionPair = partition(a, b, pivot)
      const mid = partitionPair[0]
      const alreadyPartitioned = partitionPair[1]
      wasPartitioned = alreadyPartitioned

      const leftLen = mid - a
      const rightLen = b - mid
      const balanceThreshold = length >>> 3
      if (leftLen < rightLen) {
        wasBalanced = leftLen >= balanceThreshold
        pdqSort(a, mid, limit)
        a = mid + 1
      } else {
        wasBalanced = rightLen >= balanceThreshold
        pdqSort(mid + 1, b, limit)
        b = mid
      }
    }
  }

  function insertionSortFunc(start: number, end: number): void {
    for (let i = start + 1; i < end; i++) {
      for (let j = i; j > start && compareFn(arr[j - 1], arr[j]) > 0; j--) {
        swap(j - 1, j)
      }
    }
  }

  // #region heapSort
  function heapSortFunc(a: number, b: number): void {
    const first = a
    let lo = 0
    let hi = b - a

    // Build heap with greatest element at top.
    for (let i = (hi - 1) >>> 1; ~i; i--) {
      siftDownFunc(i, hi, first)
    }

    // Pop elements, largest first, into end of data.
    for (let i = hi - 1; ~i; i--) {
      swap(first, first + i)
      siftDownFunc(lo, i, first)
    }
  }

  // siftDown_func implements the heap property on arr[lo:hi].
  // first is an offset into the array where the root of the heap lies.
  function siftDownFunc(lo: number, hi: number, first: number): void {
    let root = lo
    while (true) {
      let child = (root << 1) | 1
      if (child >= hi) {
        return
      }
      if (child + 1 < hi && compareFn(arr[first + child], arr[first + child + 1]) < 0) {
        child++
      }
      if (!(compareFn(arr[first + root], arr[first + child]) < 0)) {
        return
      }
      swap(first + root, first + child)
      root = child
    }
  }
  // #endregion heapSort

  // #region quickSort
  // breakPatterns_func scatters some elements around in an attempt to break some patterns
  // that might cause imbalanced partitions in quicksort.
  function breakPatterns(a: number, b: number): void {
    const length = b - a
    if (length >= 8) {
      let random = length
      const modulus = nextPowerOf2(length)
      for (let idx = a + (((length >> 2) << 1) - 1); idx <= a + (((length >> 2) << 1) | 1); idx++) {
        random = nextXorShift(random)
        let other = (random & (modulus - 1)) >>> 0
        if (other >= length) {
          other -= length
        }
        swap(idx, a + other)
      }
    }
  }

  function choosePivot(a: number, b: number): [pivot: number, hint: number] {
    const SHORTEST_NINTHER = 50
    const MAX_SWAPS = 4 * 3

    const l = b - a
    let swaps = 0
    const step = l >> 2
    let i = a + step
    let j = i + step
    let k = j + step

    if (l >= 8) {
      if (l >= SHORTEST_NINTHER) {
        // Tukey ninther method, the idea came from Rust's implementation.
        i = medianAdjacent(i)
        j = medianAdjacent(j)
        k = medianAdjacent(k)
      }
      // Find the median among i, j, k and stores it into j.
      j = median(i, j, k)
    }

    switch (swaps) {
      case 0:
        return [j, INCREASING_HINT]
      case MAX_SWAPS:
        return [j, DECREASING_HINT]
      default:
        return [j, UNKNOWN_HINT]
    }

    // finds the median of data[a - 1], data[a], data[a + 1] and stores the index into a.
    function medianAdjacent(i: number): number {
      return median(i - 1, i, i + 1)
    }

    // returns x where data[x] is the median of data[a],data[b],data[c], where x is a, b, or c.
    function median(a: number, b: number, c: number): number {
      if (b < a) {
        swaps++
        a ^= b
        b ^= a
        a ^= b
      }
      if (c < b) {
        swaps++
        b ^= c
        c ^= b
        b ^= c
      }
      if (b < a) {
        swaps++
        a ^= b
        b ^= a
        a ^= b
      }
      return b
    }
  }

  function partition(
    a: number,
    b: number,
    pivot: number
  ): [newPivot: number, alreadyPartitioned: boolean] {
    swap(a, pivot)
    let i = a + 1
    let j = b - 1
    while (i <= j && compareFn(arr[i], arr[a]) < 0) {
      i++
    }
    while (i <= j && !(compareFn(arr[j], arr[a]) < 0)) {
      j--
    }
    if (i > j) {
      swap(j, a)
      return [j, true]
    }
    swap(i, j)
    i++
    j--
    while (true) {
      while (i <= j && compareFn(arr[i], arr[a]) < 0) {
        i++
      }
      while (i <= j && !(compareFn(arr[j], arr[a]) < 0)) {
        j--
      }
      if (i > j) {
        break
      }
      swap(i, j)
      i++
      j--
    }
    swap(j, a)
    return [j, false]
  }

  function partitionEqual(a: number, b: number, pivot: number): number {
    swap(a, pivot)
    // i and j are inclusive of the elements remaining to be partitioned
    let i = a + 1
    let j = b - 1
    while (true) {
      while (i <= j && !(compareFn(arr[a], arr[i]) < 0)) {
        i++
      }
      while (i <= j && compareFn(arr[a], arr[j]) < 0) {
        j--
      }
      if (i > j) {
        break
      }
      swap(i, j)
      i++
      j--
    }
    return i
  }
  // #region quickSort

  function partialInsertionSort(a: number, b: number): boolean {
    const MAX_STEPS = 5 // maximum number of adjacent out-of-order pairs that will get shifted
    const SHORTEST_SHIFTING = 50 // don't shift any elements on short arrays

    let i = a + 1
    for (let j = 0; j < MAX_STEPS; j++) {
      while (i < b && compareFn(arr[i], arr[i - 1]) >= 0) {
        i++
      }

      if (i === b) {
        return true
      }

      if (b - a < SHORTEST_SHIFTING) {
        return false
      }

      swap(i, i - 1)

      // Shift the smaller one to the left.
      if (i - a >= 2) {
        for (let j = i - 1; j >= 1; j--) {
          if (compareFn(arr[j], arr[j - 1]) >= 0) {
            break
          }
          swap(j, j - 1)
        }
      }
      // Shift the greater one to the right.
      if (b - i >= 2) {
        for (let j = i + 1; j < b; j++) {
          if (compareFn(arr[j], arr[j - 1]) >= 0) {
            break
          }
          swap(j, j - 1)
        }
      }
    }

    return false
  }

  // #region utils
  function swap(i: number, j: number): void {
    const tmp = arr[i]
    arr[i] = arr[j]
    arr[j] = tmp
  }

  function reverseRange(a: number, b: number): void {
    for (let i = a, j = b - 1; i < j; i++, j--) {
      swap(i, j)
    }
  }

  function nextXorShift(xorShift: number): number {
    xorShift ^= xorShift << 13
    xorShift ^= xorShift >> 17
    xorShift ^= xorShift << 5
    return xorShift >>> 0
  }

  function nextPowerOf2(length: number): number {
    return 1 << (32 - Math.clz32(length))
  }
  // #endregion utils
}

export { sortRange }

if (require.main === module) {
  demo()
  testPerform()

  function demo(): void {
    const arr = [1, 2, 5, 4, 6, -1, 7, 1, 1]
    sortRange(arr, (a, b) => a - b, 2, 6)
    console.log(arr)
  }

  function testPerform(): void {
    const n = 1e6
    const arr = Array(n)
    for (let i = 0; i < n; i++) {
      arr[i] = (Math.random() * 1e9) | 0
    }

    const copy1 = arr.slice()
    console.time('sortRange')
    sortRange(copy1, (a, b) => a - b)
    console.timeEnd('sortRange')

    const copy2 = arr.slice()
    console.time('sort')
    copy2.sort((a, b) => a - b)
    console.timeEnd('sort')
  }
}
