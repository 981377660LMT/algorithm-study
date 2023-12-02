# Tree Isomorphism 树同构.
# https://oi-wiki.org/graph/tree-ahu/
#
# Api:
#  TreeHasher: 可能出错，但是快.
#    hashAll(graph) int
#    hasSubTrees(graph, root) List[int]
#  TreeEncoder: AHU算法，不会出错，但是慢.
#    encode(graph, root) List[int]


from collections import defaultdict
import random
from typing import List, Tuple


MOD61 = (1 << 61) - 1


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


def add(a: int, b: int) -> int:
    a += b
    if a >= MOD61:
        a -= MOD61
    return a


def mul(a: int, b: int) -> int:
    c = a * b
    return add(c >> 61, c & MOD61)


def getTreeDiameter(n: int, tree: List[List[int]], start=0) -> Tuple[int, List[int]]:
    """求无权树的(直径长度,直径路径)."""

    def dfs(start: int) -> Tuple[int, List[int]]:
        dist = [-1] * n
        dist[start] = 0
        stack = [start]
        while stack:
            cur = stack.pop()
            for next in tree[cur]:
                if dist[next] != -1:
                    continue
                dist[next] = dist[cur] + 1
                stack.append(next)
        endPoint = dist.index(max(dist))
        return endPoint, dist

    u, _ = dfs(start)
    v, dist = dfs(u)
    diameter = dist[v]
    path = [v]
    while u != v:
        for next in tree[v]:
            if dist[next] + 1 == dist[v]:
                path.append(next)
                v = next
                break

    return diameter, path


class TreeHasher:
    __slots__ = ("_rand", "_r")

    def __init__(self):
        self._rand = lambda: random.randint(1, MOD61 - 1)
        self._r = []

    def hashAll(self, tree: List[List[int]], root=-1) -> int:
        res = 0
        if root == -1:
            d, path = getTreeDiameter(len(tree), tree)
            res = self._dfsAll(tree, path[d // 2], -1)[0]
            if d % 2 == 1:
                res = min2(res, self._dfsAll(tree, path[d // 2 + 1], -1)[0])
        else:
            res = self._dfsAll(tree, root, -1)[0]
        return res

    def hashSubTrees(self, tree: List[List[int]], root: int) -> List[int]:
        hash = [0] * len(tree)
        self._dfsSubTrees(tree, hash, root, -1)
        return hash

    def _dfsAll(self, tree: List[List[int]], v: int, p: int) -> Tuple[int, int]:
        maxd = 0
        hash = []
        for c in tree[v]:
            if c != p:
                h, d = self._dfsAll(tree, c, v)
                maxd = max2(maxd, d + 1)
                hash.append(h)
        if len(self._r) == maxd:
            self._r.append(self._rand())
        res = 1
        for h in hash:
            res = mul(res, add(self._r[maxd], h))
        return res, maxd

    def _dfsSubTrees(self, tree: List[List[int]], hash: List[int], v: int, p: int) -> int:
        maxd = 0
        for c in tree[v]:
            if c != p:
                maxd = max2(maxd, self._dfsSubTrees(tree, hash, c, v) + 1)
        if len(self._r) == maxd:
            self._r.append(self._rand())
        res = 1
        for c in tree[v]:
            if c != p:
                res = mul(res, add(self._r[maxd], hash[c]))
        hash[v] = res
        return maxd


class TreeEncoder:
    __slots__ = "_mp"

    def __init__(self):
        self._mp = dict({tuple([]): 0})

    def encode(self, tree: List[List[int]], root: int) -> List[int]:
        labels = [0] * len(tree)
        self._dfs(tree, labels, root, -1)
        return labels

    def _dfs(self, tree: List[List[int]], val: List[int], v: int, p: int) -> None:
        children = []
        for c in tree[v]:
            if c != p:
                self._dfs(tree, val, c, v)
                children.append(val[c])
        children.sort()
        key = tuple(children)
        val[v] = self._mp.setdefault(key, len(self._mp))  # type: ignore


if __name__ == "__main__":
    import sys

    input = sys.stdin.readline
    sys.setrecursionlimit(int(1e6))

    # https://judge.yosupo.jp/problem/rooted_tree_isomorphism_classification
    def rootedTreeIsomorphismClassification1(tree: List[List[int]]) -> Tuple[int, List[int]]:
        encoder = TreeEncoder()
        labels = encoder.encode(tree, 0)
        return max(labels) + 1, labels

    # https://judge.yosupo.jp/problem/rooted_tree_isomorphism_classification
    def rootedTreeIsomorphismClassification2(tree: List[List[int]]) -> Tuple[int, List[int]]:
        encoder = TreeHasher()
        labels = encoder.hashSubTrees(tree, 0)
        mp = defaultdict(lambda: len(mp))
        for i, v in enumerate(labels):
            labels[i] = mp[v]
        return len(mp), labels

    n = int(input())
    parents = list(map(int, input().split()))
    adjList = [[] for _ in range(n)]
    for i, p in enumerate(parents):
        adjList[p].append(i + 1)
        adjList[i + 1].append(p)

    count, labels = rootedTreeIsomorphismClassification2(adjList)
    print(count)
    print(*labels)
