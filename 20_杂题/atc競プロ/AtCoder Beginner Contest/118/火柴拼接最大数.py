# 火柴拼接最大数
# D - Match Matching
# https://atcoder.jp/contests/abc118/tasks/abc118_d
# 给你n根(n<=1e4)火柴 求恰好用完n根能够拼接出的最大数字是多少

from functools import lru_cache
import sys
from typing import List

COST = {
    1: 2,
    2: 5,
    3: 5,
    4: 4,
    5: 5,
    6: 6,
    7: 3,
    8: 7,
    9: 6,
}

INF = int(1e18)


def matchMathing(n: int, digits: List[int]) -> int:
    @lru_cache(None)
    def dfs(remain: int) -> int:
        if remain <= 0:
            return 0 if remain == 0 else -INF

        res = -INF
        for cur in digits:
            res = max(res, dfs(remain - COST[cur]) * 10 + cur)
        return res

    res = dfs(n)
    dfs.cache_clear()
    return res


if __name__ == "__main__":
    sys.setrecursionlimit(int(1e6))
    n, _ = map(int, input().split())
    digits = list(map(int, input().split()))
    print(matchMathing(n, digits))
