from typing import Optional


# Definition for a binary tree node.
class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


# 你需要将灯饰逐一插入装饰树中，要求如下：
# 完成装饰的二叉树根结点与 root 的根结点值相同
# !若一个节点拥有父节点，则在该节点和他的父节点之间插入一个灯饰（即插入一个值为 -1 的节点）。具体地：
# 在一个 父节点 x 与其左子节点 y 之间添加 -1 节点， 节点 -1、节点 y 为各自父节点的左子节点，
# 在一个 父节点 x 与其右子节点 y 之间添加 -1 节点， 节点 -1、节点 y 为各自父节点的右子节点，
# 现给定二叉树的根节点 root ，请返回完成装饰后的树的根节点。
class Solution:
    def expandBinaryTree(self, root: Optional["TreeNode"]) -> Optional["TreeNode"]:
        def dfs(root: Optional["TreeNode"]) -> None:
            if not root:
                return
            if root.left:
                root.left = TreeNode(-1, root.left)
                dfs(root.left.left)
            if root.right:
                root.right = TreeNode(-1, None, root.right)
                dfs(root.right.right)

        dfs(root)
        return root


# [7, -1, -1, 5, null, null, 6]
# [3,-1,-1,1,null,null,7,-1,-1,null,-1,3,null,null,8,null,4]
