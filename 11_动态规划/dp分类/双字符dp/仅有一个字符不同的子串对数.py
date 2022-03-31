# Number of Substrings with Single Character Difference
# n,m<=100
class Solution2:
    def solve(self, s, t):
        # 暴力O(mnmin(m,n))
        res = 0
        for i in range(len(s)):
            for j in range(len(t)):
                diff = 0
                for k in range(min(len(s) - i, len(t) - j)):
                    if s[i + k] != t[j + k]:
                        diff += 1
                    if diff == 1:
                        res += 1
                    elif diff > 1:
                        break
        return res


class Solution:
    def solve(self, s, t):
        # 预处理O(mn)
        m, n = len(s), len(t)
        prefix = [[0] * (n + 1) for _ in range(m + 1)]
        suffix = [[0] * (n + 1) for _ in range(m + 1)]

        for i in range(1, m + 1):
            for j in range(1, n + 1):
                if s[i - 1] == t[j - 1]:
                    prefix[i][j] = prefix[i - 1][j - 1] + 1

        for i in range(m - 1, -1, -1):
            for j in range(n - 1, -1, -1):
                if s[i] == t[j]:
                    suffix[i][j] = suffix[i + 1][j + 1] + 1

        res = 0
        for i in range(1, m + 1):
            for j in range(1, n + 1):
                if s[i - 1] != t[j - 1]:
                    # 如果 s[i] != s[j]，我们找到了 (prefix[i-1][j-1] + 1) * (suffix[i-1][j-1] + 1) 个符合条件的字符组合。也就是前缀+1 和后缀长度+1 的笛卡尔积。
                    res += (prefix[i - 1][j - 1] + 1) * (suffix[i][j] + 1)
        return res
