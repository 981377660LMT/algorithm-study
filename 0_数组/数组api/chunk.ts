export function chunk<T>(arr: T[], maxSize: number): T[][] {
  if (maxSize <= 1) {
    return arr.map(item => [item])
  }
  const res: T[][] = []
  let ptr = 0
  while (ptr < arr.length) {
    res.push(arr.slice(ptr, ptr + maxSize))
    ptr += maxSize
  }
  return res
}

if (require.main === module) {
  console.log(chunk([1, 2, 3, 4, 5, 6, 7, 8, 9], 3))
  console.log(chunk([1, 2, 3, 4, 5, 6, 7, 8, 9], 1))
  console.log(chunk([1, 2, 3, 4, 5, 6, 7, 8, 9], 9))
}
