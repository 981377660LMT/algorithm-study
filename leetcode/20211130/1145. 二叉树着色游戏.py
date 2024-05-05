# 1145. 二叉树着色游戏
# The optimal y must be a neighbor of x. Choose the largest subtree of x. O(n).
# 那能够胜利的条件为

# 阻断的三个选择:
# 左节点个数 > 总数 - 左节点个数
# 右节点个数 > 总数 - 右节点个数
# 父节点个数 > 总数 - 父节点个数


from collections import defaultdict
from typing import Optional


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def btreeGameWinningMove(self, root: Optional[TreeNode], n: int, x: int) -> bool:
        def dfs(root: Optional["TreeNode"]) -> int:
            if not root:
                return 0
            left = dfs(root.left)
            right = dfs(root.right)
            leftCounter[root.val] = left
            rightCounter[root.val] = right
            subCounter[root.val] = left + right + 1
            return left + right + 1

        if not root:
            return False

        leftCounter = defaultdict(int)
        rightCounter = defaultdict(int)
        subCounter = defaultdict(int)
        dfs(root)
        all_ = subCounter[root.val]
        left, right, parent = leftCounter[x], rightCounter[x], all_ - subCounter[x]
        return any(x > all_ - x for x in (left, right, parent))
