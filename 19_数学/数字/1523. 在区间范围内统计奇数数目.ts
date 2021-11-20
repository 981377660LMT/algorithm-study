// 请你返回 low 和 high 之间（包括二者）奇数的数目。
// 0 <= low <= high <= 10^9
function countOdds(low: number, high: number): number {
  return ((high + 1) >> 1) - (low >> 1)
}

// 规律:0到n间的奇数个数为(n+1)>>1
