/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable no-console */

type NestedArray<T> = Array<T | NestedArray<T>>

/**
 * dfs实现 flat 函数.
 *
 * 更严格的返回类型参考 {@link Array.prototype.flat}.
 */
function flat1<T>(arr: NestedArray<T>, depth = 1): (T | NestedArray<T>)[] {
  const res = []

  for (let i = 0; i < arr.length; i++) {
    const item = arr[i]
    if (Array.isArray(item) && depth > 0) {
      res.push(...flat1(item, depth - 1))
    } else {
      res.push(item)
    }
  }

  return res
}

const nestedArray = ['a', ['v', 'g'], [['a', ['f'], 'v']]]
console.log(flat1(nestedArray, Infinity)) // [ 'a', 'v', 'g', 'a', 'f', 'v' ]

type Item<T> = T | NestedArray<T>

/**
 * bfs实现 flat 函数.
 *
 * 更严格的返回类型参考 {@link Array.prototype.flat}.
 */
function flat2<T>(arr: NestedArray<T>, depth = 1): Item<T>[] {
  const res: Item<T>[] = []
  const queue: [item: Item<T>, depth: number][] = arr.map(item => [item, depth])

  while (queue.length) {
    const [cur, curDepth] = queue.shift()!
    if (Array.isArray(cur) && curDepth > 0) {
      queue.push(...cur.map<[item: Item<T>, depth: number]>(item => [item, curDepth - 1]))
    } else {
      res.push(cur)
    }
  }

  return res
}

export {}

if (require.main === module) {
  const arr = [1, [2], [[[7]]], [3, [4, 9]]]
  console.log(flat2(arr))
  console.log(flat2(arr, 1))
  console.log(flat2(arr, 2))
}
