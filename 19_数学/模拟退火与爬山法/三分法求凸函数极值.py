"""
三分法求单峰函数的极值点.
更快的版本见: FibonacciSearch
"""

from typing import Callable


INF = int(4e18)


def minimize(fun: Callable[[int], int], left: int, right: int) -> int:
    """三分法求`严格凸函数fun`在`[left,right]`间的最小值"""
    res = INF
    while (right - left) >= 3:
        diff = (right - left) // 3
        mid1 = left + diff
        mid2 = right - diff
        if fun(mid1) > fun(mid2):
            left = mid1
        else:
            right = mid2

    while left <= right:
        cand = fun(left)
        res = cand if cand < res else res
        left += 1

    return res


def maximize(fun: Callable[[int], int], left: int, right: int) -> int:
    """三分法求`严格凸函数fun`在`[left,right]`间的最大值"""
    res = -INF
    while (right - left) >= 3:
        diff = (right - left) // 3
        mid1 = left + diff
        mid2 = right - diff
        if fun(mid1) < fun(mid2):
            left = mid1
        else:
            right = mid2

    while left <= right:
        cand = fun(left)
        res = cand if cand > res else res
        left += 1

    return res


def optimize(fun: Callable[[int], int], left: int, right: int, *, min: bool) -> int:
    """三分法求`严格凸函数fun`在`[left,right]`间的最值"""
    return minimize(fun, left, right) if min else maximize(fun, left, right)


if __name__ == "__main__":
    assert optimize(lambda x: x**2 + 2 * x, -1, 400, min=True) == -1
    assert optimize(lambda x: 0, -1, 400, min=True) == 0

    assert optimize(lambda x: x**2, -1, 40, min=False) == 1600
    assert optimize(lambda x: x**2, -50, 40, min=False) == 2500
