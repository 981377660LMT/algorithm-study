# F - Potion
# 魔法药水(背包dp)
# 有n个数，你可以在第0个时刻条任意多个数出来，并让计数器加上他们的和，
# 设你挑了k个数，那么之后每个时刻计数器都会加k，你要最小化计数器恰好到x的时间
# !n<=100 Ai<=1e7 1e9<=x<=1e18
# !枚举选的个数(固定k会比较方便) dp[index][remain][mod] O(n^4) =>
# !选k个数,模k等于target,求选出的数之和的最大值


from functools import lru_cache
from typing import List


def potion(magic: List[int], target: int) -> int:
    def cal(k: int) -> int:
        @lru_cache(None)
        def dfs(index: int, remain: int, mod: int) -> int:
            if remain < 0:
                return -INF
            if index == n:
                return 0 if ((remain == 0) and (mod == curTarget)) else -INF
            res = dfs(index + 1, remain, mod)  # jump
            if remain:  # select
                res = max(res, magic[index] + dfs(index + 1, remain - 1, (mod + magic[index]) % k))
            return res

        curTarget = target % k
        res = dfs(0, k, 0)
        dfs.cache_clear()
        return res

    n = len(magic)
    res = INF
    for k in range(1, n + 1):
        maxSum = cal(k)  # !选k个数,模k等于target,求选出的数的最大值
        if maxSum > 0:
            res = min(res, (target - maxSum) // k)
    return res


if __name__ == "__main__":

    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    MOD = 998244353
    INF = int(4e18)

    n, target = map(int, input().split())
    nums = list(map(int, input().split()))
    print(potion(nums, target))
