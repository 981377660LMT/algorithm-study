from typing import Optional


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


# 返回该二叉树中值为 p 的结点与值为 q 的结点间的 距离 。
class Solution:
    def findDistance(self, root: Optional['TreeNode'], p: int, q: int) -> int:
        def find_lca(root: Optional['TreeNode'], p: int, q: int) -> Optional[TreeNode]:
            if not root or root.val == p or root.val == q:
                return root
            left = find_lca(root.left, p, q)
            right = find_lca(root.right, p, q)
            if not left:
                return right
            if not right:
                return left
            return root

        def dfs(root: Optional['TreeNode'], target: int) -> int:
            if not root:
                return 0x7FFFFFFF
            if root.val == target:
                return 0
            left = dfs(root.left, target)
            right = dfs(root.right, target)
            return min(left, right) + 1

        lca = find_lca(root, p, q)
        p_dep = dfs(lca, p)
        q_dep = dfs(lca, q)
        return p_dep + q_dep

