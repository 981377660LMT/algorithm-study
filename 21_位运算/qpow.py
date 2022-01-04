def qpow(base: int, exp: int, mod: int) -> int:
    res = 1

    while exp:
        if exp & 1:
            res *= base
            res %= mod

        exp >>= 1
        base **= 2
        base %= mod

    return res


print(qpow(200, 300, int(1e9 + 7)))

from collections import defaultdict

d = defaultdict(int, {1: [2]})

print(d[1], d[2])
