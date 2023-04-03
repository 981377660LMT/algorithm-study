# 斐波那契数列第k项


MOD = int(1e9 + 7)


def kthFibonacci(k: int) -> int:
    """斐波那契数列前k(0-indexed)项:0,1,1,2,3,5..."""
    f, res = (0, 1), (1, 0)
    while k:
        a, b = f
        c, d = res
        if k & 1:
            res = ((a * c + b * d) % MOD, (b * c + (a + b) * d) % MOD)
        f = ((a * a + b * b) % MOD, (b * (a + a + b)) % MOD)
        k >>= 1
    return res[1]


res = [kthFibonacci(i) for i in range(10)]
assert res == [0, 1, 1, 2, 3, 5, 8, 13, 21, 34]
