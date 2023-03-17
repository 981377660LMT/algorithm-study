# 字典序拓扑排序(字典序最小拓扑序)
# 求字典序最小的拓扑排序/字典序最大的拓扑排序
# 0 - n-1的排列 ai要在bi前出现 求字典序最小的排列 不存在则输出-1


from heapq import heapify, heappop, heappush
from typing import List, Tuple


def lexicographicalOrderTopoSort(n: int, prerequisites: List[Tuple[int, int]]) -> List[int]:
    """(0到n-1)字典序最小的拓扑排序"""
    adjList = [[] for _ in range(n)]
    deg = [0] * (n)
    for pre, cur in prerequisites:
        adjList[pre].append(cur)
        deg[cur] += 1

    res = []
    pq = [i for i in range(n) if deg[i] == 0]
    heapify(pq)
    while pq:
        cur = heappop(pq)
        res.append(cur)
        for next in adjList[cur]:
            deg[next] -= 1
            if deg[next] == 0:
                heappush(pq, next)

    return [] if len(res) != n else res


if __name__ == "__main__":
    n, m = map(int, input().split())
    prerequisites = [tuple(map(int, input().split())) for _ in range(m)]
    res = lexicographicalOrderTopoSort(n, prerequisites)
    print(*res if res else [-1])
