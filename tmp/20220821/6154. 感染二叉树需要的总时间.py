from typing import Optional
from collections import defaultdict, deque


# Definition for a binary tree node.
class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def amountOfTime(self, root: Optional["TreeNode"], start: int) -> int:
        def dfs(cur: Optional["TreeNode"], parent: int) -> None:
            if cur is None:
                return
            if parent != -1:
                adjMap[parent].add(cur.val)
                adjMap[cur.val].add(parent)
            dfs(cur.left, cur.val)
            dfs(cur.right, cur.val)

        adjMap = defaultdict(set)
        dfs(root, -1)

        queue = deque([start])
        visited = set([start])
        res = -1
        while queue:
            len_ = len(queue)
            for _ in range(len_):
                cur = queue.popleft()
                for next in adjMap[cur]:
                    if next not in visited:
                        visited.add(next)
                        queue.append(next)
            res += 1
        return res
