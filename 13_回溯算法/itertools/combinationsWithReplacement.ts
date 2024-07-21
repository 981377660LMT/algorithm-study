/* eslint-disable no-constant-condition */

/**
 * 可取重复元素的组合，遍历所有大小为`r`的组合.
 * @param n 元素总数.
 * @param r 选取的元素个数.
 * @param f 回调函数, 用于处理每个组合的结果.返回`true`时停止遍历.
 * @example
 * ```ts
 * enumerateCombinationsWithReplacement(3, 2, indices => {
 *   console.log(indices)
 * })
 * // [ 0, 0 ]
 * // [ 0, 1 ]
 * // [ 0, 2 ]
 * // [ 1, 1 ]
 * // [ 1, 2 ]
 * // [ 2, 2 ]
 * ```
 * @complexity 2e7 => 100ms.
 */
function enumerateCombinationsWithReplacement<C>(
  n: number,
  r: number,
  f: (indicesView: readonly number[]) => boolean | void
): void {
  const ids = Array(r).fill(0)
  if (f(ids)) {
    return
  }
  while (true) {
    let i = r - 1
    for (; i >= 0; i--) {
      if (ids[i] !== n - 1) {
        break
      }
    }
    if (i === -1) {
      return
    }
    ids[i]++
    for (let j = i + 1; j < r; j++) {
      ids[j] = ids[i]
    }
    if (f(ids)) {
      return
    }
  }
}

if (require.main === module) {
  const n = 20
  const r = 10

  console.time('combinations')
  let count = 0
  enumerateCombinationsWithReplacement(n, r, indices => {
    count++
  })
  console.timeEnd('combinations') // !100ms
  console.log(count) // !20030010
}

export { enumerateCombinationsWithReplacement }
