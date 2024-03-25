# pdd的T3
# !给一个有向无环图，超级点定义为图中任意点都可以从自己到超级点或者超级点到自己；问有多少个超级点.
# 在一个有向无环图中，如果一个点u是超级点，那么它必须满足对于图中的任意其他点v，
# 要么存在一条从u到v的路径，要么存在一条从v到u的路径，求图中超级点的个数。
# 数据范围2<=节点数<=3e5，1<=边数<=3e5
#
# !等同于正图和反图上拓扑排序都只有一个选择
# 拓扑排序，如果队列中有多个入度为0的点，就优先弹出id最小的那个，然后正反都做一次，如果位置能对应上就是超级点.
# 两次拓扑排序，一次优先弹出id较小的，一次优先弹出id较大的
# 如果一个点两次拓扑序都一样，就是关键点
# 如果不是DAG，缩点即可


from heapq import heapify, heappop, heappush
from typing import List, Tuple


def topoSortByHeap(
    n: int, adjList: List[List[int]], directed=True, minFirst=True
) -> Tuple[List[int], bool]:
    """使用优先队列的拓扑排序."""
    if directed:
        deg = [0] * n
        for i in range(n):
            for j in adjList[i]:
                deg[j] += 1
    else:
        deg = [len(adj) for adj in adjList]

    startDeg = 0 if directed else 1
    pq = [v if minFirst else -v for v in range(n) if deg[v] == startDeg]
    heapify(pq)
    res = []

    if minFirst:
        while pq:
            cur = heappop(pq)
            res.append(cur)
            for next in adjList[cur]:
                deg[next] -= 1
                if deg[next] == startDeg:
                    heappush(pq, next)
    else:
        while pq:
            cur = -heappop(pq)
            res.append(cur)
            for next in adjList[cur]:
                deg[next] -= 1
                if deg[next] == startDeg:
                    heappush(pq, -next)

    if len(res) != n:
        return [], False
    return res, True


# TODO: verify
# 4个点，边为[(0,1),(0,2),(0,3)]，2这个点会被误判为关键点
def superPoints(n: int, edges: List[Tuple[int, int]]) -> List[bool]:
    # adjList = [[] for _ in range(n)]
    # rAdjList = [[] for _ in range(n)]
    # for u, v in edges:
    #     adjList[u].append(v)
    #     rAdjList[v].append(u)
    # order = topoSortByHeap(n, adjList, minFirst=True)[0]
    # rOrder = topoSortByHeap(n, rAdjList, minFirst=True)[0]
    # print(order, rOrder)
    # rOrder.reverse()
    # return [a == b for a, b in zip(order, rOrder)]
    ...


if __name__ == "__main__":
    n = 4
    edges = [(0, 1), (0, 2), (0, 3)]
    print(superPoints(n, edges))  # [True, True, True, True, True]
