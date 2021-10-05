// 将一个连分数化成最简分数
function fraction(cont: number[]): number[] {
  if (cont.length === 1) return [cont[0], 1]
  const last = fraction(cont.slice(1))
  return [cont[0] * last[0] + last[1], last[0]]
}

console.log(fraction([3, 2, 0, 2]))
// 输出：[13, 4]
// 解释：原连分数等价于3 + (1 / (2 + (1 / (0 + 1 / 2))))。注意[26, 8], [-13, -4]都不是正确答案。
