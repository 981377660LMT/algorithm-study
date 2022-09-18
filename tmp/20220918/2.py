# 给你一个仅由小写英文字母组成的字符串 s ，返回其 最长字母序连续子字符串 的长度。
class Solution:
    def longestContinuousSubstring(self, s: str) -> int:
        res = 1
        dp = 1
        for i in range(1, len(s)):
            if ord(s[i]) - ord(s[i - 1]) == 1:
                dp += 1
            else:
                dp = 1
            res = max(res, dp)
        return res
