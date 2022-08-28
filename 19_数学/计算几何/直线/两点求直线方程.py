from fractions import Fraction
from math import gcd
from typing import Tuple, Union


def calSlopeInterceptForm(
    x1: int, y1: int, x2: int, y2: int
) -> Union[Tuple[Fraction, Fraction], Tuple[None, int]]:
    """求出直线方程的斜截式`y=kx+b 或 x=b`的斜率和截距

    Returns:
        Tuple[Fraction, Fraction] | Tuple[None, int]: 斜率和截距
        当斜率为None时,截距表示在x轴上的截距,直线方程为 x=截距;
        当斜率不为None时,截距表示在y轴上的截距,直线方程为 y=斜率*x+截距
    """
    if x1 == x2:
        return None, x1
    slope = Fraction(y2 - y1, x2 - x1)
    intercept = y1 - slope * x1
    return slope, intercept


def calGeneralForm(x1: int, y1: int, x2: int, y2: int) -> Tuple[int, int, int]:
    """求出直线方程的一般式`ax+by+c=0`的系数"""
    if x1 == x2:
        return 1, 0, -x1
    dy, dx = y2 - y1, x1 - x2
    gcd_ = gcd(dy, dx)
    dy, dx = dy // gcd_, dx // gcd_
    c = dx * y1 - dy * x1
    if dx > 0:
        dx, dy, c = -dx, -dy, -c
    return dy, -dx, c
