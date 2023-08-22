# FibonacciSearch 斐波那契搜索求凸函数极值，比三分法稍快一些
# !寻找[left,right]中的一个极值点,不要求单峰性质


from typing import Callable, Tuple


INF = int(1e18)


def minimize(fun: Callable[[int], int], left: int, right: int) -> Tuple[int, int]:
    """斐波那契搜索求`凸函数fun`在`[left,right]`间的`(极小值点,极小值)`."""
    return fibonacciSearch(fun, left, right, min=True)


def maximize(fun: Callable[[int], int], left: int, right: int) -> Tuple[int, int]:
    """斐波那契搜索求`凸函数fun`在`[left,right]`间的`(极大值点,极大值)`."""
    return fibonacciSearch(fun, left, right, min=False)


def fibonacciSearch(f: Callable[[int], int], left: int, right: int, min: bool) -> Tuple[int, int]:
    """斐波那契搜索寻找[left,right]中的一个极值点,不要求单峰性质.
    Args:
        f: 目标函数.
        left: 搜索区间左端点(包含).
        right: 搜索区间右端点(包含).
        min: 是否寻找最小值.
    Returns:
        极值点的横坐标x和纵坐标f(x).
    """
    assert left <= right

    a, b, c, d = left, left + 1, left + 2, left + 3
    n = 0
    while d < right:
        b = c
        c = d
        d = b + c - a
        n += 1

    def get(i: int) -> int:
        if right < i:
            return INF
        return f(i) if min else -f(i)

    ya, yb, yc, yd = get(a), get(b), get(c), get(d)
    for _ in range(n):
        if yb < yc:
            d = c
            c = b
            b = a + d - c
            yd = yc
            yc = yb
            yb = get(b)
        else:
            a = b
            b = c
            c = a + d - b
            ya = yb
            yb = yc
            yc = get(c)

    x = a
    y = ya
    if yb < y:
        x = b
        y = yb
    if yc < y:
        x = c
        y = yc
    if yd < y:
        x = d
        y = yd

    return (x, y) if min else (x, -y)


if __name__ == "__main__":

    def f(x: int) -> int:
        return x**2

    print(fibonacciSearch(f, 0, 100, True))
    print(fibonacciSearch(f, 0, 100, False))
