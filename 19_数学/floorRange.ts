/**
 * 将 [1,n] 内的数分成O(2*sqrt(n))段, 每段内的 n//i 相同.
 * 每个段为(left,right,div)，表示 left <= i <= right 内的 n//i == div.
 */
function floorRange(n: number): { left: number; right: number; div: number }[] {
  if (n <= 0) return []
  const res: { left: number; right: number; div: number }[] = []
  let m = 1
  while (m * m <= n) {
    res.push({ left: m, right: m, div: Math.floor(n / m) })
    m++
  }
  for (let i = m; i > 0; i--) {
    const left = Math.floor(n / (i + 1)) + 1
    const right = Math.floor(n / i)
    if (left <= right && res.length > 0 && res[res.length - 1].right < left) {
      res.push({ left, right, div: Math.floor(n / left) })
    }
  }
  return res
}

export { floorRange }

if (require.main === module) {
  console.log(floorRange(10))
}
