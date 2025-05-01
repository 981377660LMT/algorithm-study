# 1092. 最短公共超序列
# 给出两个字符串 str1 和 str2，返回同时以 str1 和 str2 作为子序列的最短字符串。
# 如果答案不止一个，则可以返回满足条件的任意一个答案。
# 1 <= str1.length, str2.length <= 1000
# LCS变形
# !思路是首先获取到LCS，然后重新从头“走一遍”，得到从(0,0)到(m,n)的路径，这个路径即我们需要的超序列


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def shortestCommonSupersequence(self, str1: str, str2: str) -> str:
        n, m = len(str1), len(str2)
        # dp[i][j] = LCS 长度 of str1[i:] 和 str2[j:]
        dp = [[0] * (m + 1) for _ in range(n + 1)]
        for i in range(n - 1, -1, -1):
            for j in range(m - 1, -1, -1):
                if str1[i] == str2[j]:
                    dp[i][j] = dp[i + 1][j + 1] + 1
                else:
                    dp[i][j] = max2(dp[i + 1][j], dp[i][j + 1])

        i = j = 0
        res = []
        while i < n and j < m:
            if str1[i] == str2[j]:
                res.append(str1[i])
                i += 1
                j += 1
            else:
                if dp[i + 1][j] >= dp[i][j + 1]:
                    res.append(str1[i])
                    i += 1
                else:
                    res.append(str2[j])
                    j += 1

        if i < n:
            res.append(str1[i:])
        if j < m:
            res.append(str2[j:])

        return "".join(res)


print(Solution().shortestCommonSupersequence(str1="abac", str2="cab"))
# 输出："cabac"
# 解释：
# str1 = "abac" 是 "cabac" 的一个子串，因为我们可以删去 "cabac" 的第一个 "c"得到 "abac"。
# str2 = "cab" 是 "cabac" 的一个子串，因为我们可以删去 "cabac" 末尾的 "ac" 得到 "cab"。
# 最终我们给出的答案是满足上述属性的最短字符串。
