# https://leetcode.cn/problems/extract-kth-character-from-the-rope-tree/
# 获取Rope的第K个字符


from typing import Optional


# Definition for a rope tree node.
class RopeTreeNode:
    def __init__(self, len=0, val="", left=None, right=None):
        self.len = len
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def getKthCharacter(self, root: Optional[RopeTreeNode], k: int) -> str:
        def at(root: Optional[RopeTreeNode], pos: int) -> str:
            if not root:
                return ""
            if root.len == 0:  # 叶子节点
                return root.val[pos]
            leftLen = 0
            if root.left:
                leftLen = root.left and root.left.len or len(root.left.val)
            if pos < leftLen:
                return at(root.left, pos)
            return at(root.right, pos - leftLen)

        return at(root, k - 1)
