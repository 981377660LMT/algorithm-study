# 長さ N の整数からなる数列 A=(A 1 ​ ,…,A N ​ ) であって、以下の条件を全て満たすものは何通りありますか
# 数组每个数在 [1, M] 之间
# 数组和<=k
# 求数组个数
# n,m<=50
# k<=nm


# !1. 行间转移dp O(N*M*K) 因为要保存各个容量的和 所以不能用一维数组省空间
# !2. 前缀和优化dp O(N*K) dp 的 sum 转化为 dpSum


from itertools import accumulate
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353


def main() -> None:
    # !1. 背包 行间转移
    # n, m, k = map(int, input().split())
    # dp = [0] * (k + 1)
    # for i in range(1, m + 1):
    #     dp[i] = 1

    # for _ in range(1, n):
    #     ndp = [0] * (k + 1)
    #     for cur in range(1, m + 1):
    #         for pre in range(k + 1):
    #             if pre + cur <= k:
    #                 ndp[pre + cur] += dp[pre]
    #                 ndp[pre + cur] %= MOD
    #     dp = ndp
    # print(sum(dp) % MOD)

    # !2. 前缀和优化背包的行间转移的范围dp求和
    n, m, k = map(int, input().split())
    dp = [0] * (k + 1)
    for i in range(1, m + 1):
        dp[i] = 1

    for _ in range(1, n):
        ndp = [0] * (k + 1)
        dpSum = [0] + list(accumulate(dp, func=(lambda x, y: (x + y) % MOD)))
        # ndp[i] = sum(dp[i - j] for j in range(1, m + 1)) = dpSum[i] - dpSum[i - m]
        for i in range(k + 1):
            ndp[i] = (dpSum[i] - dpSum[max(0, i - m)]) % MOD
        dp = ndp

    print(sum(dp) % MOD)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
