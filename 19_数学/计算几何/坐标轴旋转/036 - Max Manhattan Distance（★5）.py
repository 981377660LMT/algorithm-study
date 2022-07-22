"""对每个查询 求第qi个点和其他点的最大曼哈顿距离

2≤N≤1e5

曼哈顿距离 abs(x1-x2) + abs(y1-y2) 
切比雪夫距离 max(abs(x1-x2), abs(y1-y2))
对于曼哈顿距离相同的点，他们分布在同一横纵截距且截距相同的直线上。
对于切比雪夫距离相同的点，他们分布在以原点为中心的正方形边界上。


技巧1:
坐标逆时针旋转45度,曼哈顿距离变为切比雪夫距离 
X' = x - y
Y' = x + y
应用: 每个点到其他点的最大曼哈顿距离

技巧2:
同理只需将每一个点(x,y)转化为((x+y)/2,(x-y)/2). 
新坐标系下的曼哈顿距离即为原坐标系下切比雪夫距离。
应用: 通过将曼哈顿距离的x,y排序求前缀和,可以找到一个点到其他点距离和最小(中位数)

例子:
有N个二位平面上的点,定义每一个点到其八连通的点的距离为1(切比雪夫距离)。
选一个点，使得剩下所有点到该点的距离之和最小，求出这个距离之和。
"""

import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)

N, Q = map(int, input().split())
points = []
px = []
py = []
for _ in range(N):
    x, y = map(int, input().split())
    points.append((x, y))
    px.append(x - y)
    py.append(x + y)

px.sort()
py.sort()

for _ in range(Q):
    q = int(input()) - 1
    a, b = points[q]
    x, y = a - b, a + b
    cand1, cand2 = abs(x - px[0]), abs(x - px[-1])
    cand3, cand4 = abs(y - py[0]), abs(y - py[-1])
    print(max(cand1, cand2, cand3, cand4))

