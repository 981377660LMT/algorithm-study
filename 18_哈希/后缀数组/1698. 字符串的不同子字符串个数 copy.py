# 给定一个字符串 s，返回 s 的不同子字符串的个数。
# 字符串的 子字符串 是由原字符串删除开头若干个字符（可能是 0 个）并删除结尾若干个字符（可能是 0 个）形成的字符串。
# 1 <= s.length <= 500

# 你可以以 O(n) 时间复杂度解决此问题吗？
class Solution:
    def countDistinct(self, s: str) -> int:
        return len(set(s[i:j] for i in range(len(s)) for j in range(i + 1, len(s) + 1)))


print(Solution().countDistinct("aabbaba"))
