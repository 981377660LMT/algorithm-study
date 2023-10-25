/**
 * 可取重复元素的组合，遍历所有大小为`r`的组合.
 * @param arr 待遍历的数组.
 * @param r 组合大小.
 * @param cb 回调函数, 用于处理每个组合的结果.返回`true`时停止遍历.
 * @param copy 是否浅拷贝每个组合的结果, 默认为`false`.
 * @example
 * ```ts
 * enumerateCombinationsWithReplacement([1, 2, 3], 2, comb => {
 *   console.log(comb)
 * })
 * // [ 1, 1 ]
 * // [ 1, 2 ]
 * // [ 1, 3 ]
 * // [ 2, 2 ]
 * // [ 2, 3 ]
 * // [ 3, 3 ]
 * ```
 * @complexity 2e7 => 205.486ms
 */
function enumerateCombinationsWithReplacement<C>(arr: ArrayLike<C>, r: number, cb: (comb: readonly C[]) => boolean | void, copy = false): void {
  bt(0, [])

  function bt(pos: number, path: C[]): boolean {
    if (path.length === r) {
      return !!cb(copy ? path.slice() : path)
    }

    for (let i = pos; i < arr.length; i++) {
      path.push(arr[i])
      if (bt(i, path)) return true // 可取重复的元素
      path.pop()
    }
    return false
  }
}

/**
 * 可取重复元素的组合.
 * 模拟python的`itertools.combinations_with_replacement`.
 * @complexity 2e6 => 809.43ms
 * @deprecated
 */
function* combinationsWithReplacement<C>(arr: ArrayLike<C>, r: number): Generator<C[]> {
  yield* bt(0, [])

  function* bt(pos: number, path: C[]): Generator<C[]> {
    if (path.length === r) {
      yield path.slice()
      return
    }

    for (let i = pos; i < arr.length; i++) {
      path.push(arr[i])
      yield* bt(i, path) // 可取重复的元素
      path.pop()
    }
  }
}

if (require.main === module) {
  const n = 20
  const r = 10
  const arr = Array.from({ length: n }, (_, i) => i)

  console.time('combinations')
  let count = 0
  enumerateCombinationsWithReplacement(arr, r, comb => {
    count++
  })
  console.timeEnd('combinations') // !205.486ms
  console.log(count) // !20030010
}

export { enumerateCombinationsWithReplacement }
