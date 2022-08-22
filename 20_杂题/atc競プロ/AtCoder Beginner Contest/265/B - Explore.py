# 能否到达n-1房间
# !dont WA


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n, m, timeLimit = map(int, input().split())
costs = list(map(int, input().split()))  # ai 到 ai+1 需要的时间
buff = [0] * (n + 5)
for _ in range(m):
    x, y = map(int, input().split())
    x -= 1
    buff[x] += y

cur = timeLimit
for i in range(n - 1):
    cur -= costs[i]
    if cur <= 0:
        print("No")
        exit(0)
    cur += buff[i + 1]  # 注意这里是i+1

print("Yes")
