# 这题有点像
# 2_子数组区间最值.py

# 1 <= s.length <= 500
class Solution:
    def beautySum(self, s: str) -> int:
        """枚举子串：固定左端点,dp"""
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

