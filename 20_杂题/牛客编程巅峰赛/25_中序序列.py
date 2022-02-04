# 给定一棵有n个结点的二叉树的先序遍历与后序遍历序列，求其中序遍历序列。
import sys
from typing import Generator, List, Optional

sys.setrecursionlimit(100000)

# python 默认的递归限制被设置为1000，这题需要修改
# 1≤n≤100,000
class TreeNode(object):
    def __init__(self, x: int, /, *, left: 'TreeNode' = None, right: 'TreeNode' = None):
        self.val = x
        self.left = left
        self.right = right


class Solution:
    def solve(self, n: int, pre: List[int], suf: List[int]) -> List[int]:
        """给定一棵有n个结点的二叉树的先序遍历与后序遍历序列，求其中序遍历序列。"""

        def build(pre: List[int], suf: List[int]) -> Optional[TreeNode]:
            if not pre:
                return None
            if len(pre) == 1:
                return TreeNode(pre[0])
            root = TreeNode(pre[0])
            # 前序遍历第一个一定是左子树
            index = suf.index(pre[1])
            root.left = build(pre[1 : index + 2], suf[: index + 1])
            root.right = build(pre[index + 2 :], suf[index + 1 : -1])
            return root

        def inorder(root: Optional[TreeNode]) -> Generator[int, None, None]:
            if not root:
                return
            yield from inorder(root.left)
            yield root.val
            yield from inorder(root.right)

        return list(inorder(build(pre, suf)))
