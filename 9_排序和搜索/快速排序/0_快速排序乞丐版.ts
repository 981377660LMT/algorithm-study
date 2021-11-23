const qucikSort = (arr: readonly number[]): readonly number[] => {
  if (arr.length <= 1) return arr
  // 最基础的partition
  const start = arr[0]
  const small = qucikSort(arr.slice(1).filter(ele => ele < start))
  const big = qucikSort(arr.slice(1).filter(ele => ele >= start))
  return [...small, start, ...big]
}

const arr = [4, 1, 2, 5, 3, 6, 7]
console.log(qucikSort(arr))

export {}
