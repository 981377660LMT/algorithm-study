# F - Deforestation
# 砍伐森林

# n棵树, 砍掉第i棵树的代价为h[i-1]+h[i]+h[i+1]
# 之后, h[i]变为0
# !砍掉这n课树有n!种砍树顺序,问有多少种是最小代价的砍树方案
# n<=4000

# https://blog.hamayanhamayan.com/entry/2021/07/11/154241
# 插入dp/
# !1.插入dp
# 一般的数列dp在状态转移时，都是在`末尾`加入新的元素进行迁移
# !而插入dp则是在已经确定的数列的`任意位置`加入新的元素进行迁移
# dp[i][pos] 表示砍掉前i棵树,第i棵树砍掉的顺序位于第pos时的最小代价
# 例如dp[4][2] 表示砍掉前4棵树,第4棵树砍掉的顺序位于第2棵树时的最小代价
# 1 4 2 3
# 1 4 3 2
# 2 4 1 3
# 2 4 3 1
# 3 4 1 2
# 3 4 2 1
# 转移时考虑第5棵树的位置
# 如果heights[4]<heights[5] 那么5应该在4前砍掉 即 dp[4][2] => dp[5][1]+dp[5][2]
#  1 4 2 3
# ^ ^
# 因此可以暴力更新

# !2.前缀和优化dp(改成貰うdp)


from itertools import accumulate
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9) + 7
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    heights = list(map(int, input().split()))
    dp = [[0] * n for _ in range(n)]  # dp[i][j]表示砍掉前i棵树,第i棵树砍掉的顺序位于第j时的最小代价(0-indexed)
    dp[0][0] = 1
    for i in range(1, n):  # 前面有i棵树了
        preSum = [0] + list(accumulate(dp[i - 1], lambda x, y: (x + y) % MOD))
        if heights[i - 1] <= heights[i]:
            for curPos in range(i + 1):
                # for prePos in range(curPos, i):
                #     dp[i][curPos] = (dp[i][curPos] + dp[i - 1][prePos]) % MOD
                dp[i][curPos] += preSum[i] - preSum[curPos]
                dp[i][curPos] %= MOD

        if heights[i - 1] >= heights[i]:
            for curPos in range(i + 1):
                # for prePos in range(curPos):
                #     dp[i][curPos] = (dp[i][curPos] + dp[i - 1][prePos]) % MOD
                dp[i][curPos] += preSum[curPos]
                dp[i][curPos] %= MOD

    res = 0
    for i in range(n):
        res += dp[n - 1][i]
        res %= MOD
    print(res)
