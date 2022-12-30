# 1092. 最短公共超序列
# 给出两个字符串 str1 和 str2，返回同时以 str1 和 str2 作为子序列的最短字符串。
# 如果答案不止一个，则可以返回满足条件的任意一个答案。
# 1 <= str1.length, str2.length <= 1000
# LCS变形
# 思路是首先获取到LCS，然后重新从头“走一遍”，得到从(0,0)到(m,n)的路径，这个路径即我们需要的超序列


class Solution:
    def shortestCommonSupersequence(self, str1: str, str2: str) -> str:
        m, n = len(str1), len(str2)
        lcs = getLCS(str1, str2)
        i, j = 0, 0
        sb = []
        for char in lcs:
            while i < m and str1[i] != char:
                sb.append(str1[i])
                i += 1
            while j < n and str2[j] != char:
                sb.append(str2[j])
                j += 1
            sb.append(char)
            i += 1
            j += 1

        sb.append(str1[i:])
        sb.append(str2[j:])

        return "".join(sb)


def getLCS(s: str, t: str) -> str:
    n1, n2 = len(s), len(t)
    dp = [[0] * (n2 + 1) for _ in range(n1 + 1)]
    pre = [[(0, 0)] * (n2 + 1) for _ in range(n1 + 1)]
    for i in range(1, n1 + 1):
        for j in range(1, n2 + 1):
            if s[i - 1] == t[j - 1]:
                dp[i][j] = dp[i - 1][j - 1] + 1
                pre[i][j] = (i - 1, j - 1)
            else:
                if dp[i][j - 1] > dp[i][j]:
                    dp[i][j] = dp[i][j - 1]
                    pre[i][j] = (i, j - 1)
                if dp[i - 1][j] > dp[i][j]:
                    dp[i][j] = dp[i - 1][j]
                    pre[i][j] = (i - 1, j)

    res = []
    curI, curJ = n1, n2
    while 0 not in (curI, curJ):
        if curI - 1 < n1 and curJ - 1 < n2 and s[curI - 1] == t[curJ - 1]:
            res.append(s[curI - 1])
        curI, curJ = pre[curI][curJ]
    return "".join(res[::-1])


print(Solution().shortestCommonSupersequence(str1="abac", str2="cab"))
# 输出："cabac"
# 解释：
# str1 = "abac" 是 "cabac" 的一个子串，因为我们可以删去 "cabac" 的第一个 "c"得到 "abac"。
# str2 = "cab" 是 "cabac" 的一个子串，因为我们可以删去 "cabac" 末尾的 "ac" 得到 "cab"。
# 最终我们给出的答案是满足上述属性的最短字符串。
