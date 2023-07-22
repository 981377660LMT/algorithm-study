from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串 word 和一个字符串数组 forbidden 。

# 如果一个字符串不包含 forbidden 中的任何字符串，我们称这个字符串是 合法 的。

# 请你返回字符串 word 的一个 最长合法子字符串 的长度。


# 子字符串 指的是一个字符串中一段连续的字符，它可以为空


class Solution:
    def longestValidSubstring(self, word: str, forbidden: List[str]) -> int:
       res, left, n = 0, 0, len(nums)
       curSum = 0
       for right in range(n):
           curSum += nums[right]
           while left <= right and :
               curSum -= nums[left]
               left += 1
           res = max(res, right - left + 1)
       return res
