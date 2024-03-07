/**
 * 比较两个类数组的大小.
 *
 * @param arr1 类数组1.
 * @param arr2 类数组2.
 * @param compareFn 比较类数组元素的函数.默认使用大于、等于比较.
 */
function compareArray<T>(arr1: ArrayLike<T>, arr2: ArrayLike<T>, compareFn?: (a: T, b: T) => -1 | 0 | 1): -1 | 0 | 1 {
  if (!compareFn) {
    compareFn = (a: T, b: T): -1 | 0 | 1 => {
      if (a > b) {
        return 1
      } else if (a === b) {
        return 0
      } else {
        return -1
      }
    }
  }

  const n1 = arr1.length
  const n2 = arr2.length
  for (let i = 0; i < n1; i++) {
    if (i >= n2) return 1
    const res = compareFn(arr1[i], arr2[i])
    if (res !== 0) return res
  }
  return n1 < n2 ? -1 : 0
}

export { compareArray }

if (require.main === module) {
  console.log(compareArray([1, 2, 3], [1, 2, 3])) // 0
  console.log(compareArray([1, 2, 3], [1, 2, 4])) // -1
  console.log(compareArray([1, 2, 3], [1, 2, 2])) // 1
  console.log(compareArray([1, 2, 3], [1, 2, 3, 4])) // -1
  console.log(compareArray([1, 2, 3, 4], [1, 2, 3])) // 1
  console.log(compareArray([1, 2, 3], [1, 2])) // 1
  console.log(compareArray([1, 2], [1, 2, 3])) // -1
  console.log(compareArray([1, 2], [1, 2, 2])) // -1
}
