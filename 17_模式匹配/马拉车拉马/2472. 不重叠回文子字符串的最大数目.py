# 给你一个字符串 s 和一个 正 整数 k 。
# 从字符串 s 中选出一组满足下述条件且 不重叠 的子字符串：
# 每个子字符串的长度 至少 为 k 。
# 每个子字符串是一个 回文串 。
# 返回最优方案中能选择的子字符串的 最大 数目。
# 子字符串 是字符串中一个连续的字符序列。

# 2472. 不重叠回文子字符串的最大数目
"""
贪心+马拉车 O(n)
我们只需要考虑长度为k和k+1的回文串数目就行。
如果k+2i是回文串,那么掐头去尾,肯定有长度为k的回文串,
要数目最多，我们就选最短的。
!只需要判断 [i,i+k-1] 和 [i,i+k]是否为回文串即可，
!使用 manacher 算法可以在 O(n) 时间内判断一个子串是否为回文串
"""
from Manacher import Manacher


class Solution:
    def maxPalindromes(self, s: str, k: int) -> int:
        M = Manacher(s)
        n = len(s)
        dp = [0] * (n + 1)
        for i in range(1, n + 1):
            dp[i] = dp[i - 1]
            if i - k >= 0 and M.isPalindrome(i - k, i):
                dp[i] = max(dp[i], dp[i - k] + 1)
            if i - k - 1 >= 0 and M.isPalindrome(i - k - 1, i):
                dp[i] = max(dp[i], dp[i - k - 1] + 1)
        return dp[n]


assert Solution().maxPalindromes(s="abaccdbbd", k=3) == 2
