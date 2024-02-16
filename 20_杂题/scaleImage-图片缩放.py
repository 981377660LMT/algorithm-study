# scaleImage-图片缩放

from math import gcd
from typing import Tuple


def scaleImage(width: int, height: int, x: int, y: int) -> Tuple[int, int]:
    gcd_ = gcd(x, y)  # !保证缩放到整数
    if gcd_ == 0:
        return 0, 0
    x //= gcd_
    y //= gcd_
    minRatio = min(width // x, height // y)
    return x * minRatio, y * minRatio


if __name__ == "__main__":
    # https://www.luogu.com.cn/problem/CF16C
    # 给出一矩形长 a 和宽 b，给出长宽比 x:y，要求`缩短`矩形的长和宽到整数.
    # 使它的长和宽之比等于 x:y，试求最终的长宽.
    # 如果无法缩放，输出 0,0.
    a, b, x, y = map(int, input().split())
    print(*scaleImage(a, b, x, y))
