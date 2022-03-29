# 注意是最长公共子串不是子序列
s1, s2 = input(), input()

m, n = len(s1), len(s2)
dp = [[0] * (n + 1) for _ in range(m + 1)]
res = 0
for i in range(1, m + 1):
    for j in range(1, n + 1):
        if s1[i - 1] == s2[j - 1]:
            dp[i][j] = dp[i - 1][j - 1] + 1
            res = max(res, dp[i][j])
print(res)
