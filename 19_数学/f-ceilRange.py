# https://maspypy.github.io/library/enumerate/ceil_range.hpp
# f(l,r,q)：ceil が q になる範囲が [l, r)
# q == 1 のときには r == infty<T>

from math import ceil
from typing import List, Tuple

INF = int(1e18)


def ceilRange(n: int) -> List[Tuple[int, int, int]]:
    """
    将 [1,n) 内的数分成O(2*sqrt(n))段, 每段内的 ceil(n/i) 相同

    Args:
        n (int): n>=1

    Returns:
        List[Tuple[int,int,int]]:
        每个元素为(left,right,div)
        表示 left <= i < right 内的 ceil(n/i) == div
    """
    res = []
    l, r, q = n, INF, 1
    while True:
        res.append((l, r, q))
        if q == n:
            break
        r = l
        q = ceil(n / (l - 1))
        l = ceil(n / q)
    return res[::-1][:-1]


if __name__ == "__main__":
    print(ceilRange(10))
    # [(1, 2, 10), (2, 3, 5), (3, 4, 4), (4, 5, 3), (5, 10, 2)]
