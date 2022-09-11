from typing import List, Optional


class TreeNode:
    def __init__(
        self, val: int, left: Optional["TreeNode"] = None, right: Optional["TreeNode"] = None
    ):
        self.val = val
        self.left = left
        self.right = right


# !带重复结点的前序中序二叉树
class Solution:
    def getBinaryTrees(self, preOrder: List[int], inOrder: List[int]) -> List[TreeNode]:
        n = len(preOrder)
        if n == 0:
            return [None]  # type: ignore

        res = []
        for i in range(n):
            if inOrder[i] == preOrder[0]:
                leftList = self.getBinaryTrees(preOrder[1 : i + 1], inOrder[:i])
                rightList = self.getBinaryTrees(preOrder[i + 1 :], inOrder[i + 1 :])
                for left in leftList:
                    for right in rightList:
                        root = TreeNode(preOrder[0])
                        root.left = left
                        root.right = right
                        res.append(root)

        return res
