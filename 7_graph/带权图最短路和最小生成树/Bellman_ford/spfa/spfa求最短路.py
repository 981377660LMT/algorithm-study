# 检查是否存在负权环的方法为：记录一个点的入队次数，如果超过V-1次说明存在负权环，
# 因为最短路径上除自身外至多V-1个点，故一个点不可能被更新超过V-1次
from collections import defaultdict, deque


n, m = map(int, input().split())
adjMap = defaultdict(list)
for _ in range(m):
    u, v, w = list(map(int, input().split()))
    adjMap[u].append((v, w))


def spfa(n: int, adjMap: defaultdict, start: int, target: int) -> int:
    # """spfa求单源最短路，适用于解决带有负权重的图，是Bellman-ford的常数优化版"""
    dist = [int(1e20)] * (n + 1)
    dist[start] = 0

    queue = deque([start])
    isInqueue = [False] * (n + 1)  # 在队列里的点
    isInqueue[start] = True

    while queue:
        cur = queue.popleft()
        isInqueue[cur] = False  # 点从队列出来了

        # 更新过谁，就拿谁去更新别人
        for next, weight in adjMap[cur]:
            if dist[cur] + weight < dist[next]:
                dist[next] = dist[cur] + weight
                if not isInqueue[next]:
                    isInqueue[next] = True
                    queue.append(next)

    return dist[target]


res = spfa(n, adjMap, 1, n)
print(res if res < int(1e10) else 'impossible')
