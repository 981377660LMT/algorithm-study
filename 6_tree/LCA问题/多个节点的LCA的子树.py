from typing import List, Optional


class Tree:
    def __init__(self, val, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def solve(self, root: Tree, values: List[int]):
        """子树结点互不相同"""

        def dfs(root) -> Optional[Tree]:
            """这个root是否被需要"""
            if not root:
                return None

            if root.val in need:
                return root

            left, right = dfs(root.left), dfs(root.right)
            if left and right:
                return root
            else:
                return left or right

        need = set(values)
        return dfs(root)

