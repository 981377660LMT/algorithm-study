from collections import defaultdict
from functools import lru_cache
from typing import Optional


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def treeToGraph(root: TreeNode):
    """二叉树转图 结点id作为键"""

    def dfs(root: Optional[TreeNode], parent: Optional[TreeNode]) -> None:
        if not root:
            return
        rootId = id(root)
        valueById[rootId] = root.val
        if parent is not None:
            parentId = id(parent)
            adjMap[parentId].add(rootId)
            adjMap[rootId].add(parentId)
        dfs(root.left, root)
        dfs(root.right, root)

    adjMap, valueById = defaultdict(set), defaultdict(int)
    dfs(root, None)
    return adjMap, valueById


class Solution:
    def maxPathSum(self, root: TreeNode) -> int:
        @lru_cache(None)
        def dfs(cur: int, parent: int) -> int:
            res = valueById[cur]
            for next in adjMap[cur]:
                if next == parent:
                    continue
                res = max(res, dfs(next, cur) + valueById[cur])
            return res

        adjMap, valueById = treeToGraph(root)
        return max(((dfs(start, -1)) for start in adjMap.keys()), default=root.val)

