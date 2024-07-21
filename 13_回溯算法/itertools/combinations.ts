/* eslint-disable no-constant-condition */

/**
 * 遍历所有大小为`r`的组合.
 * @param n 元素总数.
 * @param r 选取的元素个数.
 * @param f 回调函数, 用于处理每个组合的结果.返回`true`时停止遍历.
 * @example
 * ```ts
 * enumerateCombinations(4, 2, comb => {
 *   console.log(comb)
 * })
 * // [ 0, 1 ]
 * // [ 0, 2 ]
 * // [ 0, 3 ]
 * // [ 1, 2 ]
 * // [ 1, 3 ]
 * // [ 2, 3 ]
 * ```
 * @complexity C(30,10)(3e7) => 170ms.
 */
function enumerateCombinations(
  n: number,
  r: number,
  f: (indicesView: readonly number[]) => boolean | void
): void {
  const ids = Array.from({ length: r }, (_, i) => i)
  if (f(ids)) {
    return
  }
  while (true) {
    let i = r - 1
    for (; i >= 0; i--) {
      if (ids[i] !== i + n - r) {
        break
      }
    }
    if (i === -1) {
      return
    }
    ids[i]++
    for (let j = i + 1; j < r; j++) {
      ids[j] = ids[j - 1] + 1
    }
    if (f(ids)) {
      return
    }
  }
}

export { enumerateCombinations }

if (require.main === module) {
  enumerateCombinations(4, 2, comb => {
    console.log(comb)
  })

  const n = 30
  const r = 10
  console.time('enumerateCombinations')
  let count = 0
  enumerateCombinations(n, r, () => {
    count++
  })
  console.log(count) // !30045015
  console.timeEnd('enumerateCombinations') // !170ms
}
