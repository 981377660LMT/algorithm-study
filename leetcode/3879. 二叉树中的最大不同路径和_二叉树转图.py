# O(n^2) 枚举路径端点

from collections import defaultdict
from typing import Optional, Set

INF = int(1e18)


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def binaryTreeToGraph(root: Optional[TreeNode]):
    """二叉树转图.需要注意只有一个结点的情形."""
    graph = defaultdict(list)

    def dfs(node: Optional[TreeNode], parent: Optional[TreeNode]):
        if not node:
            return
        if parent:
            graph[node].append(parent)
            graph[parent].append(node)
        dfs(node.left, node)
        dfs(node.right, node)

    dfs(root, None)
    return graph


class Solution:
    def maxSum(self, root: Optional[TreeNode]) -> int:
        if not root:
            return 0
        if not root.left and not root.right:
            return root.val

        graph = binaryTreeToGraph(root)
        res = -INF

        def dfs(cur: TreeNode, curSum: int, visited: Set[int]):
            nonlocal res
            res = max(res, curSum)
            for next_ in graph[cur]:
                if next_.val not in visited:
                    visited.add(next_.val)
                    dfs(next_, curSum + next_.val, visited)
                    visited.remove(next_.val)

        for node in graph:
            dfs(node, node.val, {node.val})

        return res
