# G - Push Simultaneously
# https://atcoder.jp/contests/abc401/tasks/abc401_g
# 在一个平面直角坐标系中有N个人和N个按钮，每个人和按钮都有具体的(x,y)坐标。
# 请你为每个人分配一个按钮，然后每个人会以1m/s的速度朝着所分配到的按钮走去。
# 当所有人都走到所分配的按钮位置时，记当前时间为t。
# 请在所有的分配方案当中，找出t的最小值。
#
# !二分图完美匹配，也可以匈牙利算法.

import sys
from math import sqrt
from scipy.optimize import linear_sum_assignment

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")


if __name__ == "__main__":
    N = int(input())
    sx, sy = [0] * N, [0] * N
    tx, ty = [0] * N, [0] * N
    for i in range(N):
        sx[i], sy[i] = map(int, input().split())
    for i in range(N):
        tx[i], ty[i] = map(int, input().split())

    info = []
    for i in range(N):
        for j in range(N):
            dx, dy = abs(sx[i] - tx[j]), abs(sy[i] - ty[j])
            info.append((dx * dx + dy * dy, i, j))
    info.sort()

    # 二分info下标

    def check(mid: int) -> bool:
        costMatrix = [[0] * N for _ in range(N)]
        for i in range(mid):
            costMatrix[info[i][1]][info[i][2]] = 1
        rowIndex, colIndex = linear_sum_assignment(costMatrix, maximize=True)
        cost = sum(costMatrix[r][c] for r, c in zip(rowIndex, colIndex))
        return cost == N

    left, right = 0, len(info)
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            right = mid - 1
        else:
            left = mid + 1

    d = info[left - 1][0]
    print(sqrt(d))
