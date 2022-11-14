# 给你一个字符串 s 和一个 正 整数 k 。
# 从字符串 s 中选出一组满足下述条件且 不重叠 的子字符串:
# !每个子字符串的长度 至少 为 k 。
# !每个子字符串是一个 回文串 。
# !返回最优方案中能选择的子字符串的 最大 数目。


class Solution:
    def maxPalindromes2(self, s: str, k: int) -> int:
        """
        贪心+马拉车 O(n)
        我们只需要考虑长度为k和k+1的回文串数目就行。
        如果k+2i是回文串,那么掐头去尾,肯定有长度为k的回文串,
        要数目最多，我们就选最短的。
        !只需要判断 [i,i+k-1] 和 [i,i+k]是否为回文串即可，
        !使用 manacher 算法可以在 O(n) 时间内判断一个子串是否为回文串
        """
        # !js-algorithm\17_模式匹配\马拉车拉马\2472. 不重叠回文子字符串的最大数目.py
        ...

    def maxPalindromes1(self, s: str, k: int) -> int:
        """O(n^2)dp"""

        def expand(left: int, right: int) -> None:
            """中心扩展法求s[left:right+1]是否为回文串"""
            while left >= 0 and right < len(s) and s[left] == s[right]:
                if right - left + 1 >= k:
                    isPalindrome[left][right] = True
                left -= 1
                right += 1

        n = len(s)
        isPalindrome = [[False] * n for _ in range(n)]  # dp[i][j] 表示 s[i:j+1] 是否是回文串
        for i in range(n):
            expand(i, i)
            expand(i, i + 1)

        # 选出最多数量的区间，使得它们互不重叠 (dp)
        dp = [0] * (n + 1)  # 第i个字符结尾的最多不重叠回文子串数目
        for i in range(1, n + 1):
            dp[i] = dp[i - 1]  # jump
            for j in range(i - k + 1):  # not jump
                if isPalindrome[j][i - 1]:
                    dp[i] = max(dp[i], dp[j] + 1)
        return dp[-1]


print(Solution().maxPalindromes1(s="abaccdbbd", k=3))
print(Solution().maxPalindromes1(s="iqqibcecvrbxxj", k=1))
print(Solution().maxPalindromes1(s="i" * 2000, k=1))
