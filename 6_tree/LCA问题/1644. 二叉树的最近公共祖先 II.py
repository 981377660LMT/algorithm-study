from typing import Optional


class TreeNode:
    def __init__(
        self, val: int = 0, left: Optional['TreeNode'] = None, right: Optional['TreeNode'] = None
    ):
        self.val = val
        self.left = left
        self.right = right


# 返回给定节点 p 和 q 的最近公共祖先（LCA）节点。如果 p 或 q 之一 不存在 于该二叉树中，返回 null
# 在不检查节点是否存在的情况下，你可以遍历树找出最近公共祖先节点吗？

# 1.先判断pq都在不在树里
# 2.后续遍历+递归求LCA


class Solution:
    def lowestCommonAncestor(self, root: 'TreeNode', p: 'TreeNode', q: 'TreeNode') -> 'TreeNode':
        # 先判断p和q是不是都在树中
        self.flag_p, self.flag_q = False, False
        self.checkExist(root, p, q)
        if not (self.flag_p and self.flag_q):
            return None

        return self.findLCA(root, p, q)

    def checkExist(self, root, p, q) -> None:  # 判断p和q都在不在树中
        if root == p:
            self.flag_p = True
        if root == q:
            self.flag_q = True
        if root.left:
            self.checkExist(root.left, p, q)
        if root.right:
            self.checkExist(root.right, p, q)

    def findLCA(self, root, p, q) -> TreeNode:  # 后序遍历，找LCA
        if not root or root == p or root == q:
            return root
        L = self.findLCA(root.left, p, q)
        R = self.findLCA(root.right, p, q)
        if not L:
            return R
        if not R:
            return L
        return root

