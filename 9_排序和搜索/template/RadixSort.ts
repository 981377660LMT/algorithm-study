// radixSort 基数排序

/**
 * 基数排序.
 * 按照低位先排序，然后收集；再按照高位排序，然后再收集；依次类推，直到最高位.
 * @param arr 待排序数组.
 * @param maxDigit 最大位数.
 * @param minValue 每个位数的最小值.
 * @param maxValue 每个位数的最大值.
 * @param getValue 获取元素某一位的值.
 * @complexity `O(n*maxDigit)`
 */
function radixSort<T>(
  arr: { [index: number]: T; readonly length: number },
  maxDigit: number,
  minValue: number,
  maxValue: number,
  getValue: (item: T, digit: number) => number
): void {
  const buckets = Array(maxValue - minValue + 1)

  for (let i = 0; i < maxDigit; i++) {
    resetBuckets()
    for (let j = 0; j < arr.length; j++) {
      const value = getValue(arr[j], i)
      buckets[value - minValue].push(arr[j])
    }

    for (let j = 0, ptr = 0; j < buckets.length; j++) {
      const bucket = buckets[j]
      for (let k = 0; k < bucket.length; k++) {
        arr[ptr++] = bucket[k]
      }
    }
  }

  function resetBuckets(): void {
    for (let i = 0; i < buckets.length; i++) buckets[i] = []
  }
}

export {}

if (require.main === module) {
  const array = [151, 3, 44, 38, 5, 47, 15, 36, 26, 27, 2, 46, 4, 19, 50, 48]
  const maxDigit = Math.max(...array.map(item => item.toString().length))
  const minValue = 0
  const maxValue = 9
  const getValue = (item: number, digit: number) => {
    const offset = 10 ** digit
    return ~~(item / offset) % 10
  }

  radixSort(array, maxDigit, minValue, maxValue, getValue)
  console.log(array)
}
