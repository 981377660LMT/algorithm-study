# 979. 在二叉树中分配硬币
# 树上移动硬币/金币
# https://leetcode.cn/problems/distribute-coins-in-binary-tree/

from typing import Optional


class TreeNode:
    def __init__(
        self, val: int = 0, left: Optional["TreeNode"] = None, right: Optional["TreeNode"] = None
    ):
        self.val = val
        self.left = left
        self.right = right


# 树中的每个结点上都对应有 node.val 枚硬币，并且总共有 N 枚硬币。
# 1<= N <= 100
# 我们可以选择两个相邻的结点，然后将一枚硬币从其中一个结点移动到另一个结点。(移动可以是从父结点到子结点，或者从子结点移动到父结点。)。
# 返回使每个结点上只有一枚硬币所需的移动次数。


class Solution:
    def distributeCoins(self, root: TreeNode) -> int:
        # !以root为顶点的子树的需要多少硬币才平衡，正数表示我有多余，负数表示我需要别人给我硬币。
        def dfs(root: TreeNode | None) -> int:
            """返回当前子树需要父结点给的硬币数.
            大于0表示当前子树有多余的硬币,小于0表示当前子树需要硬币.
            """
            if not root:
                return 0

            lMoves = dfs(root.left)
            rMoves = dfs(root.right)
            self.res += abs(lMoves) + abs(rMoves)  # root需要向子节点移动的硬币数
            return lMoves + rMoves + root.val - 1  # root所在子树需要硬币还是多余硬币

        self.res = 0
        dfs(root)
        return self.res
