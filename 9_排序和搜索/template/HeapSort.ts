// HeapSort 堆排序
// 防止快速排序退化时，可以使用堆排序

type Mutable<T> = { -readonly [P in keyof T]: T[P] }

function heapSort<T>(
  arr: Mutable<ArrayLike<T>>,
  compareFn: (a: T, b: T) => number,
  start = 0,
  end = arr.length
): void {
  if (start < 0) start = 0
  if (end > arr.length) end = arr.length
  if (start >= end) return

  const first = start
  const lo = 0
  const hi = end - start

  // Build heap with greatest element at top.
  for (let i = (hi - 1) >>> 1; ~i; i--) {
    pushDown(i, hi, first)
  }

  // Pop elements, largest first, into end of data.
  for (let i = hi - 1; ~i; i--) {
    const tmp = arr[first]
    arr[first] = arr[first + i]
    arr[first + i] = tmp
    pushDown(lo, i, first)
  }

  // pushDown implements the heap property on arr[lo:hi].
  // first is an offset into the array where the root of the heap lies.
  function pushDown(lo: number, hi: number, offset: number): void {
    let root = lo
    while (true) {
      let child = (root << 1) | 1
      if (child >= hi) {
        break
      }
      if (child + 1 < hi && compareFn(arr[offset + child], arr[offset + child + 1]) < 0) {
        child++
      }
      if (compareFn(arr[offset + root], arr[offset + child]) >= 0) {
        return
      }

      const tmp = arr[offset + root]
      arr[offset + root] = arr[offset + child]
      arr[offset + child] = tmp
      root = child
    }
  }
}

export {}

if (require.main === module) {
  const arr = [4, 2, 100, 99, 10000, -1, 99, 2]
  heapSort(arr, (a, b) => a - b)
  console.log(arr)

  // https://leetcode.cn/problems/the-k-strongest-values-in-an-array/
  // eslint-disable-next-line no-inner-declarations
  function getStrongest(arr: number[], k: number): number[] {
    arr.sort((a, b) => a - b)
    heapSort(arr, (a, b) => a - b)
    const m = arr[(arr.length - 1) >> 1]

    heapSort(arr, (a, b) => {
      const diffA = Math.abs(a - m)
      const diffB = Math.abs(b - m)
      if (diffA === diffB) {
        return b - a
      }
      return diffB - diffA
    })

    return arr.slice(0, k)
  }
}
