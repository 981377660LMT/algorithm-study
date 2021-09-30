/**
 * @param {string} n  n的取值范围是 [3, 10^18]。 大于JS的Number.MAX_SAFE_INTEGER 2**53-1
 * @return {string}
 * 如果n的k（k>=2）进制数的所有数位全为1，则称 k（k>=2）是 n 的一个好进制。
 * 以字符串的形式给出 n, 以字符串的形式返回 n 的最小好进制。
 链接：https://leetcode-cn.com/problems/smallest-good-base/solution/er-fen-deng-bi-shu-lie-you-hua-huan-yuan-slrv/
 由于js精度限制 故需将字符串转成BigInt来计算结果
 */
const smallestGoodBase = function (n) {
  const target = BigInt(n)

  const pow = (n, k) => {
    let res = 1n
    for (let i = 0; i < k; i++) {
      res *= n
    }
    return res
  }

  /**
   *
   * @param base  base进制
   * @param len base进制下有len个1
   * @returns  base进制下有len个的值
   */
  const sumInBase = (base, len) => {
    return (1n - pow(base, len)) / (1n - base)
  }

  // 我们可逐一枚举base进制的长度
  // 注意到需要返回的是最小的 k 进制，结合前面说的进制越小 N 越大的知识，
  // 我们应该使用从后往前倒着遍历，这样遇到满足条件的 k 可直接返回
  for (let len = target.toString(2).length; len >= 0; len--) {
    let l = 2n
    let r = target - 1n

    while (l <= r) {
      const mid = (l + r) / 2n
      const sum = sumInBase(mid, len)
      if (sum < target) l = mid + 1n
      else if (sum > target) r = mid - 1n
      else return mid.toString()
    }
  }
}

console.log(smallestGoodBase('13'))
console.log(smallestGoodBase('4681'))
// 输出："3"
// 解释：13 的 3 进制是 111。

// n = int(n)
// // 上面提到的 base 进制转十进制公式
// func sum_with(base, N):
//     return sum(1 * base ** i for i in range(N))

// for k=2 to n - 1:
//     if sum_with(k, N) == n: return k
