# 长度为n的线段有多少种切割方法，切割时两个为一组，一组中第一个数字必须为4，两个数字都必须为`正整数`
from functools import lru_cache


MOD = 998244353

# 注意牛客网容易内存超出限制
# 一组至少分5个消息，至多分n个消息，写出递推式相减
# fn=f(n-5)+f(n-6)+f(n-7)+f(n-8)+f(n-9)+...+f(0)
# 即 fn=f(n-1)+f(n-5)


class Solution:
    def messageCount(self, N: int) -> int:
        # write code here
        if N <= 4:
            return 0
        if N <= 9:
            return 1
        dp = [0] * (N + 1)
        dp[5] = 1
        for i in range(6, N + 1):
            dp[i] = (dp[i - 1] + dp[i - 5]) % MOD
        return dp[-1]


print(Solution().messageCount(10))
print(Solution().messageCount(11))
print(Solution().messageCount(6666))
