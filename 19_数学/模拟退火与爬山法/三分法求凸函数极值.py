from typing import Callable


def optimize(
    f: Callable[[int], int],
    lower: int,
    upper: int,
) -> int:
    """三分法求`凸函数f`在`[lower,upper]`间的极值"""
    left, right = lower, upper
    while (right - left) >= 3:
        diff = right - left
        mid1 = left + diff // 3
        mid2 = left + 2 * diff // 3
        if f(mid1) < f(mid2):
            right = mid2 - 1
        else:
            left = mid1 + 1
    return f(left)


print(optimize(lambda x: x**2, -1, 4))
