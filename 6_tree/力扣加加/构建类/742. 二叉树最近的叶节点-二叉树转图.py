from typing import Optional
from collections import defaultdict, deque


class TreeNode:
    def __init__(
        self, val: int = 0, left: Optional['TreeNode'] = None, right: Optional['TreeNode'] = None
    ):
        self.val = val
        self.left = left
        self.right = right


# 给定一个 每个结点的值互不相同 的二叉树，和一个目标值 k，找出树中与目标值 k 最近的叶结点。
# 给定的二叉树中有某个结点使得 node.val == k。

# 总结：
# 我们将树转换为图，则我们可以通过宽度优先搜索找到最近的叶子节点
class Solution:
    def findClosestLeaf(self, root: TreeNode, k: int) -> int:
        adjMap = defaultdict(list)

        def dfs(root: Optional[TreeNode], parent: Optional[TreeNode]):
            if not root:
                return

            adjMap[root].append(parent)
            adjMap[parent].append(root)
            dfs(root.left, root)
            dfs(root.right, root)

        dfs(root, None)

        queue = deque(node for node in adjMap if node and node.val == k)
        visited = set(queue)

        while queue:
            cur = queue.popleft()
            if not cur:
                continue
            # 找到了叶子节点(adjMap)
            if len(adjMap[cur]) <= 1:
                return cur.val
            for next in adjMap[cur]:
                if next not in visited:
                    visited.add(next)
                    queue.append(next)


# root = [1, 3, 2], k = 1
# 二叉树图示：
#           1
#          / \
#         3   2

# 输出： 2 (或 3)

# 解释： 2 和 3 都是距离目标 1 最近的叶节点。

