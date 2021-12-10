from typing import List, Optional
from collections import Counter


class Node:
    def __init__(self, val: int = 0, left: Optional['Node'] = None, right: Optional['Node'] = None):
        self.val = val
        self.left = left
        self.right = right


# 当两棵二叉搜索树中的变量取任意值，分别求得的值都相等时，我们称这两棵二叉表达式树是等价的。
# 在本题中，我们只考虑 '+' 运算符（即加法）。
class Solution:
    def checkEquivalence(self, root1: 'Node', root2: 'Node') -> bool:
        def dfs(root: Node, path: List[int]):
            if not root:
                return
            if not root.left or not root.right:
                path.append(root.val)
            dfs(root.left, path)
            dfs(root.right, path)

        nums1, nums2 = [], []
        dfs(root1, nums1)
        dfs(root2, nums2)
        return Counter(nums1) == Counter(nums2)


# 输入：root1 = [+,a,+,null,null,b,c], root2 = [+,+,a,b,c]
# 输出：true
# 解释：a + (b + c) == (b + c) + a
