from functools import lru_cache
import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 正整数列 A=(a
# 1
# ​
#  ,a
# 2
# ​
#  ,…,a
# N
# ​
#  ) が与えられます。
# あなたは次の操作を 0 回以上何度でも繰り返せます。

# A から（空でない）連続する部分列を選び、削除する。
# x=1,2,…,M に対し、次の問題を解いてください。

# A の要素の総和をちょうど x にするために必要な操作回数の最小値を求めてください。ただし、どのように操作を行っても A の要素の総和をちょうど x にできない場合は代わりに -1 と出力してください。
# なお、A が空である時、A の要素の総和は 0 であるとします。
if __name__ == "__main__":
    n, m = map(int, input().split())
    nums = list(map(int, input().split()))

    res = [INF] * (3010)
    dp = [[INF] * 2 for _ in range(3010)]  # (curSum, preSelect): remove
    dp[0][0] = 1
    dp[nums[0]][1] = 0
    res[nums[0]] = 1 if n > 1 else 0
    for i in range(1, n):
        ndp = [[INF] * 2 for _ in range(3010)]
        for j in range(3010):
            for k in range(2):
                if dp[j][k] == INF:
                    continue
                ndp[j][0] = min(ndp[j][0], dp[j][k] + (k == 1))
                res[j] = min(res[j], ndp[j][0])
                if j + nums[i] <= m:
                    ndp[j + nums[i]][1] = min(ndp[j + nums[i]][1], dp[j][k])
                    res[j + nums[i]] = min(res[j + nums[i]], ndp[j + nums[i]][1] + (i != n - 1))
        dp = ndp

    for i in range(1, m + 1):
        print(res[i] if res[i] != INF else -1)
