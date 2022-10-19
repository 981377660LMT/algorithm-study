# https://www.luogu.com.cn/problem/solution/P3835
# FHQ:范浩强
# !FHQ Treap(无旋Treap,能代替普通treap和splay的所有功能) 实现可持久化平衡树
# 时空复杂度O(qlogq)


# TODO
from typing import Optional


INF = int(1e18)


class PersistentFHQTreap:
    """可持久化的Multiset(多重集合)"""

    __slots__ = ""

    def __init__(self) -> None:
        ...

    def add(self, version: int, value: int) -> int:
        ...

    def discard(self, version: int, value: int) -> int:
        ...

    def queryRank(self, version: int, value: int) -> int:
        ...

    def queryKth(self, version: int, k: int) -> int:
        ...

    def queryPrecursor(self, version: int, value: int) -> Optional[int]:
        ...

    def querySuccessor(self, version: int, value: int) -> Optional[int]:
        ...

    def __len__(self) -> int:
        ...
