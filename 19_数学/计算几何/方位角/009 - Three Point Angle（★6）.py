# 计算角度
# 二维平面上有N个不同的点
# 任选三个点ijk 求∠ijk的最大值 注意 0<=角度<=180°
# EPS=1e-7

# 3≤N≤2000

# !1. 计算方位角用math.atan2算 返回值的单位为弧度，取值范围为(-pi,pi]
# !2. 弧度转角度用math.degrees算
# 3. 思路是枚举每个顶点 计算出所有角度 然后再二分180度的另一半

from bisect import bisect_right
from math import degrees, atan2
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline

n = int(input())
points = []
for _ in range(n):
    points.append(tuple(map(int, input().split())))

res = 0
for x2, y2 in points:
    angles = []
    for x1, y1 in points:
        if x1 == x2 and y1 == y2:
            continue
        angles.append(degrees(atan2(y2 - y1, x2 - x1)))
    angles.sort()

    n = len(angles)
    # 寻找180度的另一半
    for a1 in angles:
        pos = bisect_right(angles, a1 + 180)
        if pos < n:
            cand = min(angles[pos] - a1, 360 - angles[pos] + a1)
            res = max(res, cand)
        if pos - 1 >= 0:
            cand = min(angles[pos - 1] - a1, 360 - angles[pos - 1] + a1)
            res = max(res, cand)

print(res)
