# Definition for a binary tree node.
import heapq
from typing import Generator, List


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def getAllElements(self, root1: TreeNode, root2: TreeNode) -> List[int]:
        def dfs(root: TreeNode) -> Generator[int, None, None]:
            if not root:
                return
            yield from dfs(root.left)
            yield root.val
            yield from dfs(root.right)

        return list(heapq.merge(dfs(root1), dfs(root2)))

