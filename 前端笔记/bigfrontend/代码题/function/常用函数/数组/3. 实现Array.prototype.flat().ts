/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable no-console */

type MultiDimensionalArray<T> = (T | MultiDimensionalArray<T>)[]

function flat<T>(arr: MultiDimensionalArray<T>, n: number): MultiDimensionalArray<T> {
  return flatten(arr, n)

  function flatten(item: MultiDimensionalArray<T>, todo: number): MultiDimensionalArray<T> {
    if (todo <= 0) return item

    const res: MultiDimensionalArray<T> = []
    item.forEach(v => {
      if (Array.isArray(v)) {
        res.push(...flatten(v, todo - 1))
      } else {
        res.push(v)
      }
    })

    return res
  }
}

export {}

// !注意只能dfs实现,不能用bfs实现

if (require.main === module) {
  const nestedArray = ['a', ['v', 'g'], [['a', ['f'], 'v']]]
  console.log(flat(nestedArray, 2)) // [ 'a', 'v', 'g', 'a', [ 'f' ], 'v' ]
  console.log(flat(nestedArray, Infinity)) // [ 'a', 'v', 'g', 'a', 'f', 'v' ]
}
