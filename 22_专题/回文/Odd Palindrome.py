# 给定字符串的所有回文子串是否都是奇数长度的

# 可以从中心扩展法求回文子串观察 不存在偶数长度回文等价于相邻字符串不相等
class Solution:
    def solve(self, s: str) -> bool:
        return all(pre != cur for pre, cur in zip(s, s[1:]))


print(Solution().solve(s="bab"))
# True
