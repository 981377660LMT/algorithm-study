from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一棵二叉树的根 root ，请你将每个节点的值替换成该节点的所有 堂兄弟节点值的和 。

# 如果两个节点在树中有相同的深度且它们的父节点不同，那么它们互为 堂兄弟 。

# 请你返回修改值之后，树的根 root 。


# 注意，一个节点的深度指的是从树根节点到这个节点经过的边数。
# Definition for a binary tree node.
class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def replaceValueInTree(self, root: Optional[TreeNode]) -> Optional[TreeNode]:
        def dfs(cur: Optional[TreeNode], dep: int) -> None:
            if not cur:
                return
            levelSum[dep] += cur.val
            curId = id(cur)
            if cur.left:
                parent[id(cur.left)] = curId
                dfs(cur.left, dep + 1)
                childSum[curId] += cur.left.val
            if cur.right:
                parent[id(cur.right)] = curId
                dfs(cur.right, dep + 1)
                childSum[curId] += cur.right.val

        levelSum = defaultdict(int)
        parent = defaultdict(int)  # id->parentId
        childSum = defaultdict(int)  # id->sum
        dfs(root, 0)

        def dfs2(cur: Optional[TreeNode], dep: int) -> None:
            if not cur:
                return
            curId = id(cur)
            if curId in parent:
                p = parent[curId]
                cur.val = levelSum[dep] - childSum[p]
            else:
                cur.val = 0
            dfs2(cur.left, dep + 1)
            dfs2(cur.right, dep + 1)

        dfs2(root, 0)
        return root
