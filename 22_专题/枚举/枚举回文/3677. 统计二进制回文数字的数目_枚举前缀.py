# 3677. 统计二进制回文数字的数目
# 给你一个 非负 整数 n。
# 如果一个 非负 整数的二进制表示（不含前导零）正着读和倒着读都一样，则称该数为 二进制回文数。
# 返回满足 0 <= k <= n 且 k 的二进制表示是回文数的整数 k 的数量。
# 注意： 数字 0 被认为是二进制回文数，其表示为 "0"。


class Solution:
    def countBinaryPalindromes(self, n: int) -> int:
        s = bin(n)[2:]
        m = len(s)
        res = 1
        for l in range(1, m):
            res += 1 << (l - 1) // 2
        halfLen = (m + 1) // 2
        half = s[:halfLen]
        palindrome = half + half[: m // 2][::-1]
        if int(palindrome, 2) <= n:
            res += 1
        halfNum = int(half, 2)
        res += halfNum - (1 << (halfLen - 1))
        return res
