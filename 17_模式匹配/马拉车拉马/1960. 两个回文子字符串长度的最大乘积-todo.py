from typing import Sequence
from Manacher import Manacher


# 哈希
# 你需要找到两个 不重叠的回文 子字符串，它们的长度都必须为 奇数 ，使得它们长度的乘积最大。
# 2 <= s.length <= 105
# https://leetcode-cn.com/problems/maximum-product-of-the-length-of-two-palindromic-substrings/comments/1177272


class Solution:
    def maxProduct(self, s: str) -> int:
        # 1.预处理前缀最长回文长度和后缀最长回文长度，处理出每个位置处的leftMax 和 rightMax数组
        # 2.枚举分割点 leftMax[i]*rightMax[i+1]即可
        manacher = Manacher(s)

