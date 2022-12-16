# 给你一个字符串 s ，请你返回它所有子字符串的 美丽值 之和。
# !一个字符串的 美丽值 定义为：出现频率最高字符与出现频率最低字符的出现次数之差

# 1 <= s.length <= 500
from collections import defaultdict


class Solution:
    def beautySum(self, s: str) -> int:
        """枚举子串：固定左端点,dp O(26*n^2)"""
        res = 0
        for left in range(len(s)):
            counter = defaultdict(int)
            for right in range(left, len(s)):
                counter[s[right]] += 1
                res += max(counter.values()) - min(counter.values())
        return res


print(Solution().beautySum(s="aabcb"))
# 输出：5
# 解释：美丽值不为零的字符串包括 ["aab","aabc","aabcb","abcb","bcb"] ，每一个字符串的美丽值都为 1 。
