from typing import Optional
from re import findall


class Node:
    def __init__(self, val='0', left: Optional['Node'] = None, right: Optional['Node'] = None):
        self.val = val
        self.left = left
        self.right = right


# 叶节点（有 0 个子节点的节点）表示操作数，
# 非叶节点（有 2 个子节点的节点）表示运算符： '+' （加）、 '-' （减）、 '*' （乘）和 '/' （除）。


class Solution:
    """根据中缀表达式构造二叉表达式树"""

    def expTree(self, s: str) -> 'Node':
        weight = dict(zip('(+-*/', [0, 1, 1, 2, 2]))
        numStack, optStack = [], []

        s += ')'

        for char in s:
            if char == '(':
                optStack.append(char)
            elif char.isdigit():
                numStack.append(Node(char))
            elif char in '+-*/':
                while optStack and weight[optStack[-1]] >= weight[char]:
                    node2, node1 = numStack.pop(), numStack.pop()
                    numStack.append(Node(optStack.pop(), node1, node2))
                optStack.append(char)
            elif char == ')':
                while optStack and optStack[-1] != '(':
                    node2, node1 = numStack.pop(), numStack.pop()
                    numStack.append(Node(optStack.pop(), node1, node2))
                if optStack:
                    optStack.pop()

        return numStack[0]


print(Solution().expTree(s="3*4-2*5"))
print(findall(r'(\()|(\d+)|([-+*/])|(\))', "3*4-2*5+"))

