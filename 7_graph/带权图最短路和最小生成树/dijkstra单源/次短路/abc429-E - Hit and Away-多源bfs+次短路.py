# abc429-E - Hit and Away-次短路
# https://atcoder.jp/contests/abc429/tasks/abc429_e
#
# 给定一个连通无向图 G，有 N 个顶点、M 条边。每条边双向、花费 1。
# 每个顶点被标记为安全 S 或危险 D（字符串 S 的第 i 个字符是该点的类型）。保证安全点至少 2 个，危险点至少 1 个。
# 对于每个危险顶点 v，求从某个安全点出发，经过 v，最终到达另一个不同的安全点时所需时间的最小可能值。
#
# !由于边权为1 ，故可以从每个安全结点出发进行多源 BFS，在每个危险结点处维护距离最近的两个安全结点即可。

from collections import deque

INF = int(1e18)

N, M = map(int, input().split())
adjList = [[] for _ in range(N)]
for _ in range(M):
    u, v = map(int, input().split())
    u -= 1
    v -= 1
    adjList[u].append(v)
    adjList[v].append(u)
S = input()

dist1, dist2 = [INF] * N, [INF] * N
dist1From, dist2From = [-1] * N, [-1] * N


def relax(v: int, d: int, fid: int) -> bool:
    # 尝试把 (d,fid) 插入 v 的前两小（不同源）
    if dist1From[v] == fid:
        if dist1[v] > d:
            dist1[v] = d
            return True
        return False
    if dist2From[v] == fid:
        if dist2[v] > d:
            dist2[v] = d
            return True
        return False
    if dist1[v] > d:
        dist2[v], dist2From[v] = dist1[v], dist1From[v]
        dist1[v], dist1From[v] = d, fid
        return True
    if dist2[v] > d:
        dist2[v], dist2From[v] = d, fid
        return True
    return False


safe = [i for i, c in enumerate(S) if c == "S"]
q = deque()  # (cur, from_id, dist)
for v in safe:
    q.append((v, v, 0))

while q:
    cur, fid, d = q.popleft()
    if not relax(cur, d, fid):
        continue
    nd = d + 1
    for nxt in adjList[cur]:
        q.append((nxt, fid, nd))

res = []
for i, c in enumerate(S):
    if c == "D":
        res.append(str(dist1[i] + dist2[i]))
print("\n".join(res))
