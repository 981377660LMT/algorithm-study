/**
 * 遍历所有大小为`r`的组合.
 * @param arr 待遍历的数组.
 * @param r 组合大小.
 * @param cb 回调函数, 用于处理每个组合的结果.
 * @param copy 是否浅拷贝每个组合的结果, 默认为`false`.
 * @example
 * ```ts
 * enumerateCombinations([1, 2, 3, 4], 2, comb => {
 *   console.log(comb)
 * })
 * // [ 1, 2 ]
 * // [ 1, 3 ]
 * // [ 1, 4 ]
 * // [ 2, 3 ]
 * // [ 2, 4 ]
 * // [ 3, 4 ]
 * ```
 * @complexity C(30,10)(3e7) => 347.592ms
 */
function enumerateCombinations<C>(
  arr: ArrayLike<C>,
  r: number,
  cb: (comb: readonly C[]) => void,
  copy = false
): void {
  bt(0, [])

  function bt(pos: number, path: C[]): void {
    if (path.length === r) {
      cb(copy ? path.slice() : path)
      return
    }

    for (let i = pos; i < arr.length; i++) {
      path.push(arr[i])
      bt(i + 1, path) // 不可取重复的元素
      path.pop()
    }
  }
}

/**
 * 模拟python的`itertools.combinations`.
 * @complexity C(24,10)(2e6) => 946.732ms
 */
function* combinations<C>(arr: ArrayLike<C>, r: number): Generator<C[]> {
  yield* bt(0, [])

  function* bt(pos: number, path: C[]): Generator<C[]> {
    if (path.length === r) {
      yield path.slice()
      return
    }

    for (let i = pos; i < arr.length; i++) {
      path.push(arr[i])
      yield* bt(i + 1, path) // 不可取重复的元素
      path.pop()
    }
  }
}

export { enumerateCombinations, combinations }

if (require.main === module) {
  enumerateCombinations([1, 2, 3, 4], 2, comb => {
    console.log(comb)
  })

  const n = 30
  const r = 10
  const arr = Array.from({ length: n }, (_, i) => i)

  console.time('enumerateCombinations')
  let count = 0
  enumerateCombinations(arr, r, () => {
    count++
  })
  console.log(count) // !1961256
  console.timeEnd('enumerateCombinations') // !34.127ms

  // console.time('combinations')
  // let count = 0
  // for (const _ of combinations(arr, r)) {
  //   count++
  // }
  // console.timeEnd('combinations') // !910.397ms
  // console.log(count) // !1961256
}
