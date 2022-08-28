# 打地鼠
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n = int(input())
goods = []

for _ in range(n):
    time, pos, score = map(int, input().split())
    goods.append((time, pos, score))

# dp[时间][位置] = 分数
dp = [-INF] * 5
dp[0] = 0
for i in range(n):
    ndp = [-INF] * 5
    dist = goods[i][0] - (goods[i - 1][0] if i > 0 else 0)
    for pre in range(5):
        for cur in range(5):
            if abs(pre - cur) <= dist:
                ndp[cur] = max(ndp[cur], dp[pre] + (goods[i][2] if cur == goods[i][1] else 0))
    dp = ndp

print(max(dp))
