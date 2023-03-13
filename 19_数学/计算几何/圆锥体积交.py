# 圆锥体积交/最强妙脆角
# 求两个距离为d的圆锥体的交的体积


from math import acosh, asin, pi


def solve(r: float, h: float, d: float) -> float:
    a = r - d / 2
    k = 1 - a / r
    res = h * r**2 / 3 * (pi / 2 - 2 * k * (1 - k**2) ** 0.5 - asin(k) + k**3 * acosh(1 / k))
    return res * 2


print(solve(*map(float, input().split())))
