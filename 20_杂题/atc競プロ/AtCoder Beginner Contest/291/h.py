import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
from collections import defaultdict, deque
from typing import Iterable, List, Mapping, Sequence, Union

AdjList = Sequence[Iterable[int]]
AdjMap = Mapping[int, Iterable[int]]
Tree = Union[AdjList, AdjMap]


def getCenter(n: int, tree: "Tree", root=0) -> List[int]:
    """求重心"""

    def dfs(cur: int, pre: int) -> None:
        subsize[cur] = 1
        for next in tree[cur]:
            if next == pre:
                continue
            dfs(next, cur)
            subsize[cur] += subsize[next]
            weight[cur] = max(weight[cur], subsize[next])
        weight[cur] = max(weight[cur], n - subsize[cur])
        if weight[cur] <= n / 2:
            res.append(cur)

    res = []
    weight = [0] * n  # 节点的`重量`，即以该节点为根的子树的最大节点数
    subsize = [0] * n  # 子树大小
    dfs(root, -1)
    return res


if __name__ == "__main__":
    n = int(input())
    edges = [tuple(map(int, input().split())) for _ in range(n - 1)]
    adjList = [[] for _ in range(n)]
    for u, v in edges:
        u -= 1
        v -= 1
        adjList[u].append(v)
        adjList[v].append(u)

    center = getCenter(n, adjList)

    newRoot = center[0]
    if len(center) == 2:
        newRoot = center[1]
    parents = [-1] * n

    queue = deque([(newRoot, -1)])
    while queue:
        cur, pre = queue.popleft()
        for next in adjList[cur]:
            if next == pre:
                continue
            queue.append((next, cur))
            parents[next] = cur

    for p in parents:
        print(-1 if p == -1 else p + 1, end=" ")
