# 1612. 检查两棵二叉表达式树是否等价
# 当两棵二叉搜索树中的变量取任意值，分别求得的值都相等时，我们称这两棵二叉表达式树是等价的。
# 在本题中，我们只考虑 '+' 运算符（即加法）。

from typing import List, Optional


class Node(object):
    def __init__(self, val=" ", left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def checkEquivalence(self, root1: "Node", root2: "Node") -> bool:
        def dfs(node: Optional["Node"], sign: int, counter: List[int]) -> None:
            if node is None:
                return
            if node.val == "+":
                dfs(node.left, sign, counter)
                dfs(node.right, sign, counter)
            elif node.val == "-":
                dfs(node.left, sign, counter)
                dfs(node.right, -sign, counter)
            else:  # 变量
                counter[ord(node.val) - 97] += sign

        cnt1 = [0] * 26
        cnt2 = [0] * 26
        dfs(root1, +1, cnt1)
        dfs(root2, +1, cnt2)
        return cnt1 == cnt2


# 输入：root1 = [+,a,+,null,null,b,c], root2 = [+,+,a,b,c]
# 输出：true
# 解释：a + (b + c) == (b + c) + a
