"""分割整除序列的方案数


给定一个序列A,现将其分隔为一些非空序列B1,B2,.. . ,Bk且
每段序列Bi的和都可以被 i 整除
n<=3000 Ai<=1e15
https://atcoder.jp/contests/abc207/editorial/2154


dp[j][count]表示前j个数分割为count段的方案数 答案为sum(dp[n][i] for i in range(1,n+1))
朴素的转移为O(n^3) 其中每次 i -> j 转移需要花费O(n)
dp[j][count] = sum(dp[i][count-1] for i in range(j) if (preSum[j]-preSum[i])%count==0))
(即前缀和模count相等时可以转移过来)


如何优化到O(n^2)呢? => 分组的前缀和实现O(1)转移
preSum[j]-preSum[i] 被count整除 <=> preSum[j]%count == preSum[i]%count
用 dpSum[count][mod] 记录dp[index][count]%count 结果为mod这一类dp的和
dp转移就变成
dp[j+1][count] = dpSum[count][preSum[j+1]%count]
dpSum[count][preSum[j+1]%count] += dp[j+1][count-1]
"""

from itertools import accumulate
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)


if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    preSum = [0] + list(accumulate(nums))

    dp = [[0] * (n + 1) for _ in range(n + 1)]
    dpSum = [[0] * (n + 1) for _ in range(n + 1)]
    dp[0][0] = 1
    dpSum[1][0] = 1

    for j in range(n):
        for count in range(1, n + 1):
            dp[j + 1][count] = dpSum[count][preSum[j + 1] % count]
        for count in range(2, n + 1):
            dpSum[count][preSum[j + 1] % count] += dp[j + 1][count - 1]
            dpSum[count][preSum[j + 1] % count] %= MOD
        # print(j, dpSum, dp)

    res = sum(dp[n]) % MOD
    print(res)
