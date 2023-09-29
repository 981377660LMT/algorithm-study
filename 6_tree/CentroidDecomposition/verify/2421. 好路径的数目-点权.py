# https://leetcode.cn/problems/number-of-good-paths/
# 一条 好路径需要满足以下条件：

# !开始节点和结束节点的值 相同。
# !开始节点和结束节点中间的所有节点值都小于等于开始节点的值（也就是说开始节点的值应该是路径上所有节点的最大值）。
# 请你返回不同好路径的数目。

from collections import defaultdict
from typing import DefaultDict, List, Tuple


class Solution:
    def numberOfGoodPaths(self, vals: List[int], edges: List[List[int]]) -> int:
        def collect(cur: int, pre: int, sub: "DefaultDict[int, int]", maxValue: int) -> None:
            """统计子树内的答案."""
            curValue = vals[cur]
            if curValue >= maxValue:
                sub[curValue] += 1
                maxValue = curValue
            for next in tree[cur]:
                if next != pre and not removed[next]:
                    collect(next, cur, sub, maxValue)

        def decomposition(cur: int, pre: int) -> None:
            nonlocal res
            removed[cur] = True
            for next in centTree[cur]:  # 点分树的子树中的答案(不经过重心)
                if not removed[next]:
                    decomposition(next, cur)
            removed[cur] = False

            counter = defaultdict(int, {vals[cur]: 1})  # 经过重心的路径
            for next in tree[cur]:
                if next == pre or removed[next]:
                    continue
                sub = defaultdict(int)
                collect(next, cur, sub, vals[cur])
                for v, count in sub.items():
                    res += count * counter[v]
                    counter[v] += count

        n = len(edges) + 1
        tree = [[] for _ in range(n)]
        for u, v in edges:
            tree[u].append(v)
            tree[v].append(u)

        centTree, root = centroidDecomposition(n, tree)
        removed = [False] * n
        res = 0

        decomposition(root, -1)
        return res + n


def centroidDecomposition(n: int, tree: List[List[int]]) -> Tuple[List[List[int]], int]:
    """
    树的重心分解, 返回点分树和点分树的根.

    Params:
        n: 树的节点数.
        tree: `无向树`的邻接表.
    Returns:
        centTree: 重心互相连接形成的有根树, 可以想象把树拎起来, 重心在树的中心，连接着各个子树的重心...
        root: 点分树的根.
    """

    def getSize(cur: int, parent: int) -> int:
        subSize[cur] = 1
        for next in tree[cur]:
            if next == parent or removed[next]:
                continue
            subSize[cur] += getSize(next, cur)
        return subSize[cur]

    def getCentroid(cur: int, parent: int, mid: int) -> int:
        for next in tree[cur]:
            if next == parent or removed[next]:
                continue
            if subSize[next] > mid:
                return getCentroid(next, cur, mid)
        return cur

    def build(cur: int) -> int:
        centroid = getCentroid(cur, -1, getSize(cur, -1) // 2)
        removed[centroid] = True
        for next in tree[centroid]:
            if not removed[next]:
                centTree[centroid].append(build(next))
        removed[centroid] = False
        return centroid

    subSize = [0] * n
    removed = [False] * n
    centTree = [[] for _ in range(n)]
    root = build(0)
    return centTree, root
