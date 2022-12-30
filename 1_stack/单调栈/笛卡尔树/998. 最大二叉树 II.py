from typing import Optional


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def insertIntoMaxTree(self, root: Optional[TreeNode], val: int) -> Optional[TreeNode]:
        """
        最大树:每个节点的值都大于其子树中任何其他值
        向最大树中插入val 返回插入后的树根节点

        递归插入
        如果root值小于val 那么将root作为新的根节点的左子树
        如果root值大于val 那么val插入到root的右子树中
        """
        if not root or root.val < val:
            return TreeNode(val, root)
        root.right = self.insertIntoMaxTree(root.right, val)
        return root
