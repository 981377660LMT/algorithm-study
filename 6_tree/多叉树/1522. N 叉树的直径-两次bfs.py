from collections import defaultdict, deque
from typing import DefaultDict, Optional, Set


class Node:
    def __init__(self, val=None, children=None):
        self.val = val
        self.children = children if children is not None else []


def treeToGraph(root: 'Node') -> DefaultDict[int, Set[int]]:
    """n叉树转图,节点值不一定互异时使用id来作为结点唯一标识"""

    def dfs(root: Optional['Node'], parent: Optional['Node']) -> None:
        if not root:
            return
        rootId, parentId = id(root), id(parent)
        if parent is not None:
            adjMap[parentId].add(rootId)
            adjMap[rootId].add(parentId)
        for child in root.children:
            dfs(child, root)

    adjMap = defaultdict(set)
    dfs(root, None)
    return adjMap


def getDiameter(adjMap: DefaultDict[int, Set[int]], start: int) -> int:
    queue = deque([start])
    visited = set([start])
    lastVisited = start  # 全局变量，好记录第一次BFS最后一个点的ID
    while queue:
        curLen = len(queue)
        for _ in range(curLen):
            lastVisited = queue.popleft()
            for next in adjMap[lastVisited]:
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
            for next in adjMap[cur]:
                if next not in visited:
                    visited.add(next)
                    queue.append(next)
        level += 1

    return level


class Solution:
    def diameter(self, root: 'Node') -> int:
        adjMap = treeToGraph(root)
        return getDiameter(adjMap, id(root))


print(Solution().diameter(Node(1, [Node(2), Node(3)])))
