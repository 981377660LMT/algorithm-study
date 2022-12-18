# 2249. 统计圆内格点数目(一个圆)
# https://atcoder.jp/contests/abc191/editorial/611
# !-1e5 <= x, y <= 1e5 , 0 <= r <= 1e5
# 求圆内(包括圆周上)的整点数目 注意浮点数的精度问题

# !注意Decimal 需要 python3.8提交 pypy3会TLE
from decimal import Decimal as D


def circleLatticePoints(x: D, y: D, r: D) -> int:
    res = 0
    lower, upper = D.__ceil__(y - r), D.__floor__(y + r)
    for dy in range(lower, upper + 1):
        dx = D.sqrt(r * r - (dy - y) * (dy - y))
        res += D.__floor__(x + dx) - D.__ceil__(x - dx) + 1
    return res


if __name__ == "__main__":
    x, y, r = map(D, input().split())
    print(circleLatticePoints(x, y, r))
