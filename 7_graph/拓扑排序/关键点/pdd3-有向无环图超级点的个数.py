# pdd的T3
# !给一个有向无环图，超级点定义为图中任意点都可以从自己到超级点或者超级点到自己；问有多少个超级点.
# 在一个有向无环图中，如果一个点u是超级点，那么它必须满足对于图中的任意其他点v，
# 要么存在一条从u到v的路径，要么存在一条从v到u的路径，求图中超级点的个数。
# 数据范围2<=节点数<=3e5，1<=边数<=3e5
#
#
# 如果不是DAG，缩点即可
# 我们可以遍历所有层，并且维护一个集合 S，表示目前不可达的顶点。
# 初始的时候为空。之后按照编号从小到大遍历每一层，考虑层中的所有顶点，
# 删除 S 中的可以被这一层的顶点抵达的那些顶点(反图上)，之后把整层中的所有顶点全部加入到 S 中。
# 如果某一层扫描完成，且 S 中只有一个元素，那么这一层中的唯一顶点就是关键点。

from collections import deque
from typing import List, Tuple


def superPointOnDag(n: int, edges: List[Tuple[int, int]]) -> List[bool]:
    adjList = [[] for _ in range(n)]
    rAdjList = [[] for _ in range(n)]
    deg = [0] * n
    for u, v in edges:
        adjList[u].append(v)
        deg[v] += 1
        rAdjList[v].append(u)

    queue = deque(i for i in range(n) if deg[i] == 0)

    unReachable = set()
    res = [False] * n
    while queue:
        len_ = len(queue)
        level = queue.copy()
        for _ in range(len_):
            cur = queue.popleft()
            for next_ in adjList[cur]:
                deg[next_] -= 1
                if deg[next_] == 0:
                    queue.append(next_)
            for prev in rAdjList[cur]:
                unReachable.discard(prev)
        unReachable.update(level)  # 之后把整层中的所有顶点全部加入到 S 中
        if len(unReachable) == 1:
            res[level[0]] = True

    return res


if __name__ == "__main__":
    n = 6
    edges = [(0, 1), (1, 2), (1, 3), (2, 4), (3, 4), (4, 5)]
    print(superPointOnDag(n, edges))
    n = 4
    edges = [(0, 1), (1, 2), (2, 3)]
    print(superPointOnDag(4, edges))
    n = 4
    edges = [(0, 1), (0, 2), (0, 3)]
    print(superPointOnDag(4, edges))
