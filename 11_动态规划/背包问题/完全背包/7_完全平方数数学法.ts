// 给你一个整数 n ，返回和为 n 的完全平方数的 最少数量 。
// 这个方法太慢了复杂度O(n**3/2)

// 根据 拉格朗日四平方和定理，可以得知答案必定为 1, 2, 3, 4 中的一个
// 拉格朗日四平方和定理:每个正整数均可表示为四个整数的平方和（包括0的平方）
// 根据勒让德三平方和定理，可以得知不满足n= (4^k) * (8m+7)的必然是最多写成三个数的平方和
// 即如果n=(4^k) * (8m+7)那么必定由四个平方数组成

// 根据上面两个定理，首先判断结果是否 1。
// 然后判断结果是否为 4。
// 接着枚举结果是否为 2。
// 最后上面条件都不满足，结果为 3。
// O(√n+logn)
const numSquares = (n: number) => {
  const isSquare = (num: number) => Math.floor(Math.sqrt(num)) ** 2 === num
  if (isSquare(n)) return 1

  while (n % 4 === 0) n /= 4
  if (n % 8 === 7) return 4

  for (let i = 0; i <= n; i++) {
    if (isSquare(n - i * i)) return 2
  }

  return 3
}

console.log(numSquares(12))

export {}
