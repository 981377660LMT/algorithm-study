from typing import List, Optional


class TreeNode:
    def __init__(
        self, val: int = 0, left: Optional['TreeNode'] = None, right: Optional['TreeNode'] = None
    ):
        self.val = val
        self.left = left
        self.right = right


# 返回 nodes 中所有节点的最近公共祖先（LCA）
class Solution:
    def lowestCommonAncestor(self, root: 'TreeNode', nodes: 'List[TreeNode]') -> 'TreeNode':
        if not root or root in nodes:
            return root
        l = self.lowestCommonAncestor(root.left, nodes)
        r = self.lowestCommonAncestor(root.right, nodes)
        if not l:
            return r
        if not r:
            return l
        return root

