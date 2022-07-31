from collections import defaultdict
from typing import List, Optional

INF = int(1e20)


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def canMerge(self, trees: List[TreeNode]) -> Optional[TreeNode]:
        """1.是二叉搜索树,2.每个叶子结点都用到

        如果 n 棵树可以合并成一棵二叉搜索树，那么这 n 棵树一定是从一棵二叉搜索树中 拆分 而来。
        因此拆分后，所有的树的 叶节点 的值 互不相同
        而题目中又给了 “所有根节点值互不相同”，那么，“叶节点” 和 “根节点” 之间的 “对应关系” 应当是唯一的
        其实这种套路在 1719. 重构一棵树的方案数 也出现过，我们可以先考虑`构造好的树有什么性质`，然后再反推解题。
        """

        # 建树，哈希表存叶子结点
        n = len(trees)
        leaves = defaultdict(list)
        for root in trees:
            if root.left:
                val = root.left.val
                if val in leaves:
                    return None
                leaves[val] = [root, "L", False]  # metadata
            if root.right:
                val = root.right.val
                if val in leaves:
                    return None
                leaves[val] = [root, "R", False]

        # 遍历根执行合并
        visited, count = [False] * n, 0
        for i, root in enumerate(trees):
            val = root.val
            if val in leaves:
                parent, direction, _ = leaves[val]
                if direction == "L":
                    parent.left = root
                else:
                    parent.right = root
                visited[i] = True
                count += 1
        if count != n - 1:
            return None

        # 验证二叉搜索树
        isValidBST = True
        root = trees[next(i for i, v in enumerate(visited) if not v)]

        def dfs(root: Optional[TreeNode], lower: int, upper: int) -> None:
            nonlocal isValidBST  # 不能写global  因为这样global后会变成scoped variable
            if not root:
                return
            if lower >= root.val or upper <= root.val:
                isValidBST = False
                return
            if root.left:
                leaves[root.left.val][2] = True
                dfs(root.left, lower, root.val)
            if root.right:
                leaves[root.right.val][2] = True
                dfs(root.right, root.val, upper)

        dfs(root, -INF, INF)
        if not isValidBST:
            return None

        # 验证是否所有叶子节点都被遍历到
        if all(v[2] for v in leaves.values()):
            return root
        return None


R5 = TreeNode(5)
R3 = TreeNode(3)
R8 = TreeNode(8)
R5.left = R3
R5.right = R8
R3 = TreeNode(3)
R2 = TreeNode(2)
R6 = TreeNode(6)
R3.left = R2
R3.right = R6

print(Solution().canMerge([R5, R3]))
