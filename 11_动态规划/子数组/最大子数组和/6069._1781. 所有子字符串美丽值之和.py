# 给你一个字符串 s ，请你返回它所有子字符串的 美丽值 之和。
# !一个字符串的 美丽值 定义为：出现频率最高字符与出现频率最低字符的出现次数之差

# 1 <= s.length <= 500
class Solution:
    def beautySum1(self, s: str) -> int:
        """枚举子串：固定左端点,dp O(26*n^2)"""
        res = 0
        for left in range(len(s)):
            counter = [0] * 26
            for right in range(left, len(s)):
                counter[ord(s[right]) - 97] += 1
                res += max(counter) - min(c for c in counter if c > 0)
        return res

    def beautySum(self, s: str) -> int:
        """枚举最高最低是哪两个字符 dp O(26*26*n)"""
        res = 0
        for left in range(len(s)):
            counter = [0] * 26
            for right in range(left, len(s)):
                counter[ord(s[right]) - 97] += 1
                res += max(counter) - min(c for c in counter if c > 0)
        return res


print(Solution().beautySum(s="aabcb"))
# 输出：5
# 解释：美丽值不为零的字符串包括 ["aab","aabc","aabcb","abcb","bcb"] ，每一个字符串的美丽值都为 1 。
