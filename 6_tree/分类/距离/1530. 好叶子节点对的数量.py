from itertools import product
from typing import List, Optional


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


# 如果二叉树中两个 `叶` 节点之间的 最短路径长度 小于或者等于 distance ，那它们就可以构成一组 好叶子节点对 。
# 返回树中 好叶子节点对的数量 。


# tree 的节点数在 [1, 2^10] 范围内
# 总结：要求到叶节点的距离，则dfs返回距离数组即可
class Solution:
    def countPairs(self, root: TreeNode, distance: int) -> int:
        def dfs(root: Optional[TreeNode]) -> List[int]:
            nonlocal res
            if not root:
                return []
            if not root.left and not root.right:
                return [1]
            left = dfs(root.left)
            right = dfs(root.right)
            res += sum(a + b <= distance for a, b in product(left, right))
            return [n + 1 for n in left + right if n + 1 < distance]

        res = 0
        dfs(root)
        return res


# print(Solution().countPairs())
