# 虫洞非常奇特，它可以看作是一条 单向 路径，
# 通过它可以使你回到过去的某个时刻（相对于你进入虫洞之前）。
# 现在农夫约翰希望能够从农场中的某片田地出发，经过一些路径和虫洞回到过去，并在他的出发时刻之前赶到他的出发地。


from collections import defaultdict, deque
from typing import Mapping

INF = int(1e18)


def spfa2(n: int, adjMap: Mapping[int, Mapping[int, int]]) -> bool:
    """spfa判断负环 存在负环返回True 否则返回False

    在原图的基础上新建一个虚拟源点,
    从该点向其他所有点连一条权值为0的有向边。
    那么原图有负环等价于新图有负环
    也等价于开始时将所有点加入队列
    """
    dist = [0] * n
    queue = deque(range(n))
    inQueue = [True] * n
    count = [0] * n

    while queue:
        cur = queue.popleft()
        inQueue[cur] = False

        for next in adjMap[cur]:
            weight = adjMap[cur][next]
            cand = dist[cur] + weight
            if cand < dist[next]:  # 如果要最长路这里需要改成 >
                dist[next] = cand
                count[next] = count[cur] + 1
                if count[next] >= n:
                    return True
                if not inQueue[next]:
                    inQueue[next] = True
                    if queue and dist[next] < dist[queue[0]]:  # !酸辣粉优化 如果要最长路这里需要改成 >
                        queue.appendleft(next)
                    else:
                        queue.append(next)

    return False


f = int(input())

# 使用酸辣粉优化后	4280 ms => 3570 ms
for _ in range(f):
    N, M, W = map(int, input().split())
    adjMap = defaultdict(lambda: defaultdict(lambda: INF))
    for _ in range(M):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        adjMap[u][v] = min(adjMap[u][v], w)
        adjMap[v][u] = min(adjMap[v][u], w)

    for _ in range(W):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        adjMap[u][v] = min(adjMap[u][v], -w)  # 回到过去

    res = spfa2(N, adjMap)
    print("YES" if res else "NO")
