/* eslint-disable no-inner-declarations */
// 用view去sort

/**
 * 区间原地排序.
 * @param keys 要排序的元素的key,必须是uint32.
 * @param compareFn 比较函数, 用于比较两个键对应的元素的大小.
 * @param start 要排序的起始位置(包含).
 * @param end 要排序的结束位置(不包含).
 */
function sortRangeStable(
  keys: Uint32Array,
  compareFn: (key1: number, key2: number) => number,
  start = 0,
  end = keys.length
): void {
  if (start < 0) start = 0
  if (end > keys.length) end = keys.length
  if (start >= end) return
  keys.subarray(start, end).sort(compareFn)
}

if (require.main === module) {
  demo()
  performDemo()

  function demo() {
    const people = [
      { name: 'a', age: 1 },
      { name: 'b', age: 2 },
      { name: 'c', age: 3 },
      { name: 'd', age: 4 },
      { name: 'e', age: 5 },
      { name: 'f', age: 6 },
      { name: 'g', age: 7 },
      { name: 'h', age: 8 },
      { name: 'i', age: 9 }
    ]
    const ids = new Uint32Array([1, 2, 3, 4, 5, 6, 7, 8, 9])
    sortRangeStable(ids, (id1, id2) => id1 - id2, 1, 8)
  }

  function performDemo() {
    const n = 1e6
    const ids = new Uint32Array(n)
    for (let i = 0; i < n; i++) {
      ids[i] = Math.floor(Math.random() * n)
    }
    console.time('sortRangeStable')
    sortRangeStable(ids, (id1, id2) => id1 - id2, 0, n)
    console.timeEnd('sortRangeStable')
  }
}

export {}
