# 给定坐标(x, y)问绕原点逆时针旋转d角度后的坐标。
# 坐标为a'= a * cos d - y * sin d, g = a * sin d ＋ y * cosd，可以用各种方法(诱导公式/旋转矩阵)推:
# !向量旋转

from math import cos, radians, sin
import sys
import os
from typing import Tuple

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def rotate(x: int, y: int, *, deg: int) -> Tuple[float, float]:
    """(x, y)逆时针旋转deg度(角度制)后的坐标"""
    rad = radians(deg)
    return x * cos(rad) - y * sin(rad), x * sin(rad) + y * cos(rad)


def main() -> None:
    x1, y1, deg = map(int, input().split())
    x2, y2 = rotate(x1, y1, deg=deg)
    print(x2, y2)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
