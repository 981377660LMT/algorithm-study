from operator import and_, or_
from typing import List, Optional, Tuple
from collections import defaultdict, Counter, deque

MOD = int(1e9 + 7)
INF = int(1e20)


# Definition for a binary tree node.
class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


MAPPING = {0: False, 1: True, 2: or_, 3: and_}


class Solution:
    def evaluateTree(self, root: Optional[TreeNode]) -> bool:
        """返回根节点 root 的布尔运算值。"""

        def dfs(root: Optional[TreeNode]) -> bool:
            if root is None:
                return True
            if not root.left and not root.right:
                return MAPPING[root.val]
            left = dfs(root.left)
            right = dfs(root.right)
            return MAPPING[root.val](left, right)

        return dfs(root)
