# https://maspypy.github.io/library/mod/range_freq_of_linear.hpp


from math import ceil


def range_freq_of_linear(L: int, R: int, a: int, b: int, mod: int, lo: int, hi: int) -> int:
    """
    count(x in [L, R) where (ax+b mod) in [lo, hi))
    """
    if lo >= hi:
        return 0
    assert 0 <= lo and lo < hi and hi <= mod
    x1 = floor_sum_of_linear(L, R, a, b - lo, mod)
    x2 = floor_sum_of_linear(L, R, a, b - hi, mod)
    return x1 - x2


def floor_sum_of_linear(L: int, R: int, a: int, b: int, mod: int) -> int:
    """
    ```
    sum((x * a + b) // mod for x in range(L, R))
    ```
    """
    if L >= R:
        return 0
    res = 0
    b += L * a
    n = R - L

    if b < 0:
        k = ceil(-b / mod)
        b += k * mod
        res -= n * k

    while n:
        q, a = a // mod, a % mod
        res += n * (n - 1) // 2 * q
        # res %= MOD
        if b >= mod:
            q, b = b // mod, b % mod
            res += n * q
            # res %= MOD
        n, b = (a * n + b) // mod, (a * n + b) % mod
        a, mod = mod, a

    return res
