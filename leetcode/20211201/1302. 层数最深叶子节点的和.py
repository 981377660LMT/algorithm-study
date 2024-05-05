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
    def deepestLeavesSum(self, root: Optional["TreeNode"]) -> int:
        queue = deque([(root, 0)])
        maxDepth, res = -1, 0
        while queue:
            cur, depth = queue.popleft()
            if not cur:
                continue
            if depth > maxDepth:
                maxDepth, res = depth, cur.val
            else:
                res += cur.val
            if cur.left:
                queue.append((cur.left, depth + 1))
            if cur.right:
                queue.append((cur.right, depth + 1))

        return res
