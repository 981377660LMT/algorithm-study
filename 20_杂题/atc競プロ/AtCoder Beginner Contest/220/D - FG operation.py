"""
给定一个长度为n的序列,每个元素都是0~9的一个数。
现在有两种操作:
!操作一:将最左端的两个值x和y变成一个值(x +y) % 10
!操作二:将最左端的两个值x和y变成一个值(x*y) % 10
一共有2^(n -1)中方案,问所有方案最后有多少个0~9

n<=2e5

线性dp
"""

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    dp = [0] * 10
    a, b = nums[0], nums[1]
    dp[(a + b) % 10] += 1
    dp[(a * b) % 10] += 1

    for i in range(2, n):
        ndp = [0] * 10
        cur = nums[i]
        for pre in range(10):
            ndp[(pre + cur) % 10] = (ndp[(pre + cur) % 10] + dp[pre]) % MOD
            ndp[(pre * cur) % 10] = (ndp[(pre * cur) % 10] + dp[pre]) % MOD
        dp = ndp

    print(*dp, sep="\n")
