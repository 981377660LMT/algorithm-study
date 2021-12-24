function baseNeg2(n: number): string {
  if (n === 0) return '0'

  const sb: number[] = []
  while (n !== 0) {
    const [div, mod] = [Math.floor(n / 2), n % 2]
    sb.push(Math.abs(mod))
    // 这里不同
    n = -div
  }

  return sb.reverse().join('')
}

function baseNeg22(n: number): string {
  if (n === 0 || n === 1) return String(n)
  // 这里不同
  return baseNeg22(-(n >> 1)) + String(n & 1)
}

console.log(baseNeg2(2))
console.log(baseNeg22(2))
// 输出："110"
// 解释：(-2) ^ 2 + (-2) ^ 1 = 2
