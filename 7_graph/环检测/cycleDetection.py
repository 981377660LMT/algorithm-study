from typing import List, Tuple


def cycleDetectionUndirected(n: int, edges: List[Tuple[int, int]]) -> Tuple[List[int], List[int]]:
    """无向图找环(有多个环时，返回任意一个环)

    Args:
        n (int): 节点数
        edges (List[Tuple[int, int]]): 无向边

    Returns:
        Tuple[List[int], List[int]] : 环上边的索引，环上点的索引
    """

    def dfs() -> Tuple[int, int, int, List[int], List[int]]:
        visited = [0] * n
        used = [0] * m
        preV = [-1] * n
        preE = [-1] * n
        for i in range(n):
            if visited[i]:
                continue
            stack = [(i, -1, -1)]  # (cur, preV, preE)
            while stack:
                v, p, e = stack.pop()
                if e != -1 and used[e]:
                    continue
                if visited[v]:
                    preV[v] = p
                    if e != -1:
                        preE[v] = e
                    return v, p, e, preV, preE
                visited[v] = 1
                if e != -1:
                    preE[v] = e
                    used[e] = 1
                preV[v] = p
                for a, e in adjList[v]:
                    if used[e]:
                        continue
                    stack.append((a, v, e))
        return -1, -1, -1, preV, preE

    def getCycle() -> Tuple[List[int], List[int]]:
        v, p, e, preV, preE = dfs()
        if v == -1:
            return [], []
        cycleEdges = [e]
        cycleNodes = [p]
        while v != p:
            e = preE[p]
            p = preV[p]
            cycleEdges.append(e)
            cycleNodes.append(p)
        return cycleEdges[::-1], cycleNodes[::-1]

    m = len(edges)
    adjList = [[] for _ in range(n)]
    for i, (u, v) in enumerate(edges):
        adjList[u].append((v, i))
        adjList[v].append((u, i))

    cycleEdges, cycleNodes = getCycle()
    return cycleEdges, cycleNodes


def cycleDetectionDirected(n: int, edges: List[Tuple[int, int]]) -> List[int]:
    """有向图找环(有多个环时，返回任意一个环)

    Args:
        n (int): 节点数
        edges (List[Tuple[int, int]]): 有向边

    Returns:
        List[int]: 环上边的索引
    """
    adjList = [[] for _ in range(n)]
    for i, (next, v) in enumerate(edges):
        adjList[next].append((v, i))

    visited = [0] * n
    finished = [0] * n
    stack = []
    for i in range(n):
        if visited[i]:
            continue
        todo = [(1, i, -1), (0, i, -1)]
        visited[i] = True
        while todo:
            kind, cur, edgeId = todo.pop()
            if kind == 0:
                if finished[cur]:
                    continue
                visited[cur] = 1
                stack.append((cur, edgeId))
                for next, id in adjList[cur]:
                    if finished[cur]:
                        continue
                    if visited[next] and finished[next] == 0:
                        cycle = [id]
                        while stack:
                            cur, id = stack.pop()
                            if cur == next:
                                break
                            cycle.append(id)
                        return cycle[::-1]
                    elif visited[next] == 0:
                        todo.append((1, next, id))
                        todo.append((0, next, id))
            else:
                if finished[cur]:
                    continue
                stack.pop()
                finished[cur] = 1
    return []


if __name__ == "__main__":
    n, m = map(int, input().split())
    edges = [tuple(map(int, input().split())) for _ in range(m)]
    res = cycleDetectionDirected(n, edges)
    if len(res) == 0:
        print(-1)
        exit(0)
    print(len(res))
    print(*res, sep="\n")
