# 長さ N の整数からなる数列 A=(A 1 ​ ,…,A N ​ ) であって、以下の条件を全て満たすものは何通りありますか
# 数组每个数在 [1, M] 之间
# 数组和<=k
# 求数组个数
# n,m<=50
# k<=nm

# 数列を先頭から決めていく際、覚えておくべき必要があるものはその時点での数列の総和のみであり、
# 具体的に各要素の値が何であったかは捨象してよい。

# !三种解法:
# !1. 行间转移dp O(N*M*K) 因为要保存各个容量的和 所以不能用一维数组省空间
# !2. 前缀和优化dp O(N*K) dp 的 sum 转化为 dpSum
# !3. 生成函数 O(K)

from collections import defaultdict
from functools import lru_cache
from itertools import accumulate
import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353


def main() -> None:
    # !0.记忆化dfs
    # @lru_cache(None)
    # def dfs(index: int, remain: int) -> int:
    #     if remain < 0:
    #         return 0
    #     if index == n:
    #         return 1 if remain == 0 else 0

    #     res = 0
    #     for cur in range(1, m + 1):
    #         res += dfs(index + 1, remain - cur)
    #         res %= MOD
    #     return res

    # n, m, k = map(int, input().split())
    # res = 0
    # for sum_ in range(n, k + 1):
    #     res += dfs(0, sum_)
    #     res %= MOD
    # dfs.cache_clear()
    # print(res)

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
    # n, m, k = map(int, input().split())
    # dp = [0] * (k + 1)
    # for i in range(1, m + 1):
    #     dp[i] = 1

    # for _ in range(1, n):
    #     ndp = [0] * (k + 1)
    #     dpSum = [0] + list(accumulate(dp, func=(lambda x, y: (x + y) % MOD)))
    #     # ndp[i] = sum(dp[i - j] for j in range(1, m + 1)) = dpSum[i] - dpSum[i - m]
    #     for i in range(k + 1):
    #         ndp[i] = (dpSum[i] - dpSum[max(0, i - m)]) % MOD
    #     dp = ndp

    # print(sum(dp) % MOD)

    # !3. 生成函数
    # f(x) = x + x^2 + ... + x^m = x*((1-x^m)/(1-x))
    # 求 f(x)^n 0次项到k次项的系数之和
    # !注意到f(x)从0次项到k次项的系数之和 等于 f(x)/(1-x)的k次项的系数 (这个可以左边乘(1-x)看出来)
    # !即要求 (1-x^m)^n * (1-x)^-(n+1) 的 k-n 次项系数
    # https://atcoder.jp/contests/abc248/submissions/31149808
    n, m, k = map(int, input().split())


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
