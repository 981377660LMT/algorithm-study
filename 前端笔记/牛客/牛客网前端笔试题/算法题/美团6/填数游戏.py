# _+_+_+...+_=n
# 这里可以填任意多个正整数（甚至可能是1个），只要这些数的和等于n即可。
# 但是，有一个额外的限制，填入的所有数必须小于等于k，大于等于1，填入的数的最大值必须大于等于d。
# 请你计算，有多少个不同的等式满足这些限制。由于答案可能很大，请将答案mod(998244353)后输出。
from functools import lru_cache
import sys

sys.setrecursionlimit(100000)

MOD = 998244353

target, k, d = map(int, input().split())


# n,k,d<=1000
@lru_cache(None)
def dfs(curSum: int, alreadySatisfied: bool) -> int:
    if curSum > target:
        return 0
    if curSum == target:
        return int(alreadySatisfied)

    res = 0
    for cur in range(1, k + 1):
        if cur >= d:
            res += dfs(curSum + cur, True) % MOD
        else:
            res += dfs(curSum + cur, alreadySatisfied) % MOD
    return res % MOD


res = dfs(0, False)
dfs.cache_clear()
print(res)
