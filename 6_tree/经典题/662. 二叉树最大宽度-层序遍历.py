from collections import deque
from typing import Optional

# Definition for a binary tree node.


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


# 题目数据保证答案将会在  32 位 带符号整数范围内。
class Solution:
    def widthOfBinaryTree(self, root: Optional["TreeNode"]) -> int:
        """层序遍历维护每层的信息"""
        if not root:
            return 0

        res = 1
        queue = deque([(root, 1)])
        while queue:
            nextQueue = deque()
            len_ = len(queue)
            for _ in range(len_):
                cur, index = queue.popleft()
                if cur.left:
                    nextQueue.append((cur.left, index << 1))
                if cur.right:
                    nextQueue.append((cur.right, index << 1 | 1))
            queue = nextQueue

            if queue:
                res = max(res, queue[-1][1] - queue[0][1] + 1)

        return res
