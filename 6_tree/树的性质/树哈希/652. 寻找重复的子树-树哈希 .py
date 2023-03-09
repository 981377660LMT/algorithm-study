# 如果两棵树具有 相同的结构 和 相同的结点值 ，则认为二者是 重复 的。
# 对于同一类的重复子树，你只需要返回其中任意 一棵 的根结点即可。
# !注意题目要求子树对应位置也要一样
# !n<=5000
# !-200<=Node.val<=200


# !使用`哈希值的编号`来代替很长的哈希字符串 减少哈希值长度
# Definition for a binary tree node.
from collections import defaultdict
import itertools
from typing import List, Optional


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


# !1.用一个三元组直接表示一棵子树(根节点的值,左子树,右子树)
# !2.向上返回时不返回三元组,而是返回一个哈希值的编号


class Solution:
    def findDuplicateSubtrees(self, root: Optional[TreeNode]) -> List[Optional[TreeNode]]:
        def dfs(node: Optional[TreeNode]) -> Optional[int]:
            if node is None:
                return
            hash_ = (node.val, dfs(node.left), dfs(node.right))
            hashId = pool[hash_]
            counter[hashId].append(node)
            return hashId

        gen = itertools.count()
        pool = defaultdict(lambda: next(gen))  # 全局自增id
        counter = defaultdict(list)
        dfs(root)
        return [nodes[0] for nodes in counter.values() if len(nodes) > 1]
