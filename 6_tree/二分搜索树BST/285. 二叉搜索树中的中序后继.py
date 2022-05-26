# Definition for a binary tree node.
from itertools import chain, dropwhile, takewhile
from typing import Generator, Optional


class TreeNode:
    def __init__(self, x):
        self.val = x
        self.left = None
        self.right = None


class Solution:
    def inorderSuccessor(self, root: 'TreeNode', p: 'TreeNode') -> Optional['TreeNode']:
        def inorder(node: Optional['TreeNode']) -> Generator['TreeNode', None, None]:
            if not node:
                return
            yield from inorder(node.left)
            yield node
            yield from inorder(node.right)

        gen = dropwhile(lambda node: node.val != p.val, inorder(root))
        next(gen)
        return next(gen, None)

    def inorderSuccessor2(self, root: 'TreeNode', p: 'TreeNode') -> Optional['TreeNode']:
        dfs = lambda node: chain(dfs(node.left), [node], dfs(node.right)) if node else []
        gen = dropwhile(lambda node: node.val != p.val, dfs(root))
        return next(gen, None) if next(gen, None) else None  # 解决了p不在树中的情况
