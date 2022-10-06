# n个人分组
# !n个人编号1到n,分成任意组(不一定要连续分组)
# !使任意一组内不存在两个人的编号模m同余
# 求分成1,2,3,...,n组的方案数,取模.
# n<=5000

# https://www.cnblogs.com/dream1024/p/15254152.html
# 有点像第二类斯特林数的dp？
# !前i个人分成j组的方案数
# 单独分成一组 dp[i-1][j-1]
# 放到前面的某一组 (前面有j个组,但是有 (i-1)//m 个组中存在和i模m同余的数)

# !dp[i][j] = dp[i-1][j-1] + dp[i-1][j] * (j-(i-1)//m)
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, m = map(int, input().split())
    dp = [[0] * (n + 1) for _ in range(n + 1)]
    dp[0][0] = 1
    for i in range(1, n + 1):
        for j in range(1, n + 1):
            dp[i][j] = (dp[i - 1][j - 1] + dp[i - 1][j] * (j - (i - 1) // m)) % MOD

    for i in range(1, n + 1):
        print(dp[n][i])
