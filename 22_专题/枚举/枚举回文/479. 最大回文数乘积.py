# 给定一个整数 n ，返回 可表示为两个 n 位整数乘积的 最大回文整数 。
# 因为答案可能非常大，所以返回它对 1337 取余 。

# 1 <= n <= 8
# 回文串最多17位


from enumeratePalindrome import emumeratePalindromeByLength

MOD = 1337


class Solution:
    def largestPalindrome(self, n: int) -> int:
        iter = emumeratePalindromeByLength(1, 2 * n, reverse=True)
        lower, upper = 10 ** (n - 1), 10**n - 1
        for p in iter:
            # 这里用大的范围，是因为从后往前找，[lower,sqrt]里的数多 [sqrt,upper]里的数少
            num = int(p)
            for factor in range(upper, int(num**0.5), -1):
                if num % factor == 0 and upper >= num // factor >= lower:
                    return num % MOD
        return -1
