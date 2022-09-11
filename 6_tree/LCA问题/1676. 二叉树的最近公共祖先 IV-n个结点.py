from collections import defaultdict
from typing import Generator, List, Optional


class TreeNode:
    def __init__(
        self, val: int = 0, left: Optional['TreeNode'] = None, right: Optional['TreeNode'] = None
    ):
        self.val = val
        self.left = left
        self.right = right


# !O(n)求二叉树多个结点的最近公共祖先
# !可以推广到求n叉树k个结点的LCA
class Solution:
    def lowestCommonAncestor(self, root: 'TreeNode', nodes: 'List[TreeNode]') -> 'TreeNode':
        def dfs(root: Optional['TreeNode']) -> Generator['TreeNode', None, None]:
            """后序dfs和从下往上拓扑排序 都是等价的
            看哪个点最先为k
            """
            if not root:
                return

            subSum[root.val] += int(root in needs)

            for child in [root.left, root.right]:
                if not child:
                    continue
                yield from dfs(child)
                subSum[root.val] += subSum[child.val]

            if subSum[root.val] == len(nodes):
                yield root

        needs = set(nodes)
        subSum = defaultdict(int)
        return next(dfs(root))

