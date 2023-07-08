"""
判断二叉树序列是否为某棵二叉树的前序遍历序列.

用一个栈模拟前序遍历过程，只要每个非根节点都找得到父节点就说明序列没问题。
假设-1是个真实存在的节点可以避免特判根节点。时空复杂度O(n)。
"""


from typing import List


class Solution:
    def isPreorder(self, nodes: List[List[int]]) -> bool:
        """
        判断序列(node,parent)是否为一棵二叉树的前序遍历序列.
        保证输入是合法的二叉树.
        """
        stack = [-1]
        for cur, pre in nodes:
            while stack and stack[-1] != pre:
                stack.pop()
            if not stack:
                return False
            stack.append(cur)
        return True
