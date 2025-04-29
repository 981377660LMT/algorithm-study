from typing import List, Optional


class TreeNode:
    def __init__(
        self, x: int, left: Optional["TreeNode"] = None, right: Optional["TreeNode"] = None
    ):
        self.val = x
        self.left = left
        self.right = right


class Solution:
    def constructMaximumBinaryTree(self, nums: List[int]) -> Optional[TreeNode]:
        """
        单调递减栈：
        对于遍历到的每个新值 x，构造节点 cur：
          1. 不断弹出栈顶所有比 x 小的节点，它们都应该成为 cur.left 最右侧的那一棵子树。
          2. 如果此时栈非空，说明栈顶比 x 大，cur 应成为它的 right 子树。
          3. 将 cur 入栈。
        最终栈底元素即为根节点。
        """
        stack: List[TreeNode] = []
        for x in nums:
            cur = TreeNode(x)
            while stack and stack[-1].val < x:
                popped = stack.pop()
                cur.left = popped
            if stack:
                stack[-1].right = cur
            stack.append(cur)
        return stack[0] if stack else None
