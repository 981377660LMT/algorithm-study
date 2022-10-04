# LCS求方案
# !用一个pre(数组/字典)来还原路径 拷贝字符串会TLE/MLE


import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

s = input()
t = input()

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
print("".join(res[::-1]))


# 注意:这里也可以直接根据dp数组倒推
# dp复原
