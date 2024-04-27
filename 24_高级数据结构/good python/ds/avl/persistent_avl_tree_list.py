# from titan_pylib.data_structures.avl_tree.persistent_avl_tree_list import PersistentAVLTreeList
from typing import Generic, Iterable, Optional, TypeVar, List, Tuple

T = TypeVar("T")


class PersistentAVLTreeList(Generic[T]):
    """挿入削除が対数時間で行える永続AVL木です。"""

    class Node:
        def __init__(self, key: T):
            self.key: T = key
            self.left: Optional[PersistentAVLTreeList.Node] = None
            self.right: Optional[PersistentAVLTreeList.Node] = None
            self.height: int = 1
            self.size: int = 1

        def copy(self) -> "PersistentAVLTreeList.Node":
            node = PersistentAVLTreeList.Node(self.key)
            node.left = self.left
            node.right = self.right
            node.height = self.height
            node.size = self.size
            return node

        def balance(self) -> int:
            return (
                (0 if self.right is None else -self.right.height)
                if self.left is None
                else (
                    self.left.height if self.right is None else self.left.height - self.right.height
                )
            )

        def __str__(self):
            if self.left is None and self.right is None:
                return f"key={self.key}, height={self.height}, size={self.size}\n"
            return f"key={self.key}, height={self.height}, size={self.size},\n left:{self.left},\n right:{self.right}\n"

        __repr__ = __str__

    def __init__(self, a: Iterable[T] = [], _root: Optional[Node] = None) -> None:
        self.root: Optional[PersistentAVLTreeList.Node] = _root
        a = list(a)
        if a:
            self._build(list(a))

    def _build(self, a: List[T]) -> None:
        Node = PersistentAVLTreeList.Node

        def build(l: int, r: int) -> Node:
            mid = (l + r) >> 1
            node = Node(a[mid])
            if l != mid:
                node.left = build(l, mid)
            if mid + 1 != r:
                node.right = build(mid + 1, r)
            self._update(node)
            return node

        self.root = build(0, len(a))

    def _update(self, node: Node) -> None:
        if node.left:
            if node.right:
                node.size = 1 + node.left.size + node.right.size
                node.height = (
                    node.left.height + 1
                    if node.left.height > node.right.height
                    else node.right.height + 1
                )
            else:
                node.size = 1 + node.left.size
                node.height = node.left.height + 1
        else:
            if node.right:
                node.size = 1 + node.right.size
                node.height = node.right.height + 1
            else:
                node.size = 1
                node.height = 1

    def _rotate_right(self, node: Node) -> Node:
        assert node.left
        u = node.left.copy()
        node.left = u.right
        u.right = node
        self._update(node)
        self._update(u)
        return u

    def _rotate_left(self, node: Node) -> Node:
        assert node.right
        u = node.right.copy()
        node.right = u.left
        u.left = node
        self._update(node)
        self._update(u)
        return u

    def _balance_left(self, node: Node) -> Node:
        assert node.right
        node.right = node.right.copy()
        u = node.right
        if u.balance() == 1:
            assert u.left
            node.right = self._rotate_right(u)
        u = self._rotate_left(node)
        return u

    def _balance_right(self, node: Node) -> Node:
        assert node.left
        node.left = node.left.copy()
        u = node.left
        if u.balance() == -1:
            assert u.right
            node.left = self._rotate_left(u)
        u = self._rotate_right(node)
        return u

    def _merge_with_root(self, l: Optional[Node], root: Node, r: Optional[Node]) -> Node:
        diff = 0
        if l:
            diff += l.height
        if r:
            diff -= r.height
        if diff > 1:
            assert l
            l = l.copy()
            l.right = self._merge_with_root(l.right, root, r)
            self._update(l)
            if l.balance() == -2:
                return self._balance_left(l)
            return l
        if diff < -1:
            assert r
            r = r.copy()
            r.left = self._merge_with_root(l, root, r.left)
            self._update(r)
            if r.balance() == 2:
                return self._balance_right(r)
            return r
        root = root.copy()
        root.left = l
        root.right = r
        self._update(root)
        return root

    def _merge_node(self, l: Optional[Node], r: Optional[Node]) -> Optional[Node]:
        if l is None and r is None:
            return None
        if l is None:
            assert r
            return r.copy()
        if r is None:
            return l.copy()
        l = l.copy()
        r = r.copy()
        l, root = self._pop_right(l)
        return self._merge_with_root(l, root, r)

    def merge(self, other: "PersistentAVLTreeList") -> "PersistentAVLTreeList":
        root = self._merge_node(self.root, other.root)
        return self._new(root)

    def _pop_right(self, node: Node) -> Tuple[Node, Node]:
        path = []
        node = node.copy()
        mx = node
        while node.right is not None:
            path.append(node)
            node = node.right.copy()
            mx = node
        path.append(node.left.copy() if node.left else None)
        for _ in range(len(path) - 1):
            node = path.pop()
            if node is None:
                path[-1].right = None
                self._update(path[-1])
                continue
            b = node.balance()
            if b == 2:
                path[-1].right = self._balance_right(node)
            elif b == -2:
                path[-1].right = self._balance_left(node)
            else:
                path[-1].right = node
            self._update(path[-1])
        if path[0] is not None:
            b = path[0].balance()
            if b == 2:
                path[0] = self._balance_right(path[0])
            elif b == -2:
                path[0] = self._balance_left(path[0])
        mx.left = None
        self._update(mx)
        return path[0], mx

    def _split_node(self, node: Optional[Node], k: int) -> Tuple[Optional[Node], Optional[Node]]:
        if node is None:
            return None, None
        tmp = k if node.left is None else k - node.left.size
        l, r = None, None
        if tmp == 0:
            return node.left, self._merge_with_root(None, node, node.right)
        elif tmp < 0:
            l, r = self._split_node(node.left, k)
            return l, self._merge_with_root(r, node, node.right)
        else:
            l, r = self._split_node(node.right, tmp - 1)
            return self._merge_with_root(node.left, node, l), r

    def split(self, k: int) -> Tuple["PersistentAVLTreeList", "PersistentAVLTreeList"]:
        l, r = self._split_node(self.root, k)
        return self._new(l), self._new(r)

    def _new(self, root: Optional["PersistentAVLTreeList.Node"]) -> "PersistentAVLTreeList":
        return PersistentAVLTreeList([], root)

    def insert(self, k: int, key: T) -> "PersistentAVLTreeList":
        s, t = self._split_node(self.root, k)
        root = self._merge_with_root(s, PersistentAVLTreeList.Node(key), t)
        return self._new(root)

    def pop(self, k: int) -> Tuple["PersistentAVLTreeList", T]:
        s, t = self._split_node(self.root, k + 1)
        assert s
        s, tmp = self._pop_right(s)
        root = self._merge_node(s, t)
        return self._new(root), tmp.key

    def tolist(self) -> List[T]:
        node = self.root
        stack = []
        a = []
        while stack or node:
            if node:
                stack.append(node)
                node = node.left
            else:
                node = stack.pop()
                a.append(node.key)
                node = node.right
        return a

    def __getitem__(self, k: int) -> T:
        if k < 0:
            k += len(self)
        node = self.root
        while True:
            assert node
            t = 0 if node.left is None else node.left.size
            if t == k:
                return node.key
            elif t < k:
                k -= t + 1
                node = node.right
            else:
                node = node.left

    def __len__(self):
        return 0 if self.root is None else self.root.size

    def __str__(self):
        return "[" + ", ".join(map(str, self.tolist())) + "]"

    def __repr__(self):
        return f"PersistentAVLTreeList({self})"
