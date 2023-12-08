# 等差数列交集/等差数列交点
# Intersection of Arithmetic Progressions

from typing import Tuple


def _extendedEuclid(a, b):
    x, y, xx, yy = 1, 0, 0, 1
    while b:
        q = a // b
        a, b = b, a % b
        x, xx = xx, x - q * xx
        y, yy = yy, y - q * yy
    return a, x, y


def _crt(a1, m1, a2, m2):
    g, p, q = _extendedEuclid(m1, m2)
    if a1 % g != a2 % g:
        return 0, -1
    m = m1 // g * m2
    p = (p % m + m) % m
    q = (q % m + m) % m
    return (p * a2 % m * (m1 // g) % m + q * a1 % m * (m2 // g) % m) % m, m


def intersect(a1: int, d1: int, a2: int, d2: int) -> Tuple[int, int]:
    """
    y = a1 + d1 * x
    y = a2 + d2 * x
    返回两个等差数列的交集(首项，公差).
    如果无交集，返回(0, 0).
    """
    x, m = _crt(a1 % d1, d1, a2 % d2, d2)
    if m == -1:
        return 0, 0
    st = max(a1, a2)
    a = st if st % m == x else st - st % m + x
    return a, m


if __name__ == "__main__":
    x, m = intersect(7, 9, 13, 12)
    print(x, m)  # 25 36
