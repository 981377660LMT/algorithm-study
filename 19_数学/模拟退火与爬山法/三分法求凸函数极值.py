"""三分法求凸函数极值"""

from typing import Callable


INF = int(4e18)


def minimize(fun: Callable[[int], int], lower: int, upper: int) -> int:
    """三分法求`严格凸函数fun`在`[lower,upper]`间的最小值"""
    res = INF
    while (upper - lower) >= 3:
        diff = upper - lower
        mid1 = lower + diff // 3
        mid2 = lower + 2 * diff // 3
        if fun(mid1) > fun(mid2):
            lower = mid1
        else:
            upper = mid2

    while lower <= upper:
        cand = fun(lower)
        res = cand if cand < res else res
        lower += 1

    return res


def maximize(fun: Callable[[int], int], lower: int, upper: int) -> int:
    """三分法求`严格凸函数fun`在`[lower,upper]`间的最大值"""
    res = -INF
    while (upper - lower) >= 3:
        diff = upper - lower
        mid1 = lower + diff // 3
        mid2 = lower + 2 * diff // 3
        if fun(mid1) < fun(mid2):
            lower = mid1
        else:
            upper = mid2

    while lower <= upper:
        cand = fun(lower)
        res = cand if cand > res else res
        lower += 1

    return res


def optimize(fun: Callable[[int], int], lower: int, upper: int, *, min: bool) -> int:
    """三分法求`严格凸函数fun`在`[lower,upper]`间的最值"""
    return minimize(fun, lower, upper) if min else maximize(fun, lower, upper)


if __name__ == "__main__":
    assert optimize(lambda x: x**2 + 2 * x, -1, 400, min=True) == -1
    assert optimize(lambda x: 0, -1, 400, min=True) == 0

    assert optimize(lambda x: x**2, -1, 40, min=False) == 1600
    assert optimize(lambda x: x**2, -50, 40, min=False) == 2500
