#  https://leetcode.cn/problems/count-valid-paths-in-a-tree/description/
#  100047. 统计树中的合法路径数目
#  给你一棵 n 个节点的无向树，节点编号为 1 到 n 。
#  给你一个整数 n 和一个长度为 n - 1 的二维整数数组 edges ，
#  其中 edges[i] = [ui, vi] 表示节点 ui 和 vi 在树中有一条边。
#  请你返回树中的 合法路径数目 。
#  !如果在节点 a 到节点 b 之间 恰好有一个 节点的编号是质数，那么我们称路径 (a, b) 是 合法的 。
#  注意：
#  路径 (a, b) 指的是一条从节点 a 开始到节点 b 结束的一个节点序列，序列中的节点 互不相同 ，且相邻节点之间在树上有一条边。
#  路径 (a, b) 和路径 (b, a) 视为 同一条 路径，且只计入答案 一次 。

from collections import defaultdict
from typing import DefaultDict, List, Tuple
from collections import Counter


class EratosthenesSieve:
    """埃氏筛"""

    __slots__ = "_minPrime"  # 每个数的最小质因数

    def __init__(self, maxN: int):
        """预处理 O(nloglogn)"""
        minPrime = list(range(maxN + 1))
        upper = int(maxN**0.5) + 1
        for i in range(2, upper):
            if minPrime[i] < i:
                continue
            for j in range(i * i, maxN + 1, i):
                if minPrime[j] == j:
                    minPrime[j] = i
        self._minPrime = minPrime

    def isPrime(self, n: int) -> bool:
        if n < 2:
            return False
        return self._minPrime[n] == n

    def getPrimeFactors(self, n: int) -> "Counter[int]":
        """求n的质因数分解 O(logn)"""
        res, f = Counter(), self._minPrime
        while n > 1:
            m = f[n]
            res[m] += 1
            n //= m
        return res

    def getPrimes(self) -> List[int]:
        return [x for i, x in enumerate(self._minPrime) if i >= 2 and i == x]


E = EratosthenesSieve(int(1e5) + 10)


class Solution:
    def countPaths(self, n: int, edges: List[List[int]]) -> int:
        def collect(cur: int, pre: int, sub: "DefaultDict[int, int]", primeCount: int) -> None:
            """统计子树内的答案."""
            primeCount += int(E.isPrime(cur + 1))
            if primeCount <= 1:
                sub[primeCount] += 1
            else:
                return
            for next in tree[cur]:
                if next != pre and not removed[next]:
                    collect(next, cur, sub, primeCount)

        def decomposition(cur: int, pre: int) -> None:
            nonlocal res
            removed[cur] = True
            for next in centTree[cur]:  # 点分树的子树中的答案(不经过重心)
                if not removed[next]:
                    decomposition(next, cur)
            removed[cur] = False

            init = int(E.isPrime(cur + 1))
            counter = defaultdict(int, {init: 1})  # 经过重心的路径
            for next in tree[cur]:
                if next == pre or removed[next]:
                    continue
                sub = defaultdict(int)  # 统计子树内0和1的路径数(不包含当前根节点cur)
                collect(next, cur, sub, init)
                if init == 0:
                    res += sub[0] * counter[1]
                    res += sub[1] * counter[0]
                    counter[0] += sub[0]
                    counter[1] += sub[1]
                else:
                    res += sub[1] * counter[1]
                    counter[1] += sub[1]

        tree = [[] for _ in range(n)]
        for u, v in edges:
            u, v = u - 1, v - 1
            tree[u].append(v)
            tree[v].append(u)

        centTree, root = centroidDecomposition(n, tree)
        removed = [False] * n
        res = 0

        decomposition(root, -1)
        return res


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
