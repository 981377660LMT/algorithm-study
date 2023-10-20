from math import log2


def rangeSum(start: int, end: int) -> int:
    """区间 [start, end) 的和."""
    if start >= end:
        return 0
    return (end - start) * (start + end - 1) // 2


def rangeSquareSum(start: int, end: int) -> int:
    """区间 [start, end) 的平方和."""
    if start >= end:
        return 0
    tmp1 = end * (end - 1) * (2 * end - 1) // 6
    tmp2 = start * (start - 1) * (2 * start - 1) // 6
    return tmp1 - tmp2


def rangeCubeSum(start: int, end: int) -> int:
    """区间 [start, end) 的立方和."""
    if start >= end:
        return 0
    tmp1 = end * (end - 1) // 2
    tmp2 = start * (start - 1) // 2
    return tmp1 * tmp1 - tmp2 * tmp2


def rangeXorSum(start: int, end: int) -> int:
    """区间 [start, end) 的异或和."""
    if start >= end:
        return 0

    def preXor(upper: int) -> int:
        """[0, upper]内所有数的异或 upper>=0"""
        mod = upper % 4
        if mod == 0:
            return upper
        elif mod == 1:
            return 1
        elif mod == 2:
            return upper + 1
        return 0

    return preXor(end - 1) ^ preXor(start - 1)


def rangePow2Sum(start: int, end: int, mod=int(1e9 + 7)) -> int:
    return (pow(2, end, mod) - pow(2, start, mod)) % mod


def rangePowKSum(start: int, end: int, k: int, mod=int(1e9 + 7)) -> int:
    """
    区间 [start,end) k次幂之和模mod.
    powerSum/powSum/rangePowSum
    """
    if start >= end:
        return 0
    if mod == 1:
        return 0

    def cal(n: int) -> int:
        sum_, p = 1, k
        start = int(log2(n)) - 1
        for d in range(start, -1, -1):
            sum_ *= p + 1
            p *= p
            if (n >> d) & 1:
                sum_ += p
                p *= k
            sum_ %= mod
            p %= mod
        return sum_

    return cal(end) - cal(start)


if __name__ == "__main__":
    from functools import reduce
    import operator

    for start in range(1, 10 + 1):
        for end in range(start, 10 + 1):
            assert rangeSum(start, end + 1) == sum(range(start, end + 1))
            assert rangeSquareSum(start, end + 1) == sum(v * v for v in range(start, end + 1))
            assert rangeCubeSum(start, end + 1) == sum(v * v * v for v in range(start, end + 1))
            assert rangeXorSum(start, end + 1) == reduce(operator.xor, range(start, end + 1))
            assert rangePowKSum(start, end + 1, 2) == sum(2**i for i in range(start, end + 1))
            assert rangePowKSum(start, end + 1, 3) == sum(3**i for i in range(start, end + 1))
            assert rangePow2Sum(start, end + 1) == sum(2**i for i in range(start, end + 1))
