# 这些子串满足替换 一个不同字符 以后，是 t 串的子串
# 1 <= s.length, t.length <= 100

# O(n^3)
class Solution:
    def countSubstrings(self, s: str, t: str) -> int:
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

