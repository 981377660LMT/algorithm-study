/**
 * @param n 给定 N，返回任意漂亮数组 A（保证存在一个）。
 * 数组 A 是整数 1, 2, ..., N 组成的排列
 * 对于每个 i < j，都不存在 k 满足 i < k < j 使得 A[k] * 2 = A[i] + A[j]。
 * @summary
 * 利用性质奇数 + 偶数 = 奇数
 * 性质 1： 如果数组 A 是 漂亮数组，那么将 A 中的每一个数 x 进行 kx + b 的映射，其仍然为漂亮数组。
 * 性质 2：如果数组 A 和 B 分别是不同奇偶性的漂亮数组，那么将 A 和 B 拼接起来仍为漂亮数组。
 */
function beautifulArray(n: number): number[] {
  let res: number[] = [1]

  while (res.length < n) {
    const tmp: number[] = []
    for (const num of res) {
      if (num * 2 - 1 <= n) tmp.push(num * 2 - 1)
    }

    for (const num of res) {
      if (num * 2 <= n) tmp.push(num * 2)
    }

    res = tmp
  }

  return res
}

console.log(beautifulArray(4))

// 所以我们假设一个{1-m}的数组是漂亮数组，可以通过下面的方式构造漂亮数组{1-2m}:

// 对{1-m}中所有的数乘以2-1，构成一个奇数漂亮数组A。如{1,3,2,4},可以得到{1,5,3,7}
// 对{1-m}中所有的数乘以2,构成一个偶数漂亮数组B,如{1,3,2,4}, 可以得到{2,6,4,8}
// A+B构成了{1-2m}的漂亮数组。{1,5,3,7}+{2,6,4,8}={1,5,3,7,2,6,4,8}
// 从中删除不要的数字即可。

export {}
