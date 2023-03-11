# https://maspypy.github.io/library/mod/floor_sum_of_linear.hpp
# !sum((x * a + b) // mod for x in range(L, R))

from math import ceil

MOD = int(1e9 + 7)


def floor_sum_of_linear(L: int, R: int, a: int, b: int, div: int) -> int:
    """
    ```
    sum((x * a + b) // div for x in range(L, R))
    ```
    """
    if L >= R:
        return 0
    res = 0
    b += L * a
    n = R - L

    if b < 0:
        k = ceil(-b / div)
        b += k * div
        res -= n * k

    while n:
        q, a = a // div, a % div
        res += n * (n - 1) // 2 * q
        # res %= MOD
        if b >= div:
            q, b = b // div, b % div
            res += n * q
            # res %= MOD
        n, b = (a * n + b) // div, (a * n + b) % div
        a, div = div, a

    return res


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    N = int(input())
    M = int(input())
    A = list(map(int, input().split()))
    B = list(map(int, input().split()))
    res = 0
    for a in A:
        for b in B:
            res += floor_sum_of_linear(1, b + 1, a, 0, b)
    res *= 2
    res %= MOD
    print(res)
