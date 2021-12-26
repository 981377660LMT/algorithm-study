# 1 <= str1.length, str2.length <= 1000
# LCS变形
# 思路是首先获取到LCS，然后重新从头“走一遍”，得到从(0,0)到(m,n)的路径，这个路径即我们需要的超序列
class Solution:
    def shortestCommonSupersequence(self, str1: str, str2: str) -> str:
        m, n = len(str1), len(str2)
        dp = [[''] * (n + 1) for _ in range(m + 1)]

        for i in range(1, m + 1):
            for j in range(1, n + 1):
                if str1[i - 1] == str2[j - 1]:
                    dp[i][j] = dp[i - 1][j - 1] + str1[i - 1]
                else:
                    if len(dp[i - 1][j]) > len(dp[i][j - 1]):
                        dp[i][j] = dp[i - 1][j]
                    else:
                        dp[i][j] = dp[i][j - 1]
        print(dp)
        i, j = 0, 0
        lcs = dp[m][n]
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

        return ''.join(sb)


print(Solution().shortestCommonSupersequence(str1="abac", str2="cab"))
# 输出："cabac"
# 解释：
# str1 = "abac" 是 "cabac" 的一个子串，因为我们可以删去 "cabac" 的第一个 "c"得到 "abac"。
# str2 = "cab" 是 "cabac" 的一个子串，因为我们可以删去 "cabac" 末尾的 "ac" 得到 "cab"。
# 最终我们给出的答案是满足上述属性的最短字符串。
