from typing import List, Optional


class TreeNode:
    def __init__(
        self, val: int = 0, left: Optional['TreeNode'] = None, right: Optional['TreeNode'] = None
    ):
        self.val = val
        self.left = left
        self.right = right


# 给你一棵二叉树的根节点 root ，按顺序返回组成二叉树 边界 的这些值

# 二叉树的”三“视图，左视图+底视图（即全部叶子节点）+右视图
class Solution:
    def boundaryOfBinaryTree(self, root: TreeNode) -> List[int]:
        if root and not root.left and not root.right:
            return [root.val]
        leftView = []
        bottomView = []
        rightView = []

        def getLeftView(root,):
            if not root:
                return
            leftView.append(root.val)
            if root.left:
                getLeftView(root.left)
            elif root.right:
                getLeftView(root.right)

        def getRightView(root,):
            if not root:
                return
            rightView.append(root.val)
            if root.right:
                getRightView(root.right)
            elif root.left:
                getRightView(root.left)

        def getBottomView(root,):
            if not root:
                return
            if not root.left and not root.right:
                bottomView.append(root.val)
            else:
                getBottomView(root.left)
                getBottomView(root.right)

        getLeftView(root.left)
        getBottomView(root)
        getRightView(root.right)

        return [root.val] + leftView[:-1] + bottomView + rightView[::-1][1:]
