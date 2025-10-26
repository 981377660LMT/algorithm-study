# abc429-E - Hit and Away-次短路
# https://atcoder.jp/contests/abc429/tasks/abc429_e
#
# 给定一个连通无向图 G，有 N 个顶点、M 条边。每条边双向、花费 1。
# 每个顶点被标记为安全 S 或危险 D（字符串 S 的第 i 个字符是该点的类型）。保证安全点至少 2 个，危险点至少 1 个。
# 对每个危险点 v，求最小的 dist(v, Sa) + dist(v, Sb)，其中 Sa、Sb 是两个不同的安全点。
#
# 思路：
# - 将所有安全点按顶点编号第 k 位的 0/1 分成两组 c1/c2；
# - 对每个 k，分别对 c1/c2 做多源 BFS，得到 d1/d2；
# - 对每个危险点 v，用 d1[v]+d2[v] 更新答案；
# - k 枚举至 ceil(log2 N) 即可，保证最近两名安全点会被某个 k 分到不同组。

import sys
from collections import deque

input = sys.stdin.readline
n, m = map(int, input().split())
g = [[] for _ in range(n)]
for _ in range(m):
    u, v = map(int, input().split())
    u, v = u - 1, v - 1
    g[u].append(v)
    g[v].append(u)
s = input()
ok, ng = [], []
for i in range(n):
    if s[i] == "S":
        ok.append(i)
    else:
        ng.append(i)

INF = 10**9


def f(starts):
    d = [INF] * n
    q = deque()
    for x in starts:
        d[x] = 0
        q.append(x)
    while q:
        x = q.popleft()
        for i in g[x]:
            if d[i] != INF:
                continue
            d[i] = d[x] + 1
            q.append(i)
    return d


res = [INF] * len(ng)
B = max(1, (n - 1).bit_length())
for k in range(B):
    c1, c2 = [], []
    for x in ok:
        if (x >> k) & 1:
            c2.append(x)
        else:
            c1.append(x)
    d1, d2 = f(c1), f(c2)
    for i in range(len(ng)):
        x = ng[i]
        res[i] = min(res[i], d1[x] + d2[x])

print("\n".join(map(str, res)))
