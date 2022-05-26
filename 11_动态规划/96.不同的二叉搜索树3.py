# 有多少种不同的二叉树满足节点个数为n且树的高度不超过m
MOD = int(1e9 + 7)


# dp[i][j] 表示 i 个节点能够组成的高度不超过 j 的树的个数
# m,n<=50
n, m = map(int, input().split())
dp = [[0] * (m + 1) for _ in range(n + 1)]
dp[0] = [1] * (m + 1)

for count in range(1, n + 1):
    for maxHeight in range(1, m + 1):
        for left in range(count):
            #  左子树节点数为k，右子树节点数为i-k-1，且左右子树都要求小于等于j-1
            dp[count][maxHeight] += dp[left][maxHeight - 1] * dp[count - left - 1][maxHeight - 1]
            dp[count][maxHeight] %= MOD


print(dp[n][m])
