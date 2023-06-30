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
): void {}

export { sortRange }
