from math import e
from typing import List, Tuple
from angleSort import angleArgSort


# !最大化向量长度
# 给定一组二维平面上的点（向量），从中选择一个子集，使得这些向量的和的范数（norm）最大化。
# 返回 `最大化的范数` 和 `选择的向量的下标`。
# https://codeforces.com/contest/1841/problem/F
# https://atcoder.jp/contests/abc139/tasks/abc139_f


def maxNormSum(points: List[Tuple[int, int]]) -> Tuple[int, List[int]]:
    order = angleArgSort(points)

    n = len(points)
    if n == 0:
        return 0, []

    res = 0
    lr = (0, 0)
    left, right = 0, 1
    c = points[0]

    def calc() -> int:
        x, y = c
        return x * x + y * y

    cand = calc()
    if cand > res:
        res = cand
        lr = (left, right)

    while left < n:
        a = points[left]
        b = points[right % n]
        if right - left < n and (
            a[0] * b[1] - a[1] * b[0] > 0
            or (a[0] * b[1] - a[1] * b[0] == 0 and a[0] * b[0] + a[1] * b[1] > 0)
        ):
            c = (c[0] + b[0], c[1] + b[1])
            right += 1
            cand = calc()
            if cand > res:
                res = cand
                lr = (left, right)
        else:
            c = (c[0] - a[0], c[1] - a[1])
            left += 1
            cand = calc()
            if cand > res:
                res = cand
                lr = (left, right)

    ids = []
    for i in range(lr[0], lr[1]):
        ids.append(order[i % n])
    return res, ids


if __name__ == "__main__":
    # https://atcoder.jp/contests/abc139/tasks/abc139_f
    # 给出n个二元组(x,y)。最初位于原点(0,0)，每次可以从这n个二元组中挑出一个，
    # 然后将当前的坐标(X,Y)变为(X+x,Y+y)，每个二元组只能被选一次。
    # !选出一些二元组，使得按照这些二元组移动后与原点的欧几里得距离最大。求这个距离。
    # n<=100
    def engines() -> None:
        n = int(input())
        points = [tuple(map(int, input().split())) for _ in range(n)]
        val, ids = maxNormSum(points)
        x, y = 0, 0
        for i in ids:
            x += points[i][0]
            y += points[i][1]
        assert val == x * x + y * y
        print(val**0.5)

    engines()
    points = [(0, 0), (1, 0), (0, 1), (1, 1), (2, 2)]
    assert maxNormSum(points) == (32, [1, 2, 3, 4])  # type: ignore
