"""
反转表达式值 树形dp枚举所有情况
"""

# Definition for a binary tree node.
from operator import and_, not_, or_, xor
from typing import List, Optional, Tuple


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


MAPPING = {0: False, 1: True, 2: or_, 3: and_, 4: xor, 5: not_}
INF = int(1e20)


class Solution:
    def minimumFlips(self, root: Optional["TreeNode"], result: bool) -> int:
        """树的叶子为操作数，树的非叶子为运算符，求叶子节点的最小翻转次数"""

        def dfs(root: Optional["TreeNode"]) -> List[int]:
            """返回(变为FALSE的最小操作次数,变为TRUE的最小操作次数)"""
            if not root:
                return [INF, INF]
            if root.val in (0, 1):
                return [int(root.val == 1), int(root.val == 0)]
            if root.val == 5:
                return dfs(root.left)[::-1] if root.left else dfs(root.right)[::-1]

            res, leftRes, rightRes = [INF, INF], dfs(root.left), dfs(root.right)
            for left, leftFlip in enumerate(leftRes):
                for right, rigthFlip in enumerate(rightRes):
                    value = MAPPING[root.val](left, right)
                    res[value] = min(res[value], leftFlip + rigthFlip)

            return res

        return dfs(root)[result]
