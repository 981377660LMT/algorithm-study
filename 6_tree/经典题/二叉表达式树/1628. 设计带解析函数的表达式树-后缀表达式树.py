from dataclasses import dataclass
from operator import add, floordiv, mul, sub
from typing import List
from abc import ABCMeta, abstractmethod


class Node(metaclass=ABCMeta):
    @abstractmethod
    def evaluate(self) -> int:
        pass


OPT = dict(zip(['+', '-', '*', '/'], [add, sub, mul, floordiv]))


@dataclass(slots=True)
class MyNode(Node):
    val: str
    left: Node | None
    right: Node | None

    def evaluate(self) -> int:
        if self.val.isdigit():
            return int(self.val)
        if self.left and self.right:
            return OPT[self.val](self.left.evaluate(), self.right.evaluate())
        raise ValueError('Invalid expression')


class TreeBuilder(object):
    def buildTree(self, postfix: List[str]) -> 'Node':
        last = postfix.pop()
        root = MyNode(last, None, None)
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
res = expTree.evaluate()
print(res)
