/* eslint-disable no-param-reassign */

/**
 * 遍历多个类数组对象的笛卡尔积.
 * @param arrs `selects`中的每个类数组对象都是一个可选项列表.
 * @param f 回调函数, 用于处理每个笛卡尔积的结果.返回`true`时停止遍历.
 * @example
 * ```ts
 * enumerateProduct([['A', 'a'], ['1'], ['B', 'b'], ['2']], select => {
 *  console.log(select)
 * })
 * // [ 'A', '1', 'B', '2' ]
 * // [ 'A', '1', 'b', '2' ]
 * // [ 'a', '1', 'B', '2' ]
 * // [ 'a', '1', 'b', '2' ]
 * ```
 * @complexity 11!(4e7) => 383.51ms
 */
function enumerateProduct<S>(
  arrs: ArrayLike<S>[],
  f: (groupView: readonly S[]) => boolean | void
): void {
  const n = arrs.length
  const bt = (pos: number, group: S[], ptr: number): boolean => {
    if (pos === n) {
      return !!f(group)
    }
    const arr = arrs[pos]
    for (let i = 0; i < arr.length; i++) {
      group[ptr] = arr[i]
      if (bt(pos + 1, group, ptr + 1)) return true
    }
    return false
  }
  bt(0, Array(n), 0)
}

/**
 * 模拟python的`itertools.product`.
 * @complexity 11!(4e7) => 17.461s
 * @deprecated
 */
function* product<S>(...arrs: ArrayLike<S>[]): Generator<S[]> {
  const n = arrs.length
  yield* bt(0, [])

  function* bt(pos: number, path: S[]): Generator<S[]> {
    if (pos === n) {
      yield path.slice()
      return
    }

    const arr = arrs[pos]
    for (let i = 0; i < arr.length; i++) {
      path.push(arr[i])
      yield* bt(pos + 1, path)
      path.pop()
    }
  }
}

export { enumerateProduct }

if (require.main === module) {
  enumerateProduct([['A', 'a'], ['1'], ['B', 'b'], ['2']], select => {
    console.log(select)
  })

  // performance test
  const n = 11
  const facs = Array.from({ length: n }, (_, i) => Array.from({ length: i + 1 }, (_, j) => j + 1))
  console.time('enumerateProduct')
  let count = 0
  enumerateProduct(facs, select => {
    count++
  })
  console.timeEnd('enumerateProduct') // !383.51ms
  console.log(count) // !39916800

  count = 0
  console.time('product')
  for (const _ of product(...facs)) {
    count++
  }
  console.timeEnd('product') // !17.461s
}
