# D - Pedometer
# abc367-环上前缀和/环形前缀和
# 给定环上n个点的间距。
# 求二元组(s,t)对数，满足从s沿着顺时针方向走到t的距离是k的倍数。


from bisect import bisect_left
from itertools import accumulate
from typing import List


def solve(dists: List[int], k: int) -> int:
    n = len(dists)
    dists2 = dists[:] + dists[:]
    preSum = [0] + list(accumulate(dists2))
    groups = [[] for _ in range(k)]
    for i in range(2 * n):
        groups[preSum[i] % k].append(i)

    res = 0
    for i in range(n):
        group = groups[preSum[i] % k]
        res += bisect_left(group, i + n) - bisect_left(group, i + 1)  # [i+1, i+n)
    return res


if __name__ == "__main__":
    n, k = map(int, input().split())
    dists = list(map(int, input().split()))
    print(solve(dists, k))
