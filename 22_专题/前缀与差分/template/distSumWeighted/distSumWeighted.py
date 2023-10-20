"""
运送仓库货物模型.
"""


from typing import Callable, List
from itertools import accumulate


def distSumWeighted(positions: List[int], weights: List[int]) -> Callable[[int, int, int], int]:
    """
    数轴上按照顺序分布着n个点,每个点的位置为positions[i],权重为weights[i].
    点i到点j的距离定义为 `weights[i]*abs(positions[i]-positions[j])`.
    求区间[start,end)内的所有点到点to的距离之和.
    """
    preSum = [0] + list(accumulate(weights))
    preMul = [0] + list(accumulate(w * p for w, p in zip(weights, positions)))

    def _cal(start: int, end: int, to: int, onLeft: bool) -> int:
        if start >= end:
            return 0
        res1 = (preSum[end] - preSum[start]) * positions[to]
        res2 = preMul[end] - preMul[start]
        return (res1 - res2) if onLeft else (res2 - res1)

    def query(start: int, end: int, to: int) -> int:
        min_ = end if end < to else to
        res1 = _cal(start, min_, to, True)
        max_ = start if start > to else to
        res2 = _cal(max_, end, to, False)
        return res1 + res2

    return query


if __name__ == "__main__":
    # P3932 浮游大陆的68号岛
    # https://www.luogu.com.cn/problem/P3932
    import sys

    input = sys.stdin.readline
    MOD = 19260817
    n, q = map(int, input().split())
    positions = [0] + list(accumulate(map(int, input().split())))
    weights = list(map(int, input().split()))

    D = distSumWeighted(positions, weights)
    for _ in range(q):
        to, start, end = map(int, input().split())
        start -= 1
        to -= 1
        print(D(start, end, to) % MOD)
