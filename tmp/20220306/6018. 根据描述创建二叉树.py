from collections import defaultdict
from typing import List, Optional, Tuple

MOD = int(1e9 + 7)


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def createBinaryTree(self, descriptions: List[List[int]]) -> Optional['TreeNode']:
        def dfs(cur: Optional['TreeNode']) -> None:
            if not cur:
                return
            for next, isLeft in adjMap[cur.val]:
                if isLeft:
                    cur.left = TreeNode(next)
                    dfs(cur.left)
                else:
                    cur.right = TreeNode(next)
                    dfs(cur.right)

        adjMap = defaultdict(list)
        visited = set()
        notRoot = set()
        for u, v, type in descriptions:
            adjMap[u].append((v, type))
            visited |= {u, v}
            notRoot |= {v}

        root = next((r for r in visited if r not in notRoot), -1)

        res = TreeNode(root)
        dfs(res)
        return res


print(
    Solution()
    .createBinaryTree(
        descriptions=[[20, 15, 1], [20, 17, 0], [50, 20, 1], [50, 80, 0], [80, 19, 1]]
    )
    .__dict__
)
