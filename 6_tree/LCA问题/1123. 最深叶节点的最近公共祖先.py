from typing import Optional, Tuple


class TreeNode:
    def __init__(
        self, val: int = 0, left: Optional["TreeNode"] = None, right: Optional["TreeNode"] = None
    ):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def lcaDeepestLeaves(self, root: Optional[TreeNode]) -> Optional[TreeNode]:
        def dfs(cur: Optional[TreeNode]) -> Tuple[int, Optional[TreeNode]]:
            if not cur:
                return 0, None
            leftHeight, leftLca = dfs(cur.left)
            rightHeight, rightLca = dfs(cur.right)
            if leftHeight > rightHeight:
                return leftHeight + 1, leftLca
            if leftHeight < rightHeight:
                return rightHeight + 1, rightLca
            return leftHeight + 1, cur

        _, res = dfs(root)
        return res
