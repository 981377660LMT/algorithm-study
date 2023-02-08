# https://medium.com/hard-mode/coding-challenge-rectangle-rotation-10e2a2416ef3

# 长、宽分别为a和b的长方形（a、b都是正整数），中心为(0, 0)点。
# 将其逆时针旋转45°后，问长方形内和边界上的整点共有几个。
# 整点的定义是，坐标(x, y)的点，其中x和y都是整数。

# 求出直线方程 如果截距不是整数 那么边界就没有整点
# 如果是整数 那么那条边的整点数就是 gcd(dx,dy)+1

from math import sqrt


def count(a: int, b: int) -> int:
    ...


print(count(6, 4))
