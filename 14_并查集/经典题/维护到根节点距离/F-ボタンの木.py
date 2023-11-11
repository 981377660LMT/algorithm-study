# 树上移动硬币/金币, 类似 979. 在二叉树中分配硬币.
# 给定一棵树 初始每个点权值为starts[i].
# !第i个结点的邻接点有k个，可以将第i个结点的权值减去k，每个邻接点的权值加1.
# 问最少几次操作可以使得每个结点的权值为targets[i].
# !数据保证有解.

# !每个节点都可以通过操作向父节点“借用”权值

from typing import List, Tuple
from UnionFindWithDist import UnionFindArrayWithDist1

INF = int(1e18)


def treeOfButton(
    n: int, edges: List[Tuple[int, int]], starts: List[int], targets: List[int]
) -> int:
    def dfs(cur: int, pre: int) -> int:
        """
        返回每个结点需要从父节点获得的硬币.
        大于0表示当前子树需要从父节点获得硬币, 小于0表示当前子树需要向父节点提供硬币.
        """
        need = targets[cur] - starts[cur]
        for next in tree[cur]:
            if next == pre:
                continue
            sub = dfs(next, cur)
            uf.union(cur, next, sub)  # cur需要向next提供sub个硬币
            need += sub
        return need

    tree = [[] for _ in range(n)]
    for u, v in edges:
        tree[u].append(v)
        tree[v].append(u)
    uf = UnionFindArrayWithDist1(n)
    dfs(0, -1)

    sum_, min_ = 0, INF
    for i in range(n):
        diff = uf.dist(i, 0)
        sum_ += diff
        min_ = min(min_, diff)
    return sum_ - min_ * n


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n = int(input())
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v))
    starts = list(map(int, input().split()))
    targets = list(map(int, input().split()))
    print(treeOfButton(n, edges, starts, targets))
