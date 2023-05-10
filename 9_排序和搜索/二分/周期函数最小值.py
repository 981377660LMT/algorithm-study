# https://sotanishy.github.io/cp-library-cpp/misc/min_periodic_function.hpp
# 周期函数最小值


from typing import Callable


def minimizePeriodicFunction(n: int, f: Callable[[int], int]) -> int:
    a, b, c = 0, n, 2 * n
    while c - a > 2:
        l, r = (a + b) >> 1, (b + c + 1) >> 1
        if f(l) < f(b):
            c, b = b, l
        elif f(b) > f(r):
            a, b = b, r
        else:
            a, c = l, r
    return f(b)


if __name__ == "__main__":
    n = 1
    f = lambda x: x % 300 + 12
    print(minimizePeriodicFunction(n, f))
