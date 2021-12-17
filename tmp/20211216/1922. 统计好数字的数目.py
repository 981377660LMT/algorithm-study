# 偶数 下标处的数字为 偶数 且 奇数 下标处的数字为 质数 （2，3，5 或 7）。
# 请你返回长度为 n 且为好数字的数字字符串 总数
class Solution:
    def countGoodNumbers(self, n: int) -> int:
        def qpow(n, k, mod):
            res = 1
            while k:
                if k & 1:
                    res *= n
                    res %= mod
                n **= 2
                n %= mod
                k >>= 1
            return res

        MOD = 10 ** 9 + 7
        a = qpow(5, (n + 1) // 2, MOD)
        b = qpow(4, n // 2, MOD)

        return (a * b) % MOD


print(Solution().countGoodNumbers(1))
# 输出：5
# 解释：长度为 1 的好数字包括 "0"，"2"，"4"，"6"，"8" 。
