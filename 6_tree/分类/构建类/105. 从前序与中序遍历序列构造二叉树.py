# https://leetcode.cn/problems/construct-binary-tree-from-preorder-and-postorder-traversal/solutions/2649218/tu-jie-cong-on2-dao-onpythonjavacgojsrus-h0o5/


from typing import List, Optional


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def buildTree(self, preorder: List[int], inorder: List[int]) -> Optional[TreeNode]:
        n = len(inorder)
        inMp = {inorder[i]: i for i in range(n)}

        def dfs(preLeft: int, preRight: int, inLeft: int, inRight: int) -> Optional[TreeNode]:
            if preLeft > preRight:
                return None
            rootVal = preorder[preLeft]
            res = TreeNode(rootVal)
            inRootIndex = inMp[rootVal]
            leftSize = inRootIndex - inLeft
            res.left = dfs(preLeft + 1, preLeft + leftSize, inLeft, inRootIndex - 1)
            res.right = dfs(preLeft + leftSize + 1, preRight, inRootIndex + 1, inRight)
            return res

        return dfs(0, n - 1, 0, n - 1)
