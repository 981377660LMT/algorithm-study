# https://mugen1337.github.io/procon/Math/HighlyCompositeNumber.hpp
# HighlyCompositeNumber
# 不超过n的约数最多的自然数


from heapq import heappush, heappop
from typing import List, Tuple

P = [
    2,
    3,
    5,
    7,
    11,
    13,
    17,
    19,
    23,
    29,
    31,
    37,
    41,
    43,
    47,
    53,
    59,
    61,
    67,
    71,
    73,
    79,
    83,
    89,
    97,
    101,
    103,
    107,
    109,
    113,
    127,
    131,
    137,
    139,
]


def highly_composite_number(n: int) -> Tuple[int, int]:
    """
    返回不超过n的自然数中`约数最多的自然数`和`约数个数`.
    如果存在多个, 返回任意一个.
    1<=n<=1e18
    """
    if n == 1:
        return 1, 1
    pq = [(-2, -2, [-1])]  # number, divisor count, exp
    res = (2, 2, [1])
    while pq:
        num, divs, e = heappop(pq)
        num, divs = -num, -divs
        if res[1] < divs:
            res = (num, divs, e)
        m = len(e)
        if e[0] == -1:
            newNum = num * P[m]
            if newNum <= n:
                newE = e[:]
                newE.append(-1)
                heappush(pq, (-newNum, -divs * 2, newE))
        e0 = -e[0]
        for i in range(m):
            if e0 > -e[i]:
                break
            num *= P[i]
            if num > n:
                break
            e[i] -= 1
            divs //= e0 + 1
            divs *= e0 + 2
            heappush(pq, (-num, -divs, e[:]))
    return res[0], res[1]


if __name__ == "__main__":
    print((highly_composite_number(int(1e18))))

    def getFactors(n: int) -> List[int]:
        """n 的所有因数 O(sqrt(n))"""
        if n <= 0:
            return []
        small, big = [], []
        upper = int(n**0.5) + 1
        for i in range(1, upper):
            if n % i == 0:
                small.append(i)
                if i != n // i:
                    big.append(n // i)
        return small + big[::-1]

    preRes = [(0, 0)]
    for i in range(1, 20000):
        factors = getFactors(i)
        if len(factors) >= preRes[-1][1]:
            preRes.append((i, len(factors)))
        else:
            preRes.append(preRes[-1])
        assert preRes[-1][1] == highly_composite_number(i)[1], i
    print("ok")
