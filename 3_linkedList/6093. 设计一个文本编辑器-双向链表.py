from dataclasses import dataclass
from typing import Generic, Optional, TypeVar

V = TypeVar('V')


@dataclass(slots=True)
class Node(Generic[V]):
    value: V
    left: Optional['Node[V]'] = None
    right: Optional['Node[V]'] = None

    def insertAfter(self, node: 'Node[V]') -> 'Node[V]':
        """在 self 后插入 node,并返回该 node"""
        node.left = self
        node.right = self.right
        node.left.right = node
        if node.right:
            node.right.left = node
        return node

    def insertBefore(self, node: 'Node[V]') -> 'Node[V]':
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
        self.left = None
        self.right = None

    def __repr__(self) -> str:
        return f'{self.value}->{self.right}'


class TextEditor:
    def __init__(self):
        self.root = Node('')  # 哨兵以及初始化双向链表
        self.root.left = self.root
        self.root.right = self.root  # !初始化双向链表，下面判断节点的 next 若为 self.root，则表示 next 为空
        self.pos = self.root

    def addText(self, text: str) -> None:
        for char in text:
            self.pos = self.pos.insertAfter(Node(char))

    def deleteText(self, k: int) -> int:
        remain = k
        while remain and self.pos != self.root:
            self.pos = self.pos.left
            self.pos.right.remove()
            remain -= 1
        return k - remain

    def cursorLeft(self, k: int) -> str:
        while k and self.pos != self.root:
            self.pos = self.pos.left
            k -= 1
        return self._getText()

    def cursorRight(self, k: int) -> str:
        while k and self.pos.right != self.root:
            self.pos = self.pos.right
            k -= 1
        return self._getText()

    def _getText(self) -> str:
        res = []
        remain, cur = 10, self.pos
        while remain and cur != self.root:
            res.append(cur.value)
            cur = cur.left
            remain -= 1
        return ''.join(res[::-1])
