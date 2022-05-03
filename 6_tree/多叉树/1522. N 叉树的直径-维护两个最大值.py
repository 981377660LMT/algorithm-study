from heapq import nlargest
from typing import Optional


class Node:
    def __init__(self, val=None, children=None):
        self.val = val
        self.children = children if children is not None else []


class Solution:
    def diameter(self, root: 'Node') -> int:
        def dfs(root: Optional['Node']) -> int:
            if not root:
                return 0
            cands = [0, 0]
            for next in root.children:
                cands.append(dfs(next))
            max1, max2 = nlargest(2, cands)
            self.res = max(self.res, max1 + max2)
            return max(max1, max2) + 1

        self.res = 0
        dfs(root)
        return self.res
