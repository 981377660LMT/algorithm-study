# bellman-ford法求带负权边的单源最短路. O(VE).


from typing import List, Tuple

INF = int(1e18)


def bellmanFord(
    n: int, adjList: List[List[Tuple[int, int]]], start: int
) -> Tuple[List[int], List[int]]:
    """
    BellmanFord算法`O(VE)`求解带负权边的单源最短路,并求出每个点的前驱.

    Args:
        n (int): 节点数.
        adjList (List[List[Tuple[int, int]]]): 邻接表.
        start (int): 起点.

    Returns:
        Tuple[List[int], List[int]]: (起点到各点的最短距离,每个点的前驱).
        距离为INF表示不可达.
        距离为-INF表示经过负环到达.
    """
    dist = [INF] * n
    pre = [-1] * n
    dist[start] = 0
    loop = 0
    while True:
        loop += 1
        updated = False
        for from_ in range(n):
            if dist[from_] == INF:
                continue
            for to, cost in adjList[from_]:
                cand = dist[from_] + cost
                if cand < -INF:
                    cand = -INF
                if cand < dist[to]:
                    updated = True
                    pre[to] = from_
                    if loop >= n:
                        cand = -INF
                    dist[to] = cand
        if not updated:
            break
    return dist, pre


if __name__ == "__main__":
    # https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=GRL_1_B
    n, m, s = map(int, input().split())
    adjList = [[] for _ in range(n)]
    for _ in range(m):
        from_, to, cost = map(int, input().split())
        adjList[from_].append((to, cost))
    dist, pre = bellmanFord(n, adjList, s)
    if min(dist) == -INF:
        print("NEGATIVE CYCLE")
    else:
        for x in dist:
            if x == INF:
                print("INF")
            else:
                print(x)
