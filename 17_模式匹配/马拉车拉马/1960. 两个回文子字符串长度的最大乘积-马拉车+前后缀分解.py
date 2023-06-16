from itertools import accumulate
from Manacher import Manacher


# 你需要找到两个 `不重叠` 的回文 子字符串，它们的长度都必须为 奇数 ，使得它们长度的乘积最大。
# 2 <= s.length <= 1e5
# https://leetcode-cn.com/problems/maximum-product-of-the-length-of-two-palindromic-substrings/comments/1177272


class Solution:
    def maxProduct(self, s: str) -> int:
        # 1.预处理前缀最长回文长度和后缀最长回文长度，处理出每个位置处的leftMax 和 rightMax数组
        # 2.枚举分割点 leftMax[i]*rightMax[i+1]即可
        n, M = len(s), Manacher(s)
        leftMax = list(accumulate([M.getLongestOddEndsAt(i) for i in range(n)], max))
        rightMax = list(
            accumulate([M.getLongestOddStartsAt(i) for i in range(n - 1, -1, -1)], max)
        )[::-1]
        res = 0
        for i in range(n - 1):
            res = max(res, leftMax[i] * rightMax[i + 1])
        return res


print(Solution().maxProduct(s="ababbb"))
print(Solution().maxProduct(s="ggbswiymmlevedhkbdhntnhdbkhdevelmmyiwsbgg"))
