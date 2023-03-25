# BinaryGcd
# 二进制gcd: 在python中反而更慢了, 内置的 math.gcd 是最快的


def binaryGcd(a: int, b: int) -> int:
    x, y = a if a >= 0 else -a, b if b >= 0 else -b
    if x == 0 or y == 0:
        return x + y
    # trailingZeros = lambda x: (x&-x).bit_length()-1
    n = (x & -x).bit_length() - 1
    m = (y & -y).bit_length() - 1
    x >>= n
    y >>= m
    while x != y:
        d = (x - y).bit_length() - 1
        f = x > y
        c = x if f else y
        if not f:
            y = x
        x = (c - y) >> d
    min_ = n if n < m else m
    return x << min_


if __name__ == "__main__":
    # time it
    import timeit

    print(
        timeit.timeit(
            "binaryGcd(123456789,987654321)",
            setup="from __main__ import binaryGcd",
            number=int(1e6),
        )
    )

    def naiveGcd(a: int, b: int) -> int:
        while a != 0:
            a, b = b % a, a
        return b

    print(
        timeit.timeit(
            "naiveGcd(123456789,987654321)",
            setup="from __main__ import naiveGcd",
            number=int(1e6),
        )
    )

    print(
        timeit.timeit(
            "gcd(123456789,987654321)",
            setup="from math import gcd",
            number=int(1e6),
        )
    )
