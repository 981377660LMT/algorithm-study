from typing import Optional
from collections import deque


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


# 总结:
# 1.dfs/bfs`覆盖式更新`


class Solution:
    def deepestLeavesSum(self, root: Optional['TreeNode']) -> int:
        q = deque([(root, 0)])
        maxdep, res = -1, 0

        while q:
            node, dep = q.popleft()
            if dep > maxdep:
                maxdep, res = dep, node.val
            else:
                res += node.val
            if node.left:
                q.append((node.left, dep + 1))
            if node.right:
                q.append((node.right, dep + 1))

        return res


print(Solution().deepestLeavesSum())
