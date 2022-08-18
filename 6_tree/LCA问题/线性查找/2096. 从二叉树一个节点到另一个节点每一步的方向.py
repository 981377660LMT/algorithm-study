from collections import defaultdict
from typing import DefaultDict, Optional, Set


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


# 5944. 从二叉树一个节点到另一个节点每一步的方向
# 寻找LCA
class Solution:
    def getDirections(self, root: Optional[TreeNode], startValue: int, destValue: int) -> str:
        """求LCA和两点间的路径只需要depth和parent信息"""

        def dfs(root: Optional[TreeNode], dep: int) -> None:
            if not root:
                return
            depth[root.val] = dep
            if root.left:
                parent[root.left.val] = (root.val, "L")
                dfs(root.left, dep + 1)
            if root.right:
                parent[root.right.val] = (root.val, "R")
                dfs(root.right, dep + 1)

        def getLCA(root1: int, root2: int) -> int:
            if depth[root1] < depth[root2]:
                root1, root2 = root2, root1
            diff = depth[root1] - depth[root2]
            for _ in range(diff):
                root1 = parent[root1][0]
            while root1 != root2:
                root1 = parent[root1][0]
                root2 = parent[root2][0]
            return root1

        parent, depth = defaultdict(tuple), defaultdict(lambda: -1)
        dfs(root, 0)
        lca = getLCA(startValue, destValue)

        res1 = (depth[startValue] - depth[lca]) * "U"
        res2 = []

        while destValue != lca:
            destValue, direction = parent[destValue]
            res2.append(direction)

        return res1 + ''.join(res2[::-1])

