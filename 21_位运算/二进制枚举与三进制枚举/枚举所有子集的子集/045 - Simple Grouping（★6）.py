# # 平面上n个点分k组
# # 最小化组内的距离最大值 输出距离的这个最小值的平方
# # 2<=k<n<=15

# # 2152.覆盖点所需要的最少直线数-子集状压dp
# # !枚举子集的子集dp k*3^n
# # dp[第i组][使用的点集j]


from itertools import combinations
import sys


input = sys.stdin.readline


n, k = map(int, input().split())
points = []
for _ in range(n):
    x, y = map(int, input().split())
    points.append((x, y))

dist = [[0] * n for _ in range(n)]
for i, j in combinations(range(n), 2):
    dx, dy = points[i][0] - points[j][0], points[i][1] - points[j][1]
    cur = dx * dx + dy * dy
    dist[i][j] = cur
    dist[j][i] = cur

# !1. 预处理
submax = [0] * (2 ** n)
for state in range(2 ** n):
    cur = [i for i in range(n) if state & (1 << i)]
    if len(cur) >= 2:
        submax[state] = max(dist[i][j] for i, j in combinations(cur, 2))

# !2. 初始化 (分为1组)
dp = submax[:]

# !3. 转移
for _ in range(1, k):
    # !这里把1<<63(9223372036854775808) 改小成 1<<63-1(9223372036854775807) 快了700ms pypy3 (碰到超过2^63-1的树就会变慢)
    ndp = [9223372036854775807] * (2 ** n)
    for state in range(2 ** n):
        g1, g2 = state, 0
        while g1:
            ndp[state] = min(ndp[state], max(submax[g1], dp[g2]))
            g1 = (g1 - 1) & state
            g2 = state ^ g1
    dp = ndp

print(dp[-1])
print(2 ** 63)

