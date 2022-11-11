from collections import deque
from typing import List

# n, m ≤ 250
# 从目标出发,bfs,回到自己就是最短环
# !bfs求无向无权图最小环长度

INF = int(1e9)


# !无权图中找到包含目标结点的最小环的长度
def minCycle1(n: int, graph: List[List[int]], start: int) -> int:
    queue = deque([(start, 0)])
    dist = [INF] * n
    dist[start] = 0
    while queue:
        len_ = len(queue)
        for _ in range(len_):
            cur, cost = queue.popleft()
            for next in graph[cur]:
                if next == start:
                    return cost + 1
                cand = cost + 1
                if cand < dist[next]:
                    dist[next] = cand
                    queue.append((next, cand))
    return INF


# !无权图求最小环的长度
# https://github.dev/EndlessCheng/codeforces-go/blob/cca30623b9ac0f3333348ca61b4894cd00b753cc/copypasta/graph.go#L470
def minCycle2(n: int, graph: List[List[int]]) -> int:
    res = INF
    dist = [-1] * n
    for start in range(n):
        visited = [False] * n
        visited[start] = True
        dist[start] = 0
        queue = deque([(start, -1)])
        while True:
            findCycle = False
            while queue:
                if findCycle:
                    break
                cur, pre = queue.popleft()
                for next in graph[cur]:
                    if dist[next] == -1:
                        dist[next] = dist[cur] + 1
                        queue.append((next, cur))
                        visited[next] = True
                    elif next != pre:
                        cand = dist[next] + dist[cur] + 1
                        if cand < res:
                            res = cand
                        findCycle = True
                        break
            if findCycle:
                break

        for i in range(n):
            if visited[i]:
                dist[i] = -1

    return res


if __name__ == "__main__":
    assert minCycle1(n=3, graph=[[1], [2], [0]], start=0) == 3
    assert minCycle2(3, graph=[[1], [2], [0]]) == 3
    assert minCycle2(4, graph=[[1], [2], [3], [0]]) == 4

    edges = [(0, 1), (1, 2), (0, 5), (1, 4), (2, 3), (3, 4), (4, 5)]
    graph = [[] for _ in range(6)]
    for u, v in edges:
        graph[u].append(v)
        graph[v].append(u)
    assert minCycle2(6, graph) == 4
