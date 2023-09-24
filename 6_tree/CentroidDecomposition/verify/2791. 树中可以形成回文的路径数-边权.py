# 2791. 树中可以形成回文的路径数-边权
# https://leetcode.cn/problems/count-paths-that-can-form-a-palindrome-in-a-tree/
# 给你一棵 树（即，一个连通、无向且无环的图），根 节点为 0 ，由编号从 0 到 n - 1 的 n 个节点组成。
# 这棵树用一个长度为 n 、下标从 0 开始的数组 parent 表示，其中 parent[i] 为节点 i 的父节点，由于节点 0 为根节点，所以 parent[0] == -1 。
# 另给你一个长度为 n 的字符串 s ，其中 s[i] 是分配给 i 和 parent[i] 之间的边的字符。s[0] 可以忽略。
# !找出满足 u < v ，且从 u 到 v 的路径上分配的字符可以 重新排列 形成 回文 的所有节点对 (u, v) ，并返回节点对的数目。
# 树中的回文路径数


# !TLE


from collections import defaultdict
from typing import DefaultDict, List, Tuple


class Solution:
    def countPalindromePaths(self, parent: List[int], s: str) -> int:
        n = len(parent)
        tree = [[] for _ in range(n)]
        for i in range(1, n):
            p = parent[i]
            cost = 1 << (ord(s[i]) - 97)
            tree[i].append((p, cost))
            tree[p].append((i, cost))

        centTree, root = centroidDecomposition(n, tree)
        removed = [False] * n
        res = 0

        def collect(cur: int, pre: int, sub: "DefaultDict[int, int]", state: int) -> None:
            """统计子树内的答案."""
            sub[state] += 1
            for next, cost in tree[cur]:
                if next != pre and not removed[next]:
                    collect(next, cur, sub, state ^ cost)

        def decomposition(cur: int, pre: int) -> None:
            nonlocal res
            removed[cur] = True
            for next in centTree[cur]:  # 点分树的子树中的答案(不经过重心)
                if not removed[next]:
                    decomposition(next, cur)
            removed[cur] = False

            counter = defaultdict(int, {0: 1})  # 经过重心的路径
            for next, cost in tree[cur]:
                if next == pre or removed[next]:
                    continue
                sub = defaultdict(int)  # 统计子树内0和1的路径数(不包含当前根节点cur)
                collect(next, cur, sub, cost)
                for state, count in sub.items():
                    res += count * counter[state]
                    for i in range(26):
                        res += count * counter[state ^ (1 << i)]
                for state, count in sub.items():
                    counter[state] += count

        decomposition(root, -1)
        return res


def centroidDecomposition(n: int, tree: List[List[Tuple[int, int]]]) -> Tuple[List[List[int]], int]:
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
        for next, _ in tree[cur]:
            if next == parent or removed[next]:
                continue
            subSize[cur] += getSize(next, cur)
        return subSize[cur]

    def getCentroid(cur: int, parent: int, mid: int) -> int:
        for next, _ in tree[cur]:
            if next == parent or removed[next]:
                continue
            if subSize[next] > mid:
                return getCentroid(next, cur, mid)
        return cur

    def build(cur: int) -> int:
        centroid = getCentroid(cur, -1, getSize(cur, -1) // 2)
        removed[centroid] = True
        for next, _ in tree[centroid]:
            if not removed[next]:
                centTree[centroid].append(build(next))
        removed[centroid] = False
        return centroid

    subSize = [0] * n
    removed = [False] * n
    centTree = [[] for _ in range(n)]
    root = build(0)
    return centTree, root
