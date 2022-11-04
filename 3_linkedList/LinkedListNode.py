from dataclasses import dataclass
from typing import Generic, Optional, TypeVar

V = TypeVar("V")


@dataclass(slots=True)
class Node(Generic[V]):
    value: V
    pre: Optional["Node[V]"] = None
    next: Optional["Node[V]"] = None

    def insertAfter(self, node: "Node[V]") -> "Node[V]":
        """在 self 后插入 node,并返回该 node"""
        node.pre = self
        node.next = self.next
        node.pre.next = node
        if node.next:
            node.next.pre = node
        return node

    def insertBefore(self, node: "Node[V]") -> "Node[V]":
        """在 self 前插入 node,并返回该 node"""
        node.next = self
        node.pre = self.pre
        node.next.pre = node
        if node.pre:
            node.pre.next = node
        return node

    def remove(self) -> "Node[V]":
        """从链表里移除自身"""
        if self.pre:
            self.pre.next = self.next
        if self.next:
            self.next.pre = self.pre
        self.pre = None
        self.next = None
        return self

    def __repr__(self) -> str:
        return f"{self.value}->{self.next}"
