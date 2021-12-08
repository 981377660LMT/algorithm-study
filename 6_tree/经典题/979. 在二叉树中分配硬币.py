from typing import Optional


class TreeNode:
    def __init__(
        self, val: int = 0, left: Optional['TreeNode'] = None, right: Optional['TreeNode'] = None
    ):
        self.val = val
        self.left = left
        self.right = right


# 树中的每个结点上都对应有 node.val 枚硬币，并且总共有 N 枚硬币。
# 1<= N <= 100
# 我们可以选择两个相邻的结点，然后将一枚硬币从其中一个结点移动到另一个结点。(移动可以是从父结点到子结点，或者从子结点移动到父结点。)。
# 返回使每个结点上只有一枚硬币所需的移动次数。

# 不知道就后序dfs
class Solution:
    def distributeCoins(self, root: TreeNode) -> int:
        self.res = 0

        # 以root为顶点的子树的需要多少硬币才平衡，正数表示我有多余，负数表示我需要别人给我硬币。
        def dfs(root: TreeNode) -> int:
            if not root:
                return 0

            lMoves = dfs(root.left)
            rMoves = dfs(root.right)
            self.res += abs(lMoves + rMoves + root.val - 1)
            return lMoves + rMoves + root.val - 1

        dfs(root)
        return self.res

