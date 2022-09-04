# https://www.acwing.com/problem/content/364/

# 给定 n 个区间 [ai,bi] 和 n 个整数 ci。
# !你需要构造一个整数集合 Z，使得 ∀i∈[1,n]，Z 中满足 ai≤x≤bi 的整数 x 不少于 ci 个。
# !求这样的整数集合 Z 最少包含多少个数。
# 1≤n≤50000,
# !0≤ai,bi≤50000

# !前缀和preSum[i]表示[0,i]选择了多少个数 题目要求preSum[50000]的最小值
# 最小值=>求最长路
# !即求在约束下0到50000的最长路


# 所有的限制要找全：
# !1. Si >= Si-1
# !2. Si - Si-1 <= 1
# !3. Sb - Sa-1 >= c


from collections import defaultdict, deque
from typing import List, Mapping

INF = int(1e18)


def spfa(n: int, adjMap: Mapping[int, Mapping[int, int]], start: int) -> List[int]:
    """spfa求单源最长路 图中无正环"""
    dist = [-INF] * n
    dist[start] = 0
    queue = deque([start])
    inQueue = [False] * n
    inQueue[start] = True

    while queue:
        cur = queue.popleft()
        inQueue[cur] = False
        for next in adjMap[cur]:
            weight = adjMap[cur][next]
            cand = dist[cur] + weight
            if cand > dist[next]:  # 最长路
                dist[next] = cand
                if not inQueue[next]:
                    inQueue[next] = True
                    if queue and dist[next] > dist[queue[0]]:
                        queue.appendleft(next)
                    else:
                        queue.append(next)
    return dist


MAX = 50005
n = int(input())
adjMap = defaultdict(lambda: defaultdict(lambda: -INF))  # 最长路
for _ in range(n):
    u, v, w = map(int, input().split())
    u, v = u + 1, v + 1  # 从1开始
    adjMap[u - 1][v] = max(adjMap[u - 1][v], w)  # Sv - Su-1 >= w

# !前缀和满足的约束
for i in range(1, MAX + 1):
    adjMap[i - 1][i] = max(adjMap[i - 1][i], 0)  # Si - Si-1 >= 0
    adjMap[i][i - 1] = max(adjMap[i][i - 1], -1)  # Si-1 - Si >= -1


dist = spfa(MAX + 1, adjMap, 0)
print(dist[-1])
