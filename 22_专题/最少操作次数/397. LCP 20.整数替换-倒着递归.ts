/**
 *
 * @param n
 * n 变为 1 所需的最小替换次数是多少？
 * 如果 n 是偶数，则用 n / 2替换 n 。
   如果 n 是奇数，则可以用 n + 1或n - 1替换 n 。
 */
function integerReplacement(n: number): number {
  if (n === 1) return 0
  if (n === 2) return 1
  if (n === 3) return 2

  let res = 0
  while (n !== 1) {
    // 偶数直接右移
    // 注意要不带符号移 带符号会有溢出可能
    if ((n & 1) === 0) n >>>= 1
    else n += (n & 2) === 0 || n === 3 ? -1 : 1 // 11直接加1 01直接减1 3特殊处理减1
    res++
  }

  return res
}

console.log(integerReplacement(6))
// 输出：3
// 解释：8 -> 4 -> 2 -> 1

function integerReplacement2(n: number, c = 0): number {
  if (n === 1) return c

  if (n % 2 === 0) {
    return integerReplacement2(n / 2, c + 1)
  } else {
    return Math.min(integerReplacement2(n + 1, c + 1), integerReplacement2(n - 1, c + 1))
  }
}

// 不能用dp因为dp会计算1到n没有用的情况
