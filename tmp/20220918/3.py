from typing import List, Tuple, Optional
from collections import defaultdict, Counter, deque
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一棵 完美 二叉树的根节点 root ，请你反转这棵树中每个 奇数 层的节点值。
# Definition for a binary tree node.
class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


# 6182. 反转二叉树的奇数层
class Solution:
    def reverseOddLevels(self, root: Optional[TreeNode]) -> Optional[TreeNode]:
        queue = deque([root])
        level = 0
        while queue:
            nQueue = deque()
            size = len(queue)
            for _ in range(size):
                node = queue.popleft()
                if node.left:
                    nQueue.append(node.left)
                if node.right:
                    nQueue.append(node.right)
            level += 1
            if level & 1:
                for i in range(len(nQueue) // 2):
                    nQueue[i].val, nQueue[~i].val = nQueue[~i].val, nQueue[i].val
            queue = nQueue
        return root
