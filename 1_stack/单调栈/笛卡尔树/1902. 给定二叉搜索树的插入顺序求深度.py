# 给定二叉搜索树的插入顺序求深度

from typing import List
from 笛卡尔树 import buildCartesianTree2


class Solution:
    def maxDepthBST(self, insertNums: List[int]) -> int:
        def preOrder(insertIndex: int, dep: int) -> None:
            """前序遍历输出插入序列形成的的BST"""
            nonlocal res
            res = max(res, dep)
            if leftChild[insertIndex] != -1:
                preOrder(leftChild[insertIndex], dep + 1)
            if rightChild[insertIndex] != -1:
                preOrder(rightChild[insertIndex], dep + 1)

        mp = {num: i for i, num in enumerate(insertNums, 1)}
        newNums = [mp[i] for i in range(1, len(insertNums) + 1)]
        rootIndex, leftChild, rightChild = buildCartesianTree2(newNums)
        res = 0
        preOrder(rootIndex, 1)
        return res
