# from titan_pylib.data_structures.heap.binomial_heap import BinomialHeap
from typing import Generic, Iterable, TypeVar, List
from itertools import chain

T = TypeVar("T")


class BinomialHeap(Generic[T]):
    """二項ヒープです。

    計算量はメチャクチャさぼってます。
    あらゆる操作が :math:`\\theta{(\\log{n})}` です。

    ``List`` の代わりに ``LinkedList`` を使用し、``push,meld`` では ``O(1)`` で連結させ、 ``delete_min`` にすべてを押し付けると ``push,meld`` が ``O(1)`` 、``delete_min`` が償却 ``O(logn)`` になるはずです。
    """

    class _Node:
        def __init__(self, key: T) -> None:
            self.key = key
            self.child: List["BinomialHeap._Node"] = []

        def rank(self) -> int:
            return len(self.child)

        def tolist(self):
            a = []

            def dfs(node):
                if not node:
                    return
                a.append(node.key)
                for c in node.child:
                    dfs(c)

            dfs(self)
            return a

        def __str__(self):
            return f"_Node({sorted(self.tolist())})"

        __repr__ = __str__

    def __init__(self, a: Iterable[T] = []) -> None:
        self.ptr = []
        self.ptr_min = None
        self.len = 0
        for e in a:
            self.insert(e)

    @staticmethod
    def _link(node1: _Node, node2: _Node) -> _Node:
        if node1.key > node2.key:
            node1, node2 = node2, node1
        node1.child.append(node2)
        return node1

    def _set_min_node(self) -> None:
        self.ptr_min = None
        for node in self.ptr:
            if not node:
                continue
            if (self.ptr_min is None) or (self.ptr_min.key > node.key):
                self.ptr_min = node

    # O(logN)
    def push(self, key: T) -> None:
        self.len += 1
        new_node = BinomialHeap._Node(key)
        ptr = self.ptr
        if self.ptr_min is None or self.ptr_min.key > key:
            self.ptr_min = new_node
        for rank, node in enumerate(self.ptr):
            if node is None:
                ptr[rank] = new_node
                return
            new_node = BinomialHeap._link(new_node, node)
            ptr[rank] = None
        ptr.append(new_node)

    # O(1)
    def get_min(self) -> T:
        assert self.ptr_min, "IndexError: find_min from empty BinomialHeap"
        return self.ptr_min.key

    # O(logN)
    def pop_min(self) -> T:
        _len = self.len - 1
        node = self.ptr_min
        new_bheap = BinomialHeap()
        new_bheap.ptr = node.child
        new_bheap._set_min_node()
        self.ptr[node.rank()] = None
        while self.ptr and self.ptr[-1] is None:
            self.ptr.pop()
        self._set_min_node()
        self.meld(new_bheap)
        self.len = _len

    # O(logN)
    def meld(self, other: "BinomialHeap") -> None:
        self.len += other.len
        h0, h1 = self.ptr, other.ptr
        if len(h0) > len(h1):
            h0, h1 = h1, h0
        n = len(h1)
        new_node = None
        min_node = self.ptr_min
        if (other.ptr_min) and (min_node is None or min_node.key > other.ptr_min.key):
            min_node = other.ptr_min
        for rank in range(n):
            if rank >= len(h0) and new_node is None:
                break
            cnt = (rank < len(h0) and h0[rank] is not None) + (h1[rank] is not None)
            if new_node:
                if cnt == 2:
                    x = new_node
                    new_node = BinomialHeap._link(h0[rank], h1[rank])
                    h1[rank] = x
                elif cnt == 1:
                    new_node = BinomialHeap._link(new_node, (h1[rank] if h1[rank] else h0[rank]))
                    h1[rank] = None
                else:
                    h1[rank] = new_node
                    new_node = None
            else:
                if cnt == 2:
                    new_node = BinomialHeap._link(h0[rank], h1[rank])
                    h1[rank] = None
                if cnt == 1:
                    if h1[rank] is None:
                        h1[rank] = h0[rank]
        if new_node:
            h1.append(new_node)
        self.ptr = h1
        self.ptr_min = min_node

    def tolist(self) -> List[T]:
        return sorted(chain(*[node.tolist() for node in self.ptr if node]))

    def __len__(self):
        return self.len

    def __str__(self):
        return str(self.tolist())

    def __repr__(self):
        return f"BinomialHeap({self})"
