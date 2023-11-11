# !子树中特定深度的结点个数
# 给你一棵树，q 次询问；
# 每次询问输入两个数，u，d；
# 要求计算出在树的第 d 层并且是节点 u 的子孙节点的节点个数

# !1. depth 邻接表存递增时间戳+二分
# !2. 离线查询 + 树上启发式合并


from bisect import bisect_left, bisect_right
from collections import defaultdict
import sys
from typing import List, Tuple


def countDescendants1(n: int, parents: List[int], queries: List[Tuple[int, int]]) -> List[int]:
    """depth 邻接表存递增时间戳 + 二分"""

    def dfs(cur: int, pre: int, dep: int) -> None:
        nonlocal dfsId
        ins[cur] = dfsId
        for next in adjList[cur]:
            if next != pre:
                dfs(next, cur, dep + 1)
        outs[cur] = dfsId
        depth[dep].append(dfsId)
        dfsId += 1

    adjList = [[] for _ in range(n)]
    for i, p in enumerate(parents):
        if p == -1:
            continue
        adjList[p].append(i)
        adjList[i].append(p)

    ins, outs, dfsId = [0] * (n + 10), [0] * (n + 10), 1
    depth = defaultdict(list)
    dfs(0, -1, 0)

    res = []
    for i in range(len(queries)):
        queryRoot, queryDepth = queries[i]
        left = ins[queryRoot]
        right = outs[queryRoot]
        ids = depth[queryDepth]
        res.append(bisect_right(ids, right) - bisect_left(ids, left))
    return res


def countDescendants2(n: int, parents: List[int], queries: List[Tuple[int, int]]) -> List[int]:
    """离线查询 + 树上启发式合并"""

    def dfs(cur: int, pre: int, dep: int):
        """深度->结点个数"""
        subTree = {dep: 1}
        for next in adjList[cur]:
            if next != pre:
                nextTree = dfs(next, cur, dep + 1)
                if len(nextTree) > len(subTree):
                    subTree, nextTree = nextTree, subTree
                for k, v in nextTree.items():
                    subTree[k] = subTree.get(k, 0) + v
        for qi, qd in queryGroup[cur]:
            res[qi] = subTree.get(qd, 0)
        return subTree

    adjList = [[] for _ in range(n)]
    for i, p in enumerate(parents):
        if p == -1:
            continue
        adjList[p].append(i)
        adjList[i].append(p)

    q = len(queries)
    res = [0] * q
    queryGroup = [[] for _ in range(n)]
    for qi, (qRoot, qDepth) in enumerate(queries):
        queryGroup[qRoot].append((qi, qDepth))
    dfs(0, -1, 0)
    return res


if __name__ == "__main__":
    sys.setrecursionlimit(int(1e6))
    n = int(input())
    parents = [-1] + [int(x) - 1 for x in input().split()]
    q = int(input())
    queries = []
    for _ in range(q):
        u, d = map(int, input().split())
        queries.append((u - 1, d))
    # print(*countDescendants1(n, parents, queries), sep="\n")
    print(*countDescendants2(n, parents, queries), sep="\n")
