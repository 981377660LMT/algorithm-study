# https://maspypy.github.io/library/mod/min_of_linear_segments.hpp


from math import ceil, gcd
from typing import List, Tuple


def max_of_linear_segments(a: int, b: int, mod: int) -> Tuple[List[int], List[int]]:
    """
    `ax + b (x>=0)` が最大となるところの情報を返す。
    prefix max を更新する x 全体が、等差数列の和集合。

    次を返す:

    ・等差数列の境界となる x_0, x_1, ..., x_n
    ・各境界の間での交差 dx_0, ..., dx_{n-1}
    """
    if a != 0:
        a = mod - a
    b = mod - 1 - b
    return _min_of_linear_segments(a, b, mod)


def _min_of_linear_segments(a: int, b: int, mod: int) -> Tuple[List[int], List[int]]:
    """
    `ax + b (x>=0)` が最小となるところの情報を返す。
    prefix min を更新する x 全体が、等差数列の和集合。

    次を返す:

    ・等差数列の境界となる x_0, x_1, ..., x_n
    ・各境界の間での交差 dx_0, ..., dx_{n-1}
    """
    assert 0 <= a < mod
    assert 0 <= b < mod
    X, DX = [0], []
    g = gcd(a, mod)
    a, b, mod = a // g, b // g, mod // g
    # p/q <= (mod-a)/mod <= r/s
    p, q, r, s = 0, 1, 1, 1
    det_l, det_r = mod - a, a
    x, y = 0, b

    while y:
        # upd r/s
        k = det_r // det_l
        det_r %= det_l
        if det_r == 0:
            k -= 1
            det_r = det_l
        r += k * p
        s += k * q
        while True:
            k = ceil((det_l - y) / det_r)
            if k < 0:
                k = 0
            if det_l - k * det_r <= 0:
                break
            det_l -= k * det_r
            p += k * r
            q += k * s
            # p/q <= a/mod
            # (aq - pmod) = det_l を y から引く
            k = y // det_l
            y -= k * det_l
            x += q * k
            X.append(x)
            DX.append(q)
        k = det_l // det_r
        det_l -= k * det_r
        p += k * r
        q += k * s
        # assert min(p, q, r, s) >= 0
    return X, DX


print(max_of_linear_segments(1, 0, 5))
