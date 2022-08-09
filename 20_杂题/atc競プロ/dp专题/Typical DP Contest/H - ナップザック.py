# 重量不能超过w 颜色种数不能超过c 求最大价值
# n<=100 w<=1e4 c<=50

# !按照颜色组遍历 分组背包

from collections import defaultdict
import sys

INF = int(1e9)

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")


n, w, c = map(int, input().split())
goods = defaultdict(list)

for _ in range(n):
    weight, score, color = map(int, input().split())
    goods[color].append((weight, score))

dp = [[0] * (w + 1) for _ in range(c + 1)]  # dp[i][j] 表示i种颜色重量为j的最大价值

"""遍历每个颜色组"""
keys = sorted(goods)
for color in keys:
    for i in range(min(color, c), 0, -1):
        ndp = dp[i - 1][::]
        for weight, score in goods[color]:  # 物品
            for j in range(w, weight - 1, -1):  # 容量
                ndp[j] = max(ndp[j], ndp[j - weight] + score)
        for j in range(w + 1):
            dp[i][j] = max(dp[i][j], ndp[j])

res = 0
for row in dp:
    res = max(res, max(row))
print(res)


# TODO
