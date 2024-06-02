# 黄金比例搜索/黄金搜索 求单峰凸函数极小值

from typing import Callable, Tuple


def goldenSearch(
    fun: Callable[[float], float], min: float, max: float, iter=50
) -> Tuple[float, float]:
    invPhi = (5**0.5 - 1.0) * 0.5
    invPhi2 = invPhi * invPhi
    x1, x4 = min, max
    x2 = x1 + invPhi2 * (x4 - x1)
    x3 = x1 + invPhi * (x4 - x1)
    y2, y3 = fun(x2), fun(x3)
    for _ in range(iter):
        if y2 < y3:
            x4, x3, y3 = x3, x2, y2
            x2 = x1 + invPhi2 * (x4 - x1)
            y2 = fun(x2)
        else:
            x1, x2, y2 = x2, x3, y3
            x3 = x1 + invPhi * (x4 - x1)
            y3 = fun(x3)
    return (y2, x2) if y2 < y3 else (y3, x3)
