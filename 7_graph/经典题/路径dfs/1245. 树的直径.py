from collections import defaultdict, deque
from functools import lru_cache
from typing import DefaultDict, Iterable, List, Mapping, Optional, Sequence, Set, Union


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def treeToGraph(root: TreeNode) -> DefaultDict[int, Set[int]]:
    """二叉树转图，注意节点值需要不同"""

    def dfs(root: Optional[TreeNode], parent: int) -> None:
        if not root:
            return
        if parent != -1:
            adjMap[parent].add(root.val)
            adjMap[root.val].add(parent)
        dfs(root.left, root.val)
        dfs(root.right, root.val)

    adjMap = defaultdict(set)
    dfs(root, -1)
    return adjMap


UndirectedGraph = Union[List[Iterable[int]], Mapping[int, Iterable[int]]]


def getDiameter(undirectedGraph: UndirectedGraph, start: int) -> int:
    queue = deque([start])
    visited = set([start])
    lastVisited = 0  # 全局变量，好记录第一次BFS最后一个点的ID
    while queue:
        curLen = len(queue)
        for _ in range(curLen):
            lastVisited = queue.popleft()
            for next in undirectedGraph[lastVisited]:
                if next not in visited:
                    visited.add(next)
                    queue.append(next)

    queue = deque([lastVisited])  # 第一次最后一个点，作为第二次BFS的起点
    visited = set([lastVisited])
    level = -1  # 记好距离
    while queue:
        curLen = len(queue)
        for _ in range(curLen):
            cur = queue.popleft()
            for next in undirectedGraph[cur]:
                if next not in visited:
                    visited.add(next)
                    queue.append(next)
        level += 1

    return level


######################################################################################
class Solution:
    def treeDiameter(self, edges: List[List[int]]) -> int:
        @lru_cache(None)
        def dfs(cur: int, pre: int) -> int:
            res = 1
            for next in adjMap[cur]:
                if next == pre:
                    continue
                res = max(res, dfs(next, cur) + 1)
            return res

        adjMap = defaultdict(set)
        for u, v in edges:
            adjMap[u].add(v)
            adjMap[v].add(u)

        n = len(adjMap)
        res = 1
        for i in range(n):
            res = max(res, dfs(i, -1))
        dfs.cache_clear()
        return res - 1

