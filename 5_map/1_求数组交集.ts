// 时间复杂度O(n+m)
// 空间复杂度O(m)
const intersection = (arr1: number[], arr2: number[]) => {
  const map = new Map<number, boolean>()
  const res: number[] = []

  arr1.forEach(n => {
    map.set(n, true)
  })
  arr2.forEach(n => {
    if (map.get(n)) {
      res.push(n)
      map.delete(n)
    }
  })

  return res
}

console.log(intersection([1, 2, 3], [1]))
export {}
