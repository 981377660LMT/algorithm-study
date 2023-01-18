# 二叉树转无向图 / n叉树转无向图
# 采用先序遍历给每个结点分配id

from collections import deque
from typing import List, Optional, Tuple


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def treeToGraph1(
    root: Optional["TreeNode"],
) -> Tuple[List[List[int]], List[Tuple[int, int]], List[int]]:
    """二叉树转无向图,返回邻接表,边列表,结点值列表"""

    def dfs(root: Optional["TreeNode"]) -> int:
        if not root:
            return -1
        nonlocal globalId
        values.append(root.val)
        curId = globalId
        globalId += 1
        if root.left:
            nextId = dfs(root.left)
            edges.append((curId, nextId))
        if root.right:
            nextId = dfs(root.right)
            edges.append((curId, nextId))
        return curId

    globalId = 0  # 0 - n-1
    edges = []
    values = []
    dfs(root)

    n = len(values)
    adjList = [[] for _ in range(n)]
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)
    return adjList, edges, values


if __name__ == "__main__":
    # 863. 二叉树中所有距离为 K 的结点
    class Solution:
        def distanceK(self, root: TreeNode, target: TreeNode, k: int) -> List[int]:
            adjList, _, values = treeToGraph1(root)
            n = len(values)
            targetId = values.index(target.val)

            queue = deque([targetId])
            visited = [False] * n
            visited[targetId] = True
            for _ in range(k):
                len_ = len(queue)
                for _ in range(len_):
                    cur = queue.popleft()
                    for next in adjList[cur]:
                        if not visited[next]:
                            queue.append(next)
                            visited[next] = True

            return [values[i] for i in queue]


class Node:
    def __init__(self, val=None, children=None):
        self.val = val
        self.children = children if children is not None else []


def treeToGraph2(
    root: Optional["Node"],
) -> Tuple[List[List[int]], List[Tuple[int, int]], List[int]]:
    """n叉树转无向图,返回邻接表,边列表,结点值列表"""

    def build(root: Optional["Node"]) -> int:
        if not root:
            return -1
        nonlocal globalId
        values.append(root.val)
        curId = globalId
        globalId += 1
        for next in root.children:
            nextId = build(next)
            edges.append((curId, nextId))
        return curId

    globalId = 0  # 0 - n-1
    edges = []
    values = []
    build(root)

    n = len(values)
    adjList = [[] for _ in range(n)]
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)

    return adjList, edges, values
