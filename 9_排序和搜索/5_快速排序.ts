const fastSort = (arr: number[]): number[] => {
  if (arr.length <= 1) return arr
  const start = arr[0]
  const small = fastSort(arr.slice(1).filter(ele => ele < start))
  const big = fastSort(arr.slice(1).filter(ele => ele >= start))
  return [...small, start, ...big]
}

const arr = [4, 1, 2, 5, 3, 6, 7]
console.log(fastSort(arr))

export {}
