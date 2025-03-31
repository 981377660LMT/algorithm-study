# https://leetcode.cn/problems/construct-binary-tree-from-preorder-and-postorder-traversal/solutions/2649218/tu-jie-cong-on2-dao-onpythonjavacgojsrus-h0o5/
#
# !1. 构建一个哈希表，将中序遍历中每个值及其对应的索引记录下来，这样在递归构造树时可以快速定位根节点在中序数组中的位置。
# !2. 再利用后序遍历的特点：后序数组的最后一个元素即为当前子树的根节点，然后根据中序数组将数组划分为左右子树，递归构造即可。

from typing import List, Optional


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def buildTree(self, inorder: List[int], postorder: List[int]) -> Optional[TreeNode]:
        n = len(inorder)
        inMp = {inorder[i]: i for i in range(n)}

        def dfs(inLeft: int, inRight: int, postLeft: int, posRight: int) -> Optional[TreeNode]:
            if inLeft > inRight:
                return None
            rootVal = postorder[posRight]
            res = TreeNode(rootVal)
            inRootIndex = inMp[rootVal]
            leftSize = inRootIndex - inLeft

            res.left = dfs(inLeft, inRootIndex - 1, postLeft, postLeft + leftSize - 1)
            res.right = dfs(inRootIndex + 1, inRight, postLeft + leftSize, posRight - 1)

            return res

        return dfs(0, n - 1, 0, n - 1)
