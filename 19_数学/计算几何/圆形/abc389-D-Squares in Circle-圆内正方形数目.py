# abc389-D-Squares in Circle (圆内正方形数目)
# https://atcoder.jp/contests/abc389/tasks/abc389_d
# 题意：有一个无限的平铺大小为1×1的正方形的平面，你在一个正方形中心画了一个半径为R的圆，问有多少个正方形被包含。
# R<=1e6
#
# !滑动窗口解法.
# !对固定的纵坐标，利用滑动窗口找到横坐标的边界，然后计算这个纵坐标下的正方形数目。


def squaresInCircle(r: int) -> int:
    d = 2 * r
    x = -1
    y = d - 1
    res = 0
    while y > 0:
        while y * y + (2 * x + 1) ** 2 <= d * d:
            x += 1
        if y > 1:
            res += 4 * x - 2
        else:
            res += 2 * x - 1
        y -= 2
    return res


if __name__ == "__main__":
    r = int(input())
    print(squaresInCircle(r))
