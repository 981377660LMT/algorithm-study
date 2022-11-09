const odd = [1, 2, 3, 4, 5]
const even = [1, 2, 3, 4, 5, 6]

// !求中位数
function calMid(arr: number[]): number {
  arr = arr.slice().sort((a, b) => a - b)
  const n = arr.length
  return (arr[n >> 1] + arr[(n - 1) >> 1]) / 2
}

console.log(calMid(odd))
console.log(calMid(even))

export {}
