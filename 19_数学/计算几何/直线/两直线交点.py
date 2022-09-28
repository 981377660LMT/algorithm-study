# 两直线交点


from typing import Tuple, Union


def solve(
    a1: float, b1: float, c1: float, a2: float, b2: float, c2: float
) -> Union[Tuple[float, float], Tuple[None, None]]:
    """直线表达式为ax+by=c  求两直线交点"""
    if a1 * b2 == a2 * b1:  # 平行
        return None, None
    x = (c1 * b2 - c2 * b1) / (a1 * b2 - a2 * b1)
    y = (c2 * a1 - c1 * a2) / (a1 * b2 - a2 * b1)
    return x, y
