# https://atcoder.jp/contests/typical90/tasks/typical90_ay
# abc051 - 选k个数和不超过p的方案数
# n个物品 买k个 求价格不超过P的方案数
# k<=N<=40

# !半分全列挙 时间复杂度 n*2^(n/2)
# 左侧排序, 右侧枚举+二分查找

from bisect import bisect_right
from itertools import combinations
from typing import List


def typicalShop(prices: List[int], k: int, limit: int) -> int:
    n = len(prices)
    left, right = prices[: n // 2], prices[n // 2 :]
    leftCounter = [[] for _ in range(k + 1)]
    for i in range(k + 1):
        for sub in combinations(left, i):
            leftCounter[i].append(sum(sub))
    for v in leftCounter:
        v.sort()

    res = 0
    for i in range(k + 1):
        for sub in combinations(right, i):
            sum_ = sum(sub)
            res += bisect_right(leftCounter[k - i], limit - sum_)
    return res


if __name__ == "__main__":
    import sys

    input = sys.stdin.readline

    n, k, p = map(int, input().split())
    prices = list(map(int, input().split()))
    print(typicalShop(prices, k, p))
