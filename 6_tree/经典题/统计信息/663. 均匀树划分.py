from typing import Optional
from collections import defaultdict


class TreeNode:
    def __init__(
        self, val: int = 0, left: Optional['TreeNode'] = None, right: Optional['TreeNode'] = None
    ):
        self.val = val
        self.left = left
        self.right = right


# 你的任务是检查是否可以通过去掉树上的一条边将树分成两棵，且这两棵树结点之和相等。
# 树上结点的权值范围 [-100000, 100000]。
class Solution:
    def checkEqualTree(self, root: TreeNode) -> bool:
        counter = defaultdict(int)

        def dfs(root: Optional[TreeNode]) -> int:
            if not root:
                return 0
            leftSum = dfs(root.left)
            rightSum = dfs(root.right)
            total = leftSum + rightSum + root.val
            counter[total] += 1
            return total

        total = dfs(root)
        # 检查整棵树的一半权值是否出现（但不能是整棵树之和）。
        counter[total] -= 1
        target = total / 2
        return counter[target] >= 1

