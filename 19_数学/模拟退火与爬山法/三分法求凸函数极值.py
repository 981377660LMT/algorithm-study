from typing import Callable


INF = int(4e18)


def minimize(f: Callable[[int], int], lower: int, upper: int) -> int:
    """三分法求`凸函数f`在`[lower,upper]`间的最小值"""
    res = INF
    while (upper - lower) >= 3:
        diff = upper - lower
        mid1 = lower + diff // 3
        mid2 = lower + 2 * diff // 3
        if f(mid1) > f(mid2):
            lower = mid1
        else:
            upper = mid2

    while lower <= upper:
        res = min(res, f(lower))
        lower += 1

    return res


def maximize(f: Callable[[int], int], lower: int, upper: int) -> int:
    """三分法求`凸函数f`在`[lower,upper]`间的最大值"""
    res = -INF
    while (upper - lower) >= 3:
        diff = upper - lower
        mid1 = lower + diff // 3
        mid2 = lower + 2 * diff // 3
        if f(mid1) < f(mid2):
            lower = mid1
        else:
            upper = mid2

    while lower <= upper:
        res = max(res, f(lower))
        lower += 1

    return res


def optimize(f: Callable[[int], int], lower: int, upper: int, *, needMin: bool) -> int:
    """三分法求`凸函数f`在`[lower,upper]`间的最值"""
    return minimize(f, lower, upper) if needMin else maximize(f, lower, upper)


if __name__ == "__main__":
    assert optimize(lambda x: x**2 + 2 * x, -1, 400, needMin=True) == -1
    assert optimize(lambda x: 0, -1, 400, needMin=True) == 0

    assert optimize(lambda x: x**2, -1, 40, needMin=False) == 1600
    assert optimize(lambda x: x**2, -50, 40, needMin=False) == 2500
