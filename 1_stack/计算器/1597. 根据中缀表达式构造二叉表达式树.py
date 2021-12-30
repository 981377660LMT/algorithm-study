from typing import Optional
from re import findall


class Node:
    def __init__(self, val: int = 0, left: Optional['Node'] = None, right: Optional['Node'] = None):
        self.val = val
        self.left = left
        self.right = right


# 叶节点（有 0 个子节点的节点）表示操作数，
# 非叶节点（有 2 个子节点的节点）表示运算符： '+' （加）、 '-' （减）、 '*' （乘）和 '/' （除）。


class Solution:
    def expTree(self, s: str) -> 'Node':
        weight = dict(zip('(+-*/', [0, 1, 1, 2, 2]))
        nodeStack, optStack = [], []

        s += ')'

        for left, num, opt, right in findall(r'(\()|(\d+)|([-+*/])|(\))', s):
            if left:
                optStack.append(left)
            elif num:
                nodeStack.append(Node(num))
            elif opt:
                while optStack and weight[optStack[-1]] >= weight[opt]:
                    node1, node2 = nodeStack.pop(), nodeStack.pop()
                    nodeStack.append(Node(optStack.pop(), node2, node1))
                optStack.append(opt)
            else:
                while optStack and optStack[-1] != '(':
                    node1, node2 = nodeStack.pop(), nodeStack.pop()
                    nodeStack.append(Node(optStack.pop(), node2, node1))
                optStack and optStack.pop()

        return nodeStack[0]


print(Solution().expTree(s="3*4-2*5"))
print(findall(r'(\()|(\d+)|([-+*/])|(\))', "3*4-2*5+"))
