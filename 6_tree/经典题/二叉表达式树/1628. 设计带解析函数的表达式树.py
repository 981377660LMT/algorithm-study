from typing import List

"""
This is the interface for the expression tree Node.
You should not remove it, and you can define some classes to implement it.
"""


class Node:
    def __init__(self, val):
        self.val = val
        self.left: Node = None
        self.right = None

    def evaluate(self) -> int:
        if self.val.isdigit():
            return int(self.val)
        if self.val == '+':
            return self.left.evaluate() + self.right.evaluate()
        if self.val == '-':
            return self.left.evaluate() - self.right.evaluate()
        if self.val == '*':
            return self.left.evaluate() * self.right.evaluate()
        if self.val == '/':
            return self.left.evaluate() // self.right.evaluate()


"""    
This is the TreeBuilder class.
You can treat it as the driver code that takes the postinfix input
and returns the expression tree represnting it as a Node.
"""


class TreeBuilder(object):
    def buildTree(self, postfix: List[str]) -> 'Node':
        last = postfix.pop()
        root = Node(last)
        if root.val.isdigit():
            return root
        root.right = self.buildTree(postfix)
        root.left = self.buildTree(postfix)
        return root


"""
Your TreeBuilder object will be instantiated and called as such:
obj = TreeBuilder();
expTree = obj.buildTree(postfix);
ans = expTree.evaluate();
"""
obj = TreeBuilder()
expTree = obj.buildTree(["3", "4", "+", "2", "*", "7", "/"])
ans = expTree.evaluate()
print(ans)
