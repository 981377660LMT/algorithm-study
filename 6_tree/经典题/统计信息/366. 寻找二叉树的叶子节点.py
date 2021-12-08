from typing import List, Optional
from collections import defaultdict


class TreeNode:
    def __init__(
        self, val: int = 0, left: Optional['TreeNode'] = None, right: Optional['TreeNode'] = None
    ):
        self.val = val
        self.left = left
        self.right = right


# 给你一棵二叉树，请按以下要求的顺序收集它的全部节点：

# 依次从左到右，每次收集并删除所有的叶子节点
# 重复如上过程直到整棵树为空

# 总结：
# 根据`子树的高度`来分组
class Solution:
    def findLeaves(self, root: TreeNode) -> List[List[int]]:
        nodes = defaultdict(list)

        def dfs(root: TreeNode) -> int:
            if not root:
                return 0
            l = dfs(root.left)
            r = dfs(root.right)
            height = max(l, r)
            nodes[height].append(root.val)
            return height + 1

        dfs(root)

        return [list(n) for n in nodes.values()]

