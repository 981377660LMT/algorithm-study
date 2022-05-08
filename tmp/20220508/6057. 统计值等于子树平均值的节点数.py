from typing import List, Optional, Tuple


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def averageOfSubtree(self, root: Optional[TreeNode]) -> int:
        def dfs(node: Optional[TreeNode]) -> Tuple[int, int]:
            nonlocal res
            if not node:
                return 0, 0
            leftSum, leftCount = dfs(node.left)
            rightSum, rightCount = dfs(node.right)
            allCount = leftCount + rightCount + 1
            allSum = leftSum + rightSum + node.val
            if node.val == allSum // allCount:
                res += 1
            return allSum, allCount

        res = 0
        dfs(root)
        return res

