from typing import Optional


class TreeNode:
    def __init__(
        self, val: int = 0, left: Optional['TreeNode'] = None, right: Optional['TreeNode'] = None
    ):
        self.val = val
        self.left = left
        self.right = right


# 找到它最深的叶节点的最近公共祖先。
# 总结：最深叶子节点的公共祖先的左右子树高度相同，也就是最深叶子节点的深度一定相同
class Solution:
    def lcaDeepestLeaves(self, root: Optional[TreeNode]) -> Optional[TreeNode]:
        def getDepth(root: Optional[TreeNode]) -> int:
            if not root:
                return 0
            return max(getDepth(root.left), getDepth(root.right)) + 1

        # 其实每次求最深高度可以加个备忘录，这里就不加了
        ld = getDepth(root.left)
        rd = getDepth(root.right)
        if ld == rd:
            return root
        return self.lcaDeepestLeaves(root.left) if ld > rd else self.lcaDeepestLeaves(root.right)
