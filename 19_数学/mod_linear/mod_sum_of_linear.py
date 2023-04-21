## mod sum of linear

from math import ceil


def mod_sum_of_linear(L: int, R: int, a: int, b: int, mod: int) -> int:
    """
    ```
    sum((a * x + b) % mod for x in range(L, R))
    ```
    """
    s = a * L + b
    t = a * (R - 1) + b
    sum = (s + t) * (R - L) // 2
    fsum = floor_sum_of_linear(L, R, a, b, mod)
    return sum - fsum * mod


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
