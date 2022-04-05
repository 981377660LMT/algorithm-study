from typing import List


class Solution:
    # 求的就是扩展 KMP（Z 数组）的所有元素之和
    def sumScores(self, s: str) -> int:
        def getZ(string: str) -> List[int]:
            """z算法，求字符串公共前后缀的长度"""

            n = len(string)
            z = [0] * n  # z[i]即是s[i:n]与s的最长公共前缀的长度 (i>=1)
            left, right = 0, 0
            for i in range(1, n):
                z[i] = max(
                    min(z[i - left], right - i + 1), 0
                )  # 注：不用 min max，拆开用 < > 比较会更快（仅限于 Python）
                while i + z[i] < n and string[z[i]] == string[i + z[i]]:
                    left, right = i, i + z[i]
                    z[i] += 1
            return z

        z = getZ(s)
        return sum(z) + len(s)

