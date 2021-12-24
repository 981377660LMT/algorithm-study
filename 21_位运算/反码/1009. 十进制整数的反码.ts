// 找到前导0个数
// 对于 n 位二进制而言，按位取反相当于用 n 个 1 减去原数
// 那么如何得到 n 个 1 呢？首先找原数最高位的 1 在哪，然后将其左移一位，再减去 1，就得到了 n 个 1。
function bitwiseComplement(n: number): number {
  // Count Leading Zeros
  if (n === 0) return 1
  const clz = Math.clz32(n)
  const allOne = (1 << (32 - clz)) - 1
  return n ^ allOne
}

console.log(bitwiseComplement(10))
// 输出：5
// 解释：10 的二进制表示为 "1010"，其二进制反码为 "0101"，也就是十进制中的 5 。
console.log(bitwiseComplement(0))
// 1

function bitwiseComplement2(n: number): number {
  if (n === 0) return 1

  let upper = 1
  while (upper <= n) {
    upper <<= 1
  }

  return n ^ (upper - 1)
}
