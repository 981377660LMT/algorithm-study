# N-所有点对的曼哈顿距离之和
# https://atcoder.jp/contests/abc351/submissions/52863257
#
# 切比雪夫距离: max(|x1-x2|, |y1-y2|)
# 曼哈顿距离: |x1-x2| + |y1-y2|
# !转换关系： (x,y) <-> (y+x, y-x)

from typing import List


def manhattanDistSum(xs: List[int], ys: List[int]) -> int:
    """计算所有点对的曼哈顿距离之和."""
    res = 0
    xs, ys = sorted(xs), sorted(ys)
    res = 0
    curSum = 0
    for i, v in enumerate(xs):
        res += v * i - curSum
        curSum += v
    curSum = 0
    for i, v in enumerate(ys):
        res += v * i - curSum
        curSum += v
    return res


def chebyshevDistSum(xs: List[int], ys: List[int]) -> int:
    """计算所有点对的切比雪夫距离之和."""
    newXs, newYs = [], []
    for x, y in zip(xs, ys):
        newXs.append(y + x)
        newYs.append(y - x)
    return manhattanDistSum(newXs, newYs) >> 1


if __name__ == "__main__":
    # https://atcoder.jp/contests/abc351/tasks/abc351_e
    # 所有点对的切比雪夫距离之和->切比雪夫距离转曼哈顿距离
    n = int(input())
    xs, ys = [0] * n, [0] * n
    for i in range(n):
        xs[i], ys[i] = map(int, input().split())
    xs1, ys1 = [], []
    xs2, ys2 = [], []
    for x, y in zip(xs, ys):
        if (x + y) & 1:
            xs1.append(x)
            ys1.append(y)
        else:
            xs2.append(x)
            ys2.append(y)
    print(chebyshevDistSum(xs1, ys1) + chebyshevDistSum(xs2, ys2))
