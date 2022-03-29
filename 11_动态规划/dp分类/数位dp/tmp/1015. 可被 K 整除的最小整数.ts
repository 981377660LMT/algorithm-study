/**
 * @param {number} k  1 <= K <= 10^5
 * @return {number}
 * 给定正整数 K，你需要找出可以被 K 整除的、仅包含数字 1 的最小正整数 N。
 * 返回 N 的长度。如果不存在这样的 N，就返回 -1。
 * @description
 * 1 % 10 = 1
   11 % 6 = 5
   111 % 6 = 3
   1111 % 6 = 1
   11111 % 6 = 5
   111111 % 6 = 3
 */
function smallestRepunitDivByK(k: number): number {
  const shouldEnd = new Set([1, 3, 7, 9])
  if (!shouldEnd.has(k % 10)) return -1

  let mod = 0
  const modSet = new Set<number>()

  for (let i = 1; i < k + 1; i++) {
    mod = (10 * mod + 1) % k // 下一个数
    if (mod === 0) return i
    if (modSet.has(mod)) return -1
    modSet.add(mod)
  }

  return -1
}

console.log(smallestRepunitDivByK(3))

export {}
