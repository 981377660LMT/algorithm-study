from typing import Optional


class TreeNode:
    def __init__(
        self, val: int = 0, left: Optional['TreeNode'] = None, right: Optional['TreeNode'] = None
    ):
        self.val = val
        self.left = left
        self.right = right


# 找出存在于 不同 节点 A 和 B 之间的最大值 V，其中 V = |A.val - B.val|，且 A 是 B 的祖先。
# 思路:dfs 维护max和min即可
class Solution:
    def maxAncestorDiff(self, root: TreeNode) -> int:
        res = 0

        def dfs(root):
            if not root:
                return float('inf'), -float('inf')
            lmin, lmax = dfs(root.left)
            rmin, rmax = dfs(root.right)
            rootmin = min(root.val, lmin, rmin)
            rootmax = max(root.val, lmax, rmax)
            nonlocal res
            res = max(res, abs(root.val - rootmin), abs(root.val - rootmax))
            return rootmin, rootmax

        dfs(root)
        return res

