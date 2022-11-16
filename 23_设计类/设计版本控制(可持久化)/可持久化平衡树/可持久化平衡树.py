# 大根堆无旋Treap
# https://www.luogu.com.cn/problem/P3835
# https://www.luogu.com.cn/record/41672571

# golang实现
# !js-algorithm\算法竞赛进阶指南\GoDS (Go Data Structures)\src\tree\fhqtreap\practice\持久化平衡树


import sys
from random import random
from typing import Optional, Tuple


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")


class Node:
    __slots__ = "value", "priority", "size", "left", "right"

    def __init__(self, value: int, priority=None):
        self.value = value
        self.priority = priority if priority is not None else random()
        self.size = 1
        self.left = None
        self.right = None


class PersistentFHQTreap:
    @classmethod
    def _merge(cls, left: Optional["Node"], right: Optional["Node"]) -> Optional["Node"]:
        if not left or not right:
            return left if left else right
        size = left.size + right.size
        if left.priority < right.priority:
            newNode = Node(left.value, left.priority)
            newNode.left = left.left
            newNode.right = cls._merge(left.right, right)  # type: ignore
            newNode.size = size  # pushUp
            return newNode
        else:
            newNode = Node(right.value, right.priority)
            newNode.right = right.right
            newNode.left = cls._merge(left, right.left)  # type: ignore
            newNode.size = size
            return newNode

    @classmethod
    def _split(cls, root: "Node", value: int) -> Tuple[Optional["Node"], Optional["Node"]]:
        if value < root.value:
            if root.left:
                newNode = Node(root.value, root.priority)
                newNode.right = root.right
                left, right = cls._split(root.left, value)
                newNode.left = right  # type: ignore
                newNode.size = root.size - (left.size if left else 0)
                return left, newNode
            return None, root
        else:
            if root.right:
                newNode = Node(root.value, root.priority)
                newNode.left = root.left
                left, right = cls._split(root.right, value)
                newNode.right = left  # type: ignore
                newNode.size = root.size - (right.size if right else 0)
                return newNode, right
            return root, None

    __slots__ = "roots"

    def __init__(self):
        self.roots = []  # 记录历史版本的根节点

    def add(self, version: int, value: int) -> None:
        node = Node(value)
        if not (0 <= version < len(self.roots)):
            self.roots.append(node)
            return
        root = self.roots[version]
        if not root:
            self.roots.append(node)
            return
        left, right = self._split(root, value)
        self.roots.append(self._merge(self._merge(left, node), right))

    def discard(self, version: int, value: int) -> bool:
        if not (0 <= version < len(self.roots)):
            self.roots.append(None)
            return False
        mid, right = self._split(self.roots[version], value)
        if not mid:
            self.roots.append(self.roots[version])
            return False
        left, mid = self._split(mid, value - 1)
        if not mid:
            self.roots.append(self.roots[version])
            return False
        mid = self._merge(mid.left, mid.right)
        self.roots.append(self._merge(self._merge(left, mid), right))
        return True

    def at(self, version: int, index: int) -> Optional[int]:
        """1<=index<=size"""
        if not (0 <= version < len(self.roots)):
            self.roots.append(None)
            return None
        node = self.roots[version]
        self.roots.append(node)
        if index > node.size or index < 1:
            return None
        while True:
            left_size = node.left.size if node.left else 0
            if index == left_size + 1:
                return node.value
            elif index <= left_size:
                node = node.left
            else:
                index -= 1 + left_size
                node = node.right

    def bisect_left(self, version: int, value: int) -> int:
        if not (0 <= version < len(self.roots)):
            self.roots.append(None)
            return 1
        root = self.roots[version]
        self.roots.append(root)
        left, _ = self._split(root, value - 1)
        pos = left.size if left else 0
        return pos

    def lower(self, version: int, value: int) -> Optional[int]:
        if not (0 <= version < len(self.roots)):
            self.roots.append(None)
            return None
        root = self.roots[version]
        self.roots.append(root)
        left, _ = self._split(self.roots[version], value - 1)
        if not left:
            return None
        x = left
        while x.right:
            x = x.right
        return x.value

    def upper(self, version: int, value: int) -> Optional[int]:
        if not (0 <= version < len(self.roots)):
            self.roots.append(None)
            return None
        root = self.roots[version]
        self.roots.append(root)
        _, right = self._split(self.roots[version], value)
        if not right:
            return None
        x = right
        while x.left:
            x = x.left
        return x.value


if __name__ == "__main__":

    tree = PersistentFHQTreap()
    n = int(input())
    for _ in range(n):
        version, op, num = map(int, input().split())
        version -= 1  # !注意开始版本为-1 每次操作(不管查询还是更新都会产生新的版本)
        if op == 1:
            tree.add(version, num)
        elif op == 2:
            tree.discard(version, num)
        elif op == 3:
            print(tree.bisect_left(version, num) + 1)
        elif op == 4:
            print(tree.at(version, num))
        elif op == 5:
            res = tree.lower(version, num)
            print(res if res is not None else -(2**31) + 1)
        elif op == 6:
            res = tree.upper(version, num)
            print(res if res is not None else 2**31 - 1)
