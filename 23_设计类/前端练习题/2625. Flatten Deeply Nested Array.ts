// 扁平化数组

// !伪代码
// private flatten(nestedList: NestedInteger[]) {
//   for (const nestedInteger of nestedList) {
//     if (nestedInteger.isInteger()) {
//       this.list.push(nestedInteger)
//     } else {
//       this.flatten(nestedInteger.getList())
//     }
//   }
// }

// NestedArray
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

if (require.main === module) {
  // arr = [1, 2, 3, [4, 5, 6], [7, 8, [9, 10, 11], 12], [13, 14, 15]]
  // [1, 2, 3, 4, 5, 6, 7, 8, [9, 10, 11], 12, 13, 14, 15]
  console.log(flat([1, 2, 3, [4, 5, 6], [7, 8, [9, 10, 11], 12], [13, 14, 15]], 1))
}
