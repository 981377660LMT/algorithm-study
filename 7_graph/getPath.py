# !给定一张有向图，图中可能有环
# !再给定起点和终点，找任意一条连接两点之间的路径


from collections import deque
from typing import List


def findPathDfs(graph: List[List[int]], start: int, end: int) -> List[int]:
    if start == end:
        return [start]

    n = len(graph)
    visited = [False] * n
    path = []

    def dfs(cur: int) -> bool:
        if cur == end:
            path.append(cur)
            return True
        visited[cur] = True
        path.append(cur)
        for next_ in graph[cur]:
            if not visited[next_] and dfs(next_):
                return True
        path.pop()
        return False

    dfs(start)
    return path


def findPathBfs(graph: List[List[int]], start: int, end: int) -> List[int]:
    if start == end:
        return [start]
    n = len(graph)
    queue = deque([start])
    visited = [False] * n
    visited[start] = True
    pre = [-1] * n
    while queue:
        cur = queue.popleft()
        for next_ in graph[cur]:
            if visited[next_]:
                continue
            visited[next_] = True
            queue.append(next_)
            pre[next_] = cur
            if next_ == end:
                path = []
                while next_ != -1:
                    path.append(next_)
                    next_ = pre[next_]
                return path[::-1]
    return []


if __name__ == "__main__":
    # 添加边，构建有环图
    n = 4
    graph = [[] for _ in range(n)]
    edges = [(0, 1), (0, 2), (1, 2), (2, 0), (2, 3), (3, 3)]

    for u, v in edges:
        graph[u].append(v)

    start, end = 0, 3
    print(findPathDfs(graph, start, end))  # [0, 1, 2, 3]
    print(findPathBfs(graph, start, end))  # [0, 2, 3]
