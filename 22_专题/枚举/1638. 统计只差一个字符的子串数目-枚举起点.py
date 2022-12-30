# 这些子串满足替换 一个不同字符 以后，是 t 串的子串
# 1 <= s.length, t.length <= 100

# O(n^3)
# 1638. 统计只差一个字符的子串数目
class Solution:
    def countSubstrings(self, s: str, t: str) -> int:
        """只有一个字符不同的子串对数"""
        m, n = len(s), len(t)
        res = 0

        # 枚举起点
        for i in range(m):
            for j in range(n):
                k = 0
                diff = 0
                while i + k < m and j + k < n:
                    if s[i + k] != t[j + k]:
                        diff += 1
                    if diff >= 2:
                        break
                    if diff == 1:
                        res += 1
                    k += 1

        return res

    def countSubstrings2(self, s: str, t: str) -> int:
        """不同字符=相同前缀+不同字符+相同后缀

        因此需要处理出前后缀的lcp(最长公共前缀)
        """
        m, n = len(s), len(t)
        lcp1 = [[0] * (n + 1) for _ in range(m + 1)]  # lcp[i][j] 表示前缀 num[:i+1]与num[:j+1]的最长公共前缀长度
        lcp2 = [[0] * (n + 1) for _ in range(m + 1)]  # lcp[i][j] 表示后缀 num[i:]与num[j:]的最长公共前缀长度

        for i in range(1, m + 1):
            for j in range(1, n + 1):
                if s[i - 1] == t[j - 1]:
                    lcp1[i][j] = lcp1[i - 1][j - 1] + 1

        for i in range(m - 1, -1, -1):
            for j in range(n - 1, -1, -1):
                if s[i] == t[j]:
                    lcp2[i][j] = lcp2[i + 1][j + 1] + 1

        res = 0
        for i in range(1, m + 1):
            for j in range(1, n + 1):
                if s[i - 1] != t[j - 1]:
                    # 如果 s[i] != s[j]，我们找到了 (prefix[i-1][j-1] + 1) * (suffix[i-1][j-1] + 1) 个符合条件的字符组合。
                    # 也就是前缀+1 和后缀长度+1 的笛卡尔积。
                    res += (lcp1[i - 1][j - 1] + 1) * (lcp2[i][j] + 1)
        return res


print(Solution().countSubstrings(s="aba", t="baba"))
# 输出：6
# 解释：以下为只相差 1 个字符的 s 和 t 串的子字符串对：
# ("aba", "baba")
# ("aba", "baba")
# ("aba", "baba")
# ("aba", "baba")
# ("aba", "baba")
# ("aba", "baba")
# 加粗部分分别表示 s 和 t 串选出来的子字符串。
