from typing import Tuple


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


# 交错路径的长度定义为：访问过的节点数目 - 1（单个节点的路径长度为 0 ）。
# 请你返回给定树中最长 交错路径 的长度。

# 要求子节点对父节点的贡献
class Solution:
    def longestZigZag(self, root: TreeNode) -> int:

        self.res = 0

        def dfs(root: TreeNode) -> Tuple[int, int]:
            if not root:
                return (0, 0)
            l1, r1 = dfs(root.left)
            l2, r2 = dfs(root.right)
            # 若左子树存在，则选左子树向右拐的最大路径
            L = r1 + 1 if root.left else 0
            # 若右子树存在，则选右子树向左拐的最大路径
            R = l2 + 1 if root.right else 0
            # 每个节点上更新最长路径
            self.res = max(self.res, L, R)

            return (L, R)

        dfs(root)

        return self.res

