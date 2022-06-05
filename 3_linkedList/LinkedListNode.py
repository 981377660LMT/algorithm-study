from dataclasses import dataclass
from typing import Generic, Optional, TypeVar

V = TypeVar('V')


@dataclass(slots=True)
class Node(Generic[V]):
    value: V
    left: Optional['Node[V]'] = None
    right: Optional['Node[V]'] = None

    def insertRight(self, node: 'Node[V]') -> 'Node[V]':
        """在 self 后插入 node,并返回该 node"""
        node.left = self
        node.right = self.right
        node.left.right = node
        if node.right:
            node.right.left = node
        return node

    def insertLeft(self, node: 'Node[V]') -> 'Node[V]':
        """在 self 前插入 node,并返回该 node"""
        node.right = self
        node.left = self.left
        node.right.left = node
        if node.left:
            node.left.right = node
        return node

    def remove(self) -> None:
        """从链表里移除自身"""
        if self.left:
            self.left.right = self.right
        if self.right:
            self.right.left = self.left

    def __repr__(self) -> str:
        return f'{self.value}->{self.right}'

