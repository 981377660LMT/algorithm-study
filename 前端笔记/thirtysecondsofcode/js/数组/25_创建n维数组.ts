// 创建具有给定值的 n 维数组。
// Create a n-dimensional array with given value.
initializeNDArray(1, 3) // [1, 1, 1]
console.log(initializeNDArray(5, 2, 2, 2)) // [[[5, 5], [5, 5]], [[5, 5], [5, 5]]]

/**
 * @description n维数组
 */
type NDArray<T> = Array<NDArray<T> | T> | T

function initializeNDArray(initialValue: number, ...sizeOfDimensions: number[]): NDArray<number> {
  return sizeOfDimensions.length === 0
    ? initialValue
    : Array.from({ length: sizeOfDimensions[0] }, () =>
        initializeNDArray(initialValue, ...sizeOfDimensions.slice(1))
      )
}
