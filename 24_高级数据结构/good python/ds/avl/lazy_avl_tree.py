# from titan_pylib.data_structures.avl_tree.lazy_avl_tree import LazyAVLTree
from typing import Generic, Iterable, TypeVar, Callable, List, Tuple, Optional

T = TypeVar("T")
F = TypeVar("F")


class LazyAVLTree(Generic[T, F]):
    """遅延伝播反転可能平衡二分木です。"""

    class Node:
        def __init__(self, key: T, id: F):
            self.key: T = key
            self.data: T = key
            self.left: Optional[LazyAVLTree.Node] = None
            self.right: Optional[LazyAVLTree.Node] = None
            self.lazy: F = id
            self.rev: int = 0
            self.height: int = 1
            self.size: int = 1

        def __str__(self):
            if self.left is None and self.right is None:
                return f"key:{self.key, self.height, self.size, self.data, self.lazy, self.rev}\n"
            return f"key:{self.key, self.height, self.size, self.data, self.lazy, self.rev},\n left:{self.left},\n right:{self.right}\n"

    def __init__(
        self,
        a: Iterable[T],
        op: Callable[[T, T], T],
        mapping: Callable[[F, T], T],
        composition: Callable[[F, F], F],
        e: T,
        id: F,
        node: Node = None,
    ) -> None:
        self.root: Optional[LazyAVLTree.Node] = node
        self.op: Callable[[T, T], T] = op
        self.mapping: Callable[[F, T], T] = mapping
        self.composition: Callable[[F, F], F] = composition
        self.e: T = e
        self.id: F = id
        a = list(a)
        if a:
            self._build(a)

    def _build(self, a: List[T]) -> None:
        Node = LazyAVLTree.Node
        id = self.id

        def sort(l: int, r: int) -> Node:
            mid = (l + r) >> 1
            node = Node(a[mid], id)
            if l != mid:
                node.left = sort(l, mid)
            if mid + 1 != r:
                node.right = sort(mid + 1, r)
            self._update(node)
            return node

        self.root = sort(0, len(a))

    def _propagate(self, node: Node) -> None:
        l, r = node.left, node.right
        if node.rev:
            node.left, node.right = r, l
            if l:
                l.rev ^= 1
            if r:
                r.rev ^= 1
            node.rev = 0
        if node.lazy != self.id:
            lazy = node.lazy
            if l:
                l.data = self.mapping(lazy, l.data)
                l.key = self.mapping(lazy, l.key)
                l.lazy = lazy if l.lazy == self.id else self.composition(lazy, l.lazy)
            if r:
                r.data = self.mapping(lazy, r.data)
                r.key = self.mapping(lazy, r.key)
                r.lazy = lazy if r.lazy == self.id else self.composition(lazy, r.lazy)
            node.lazy = self.id

    def _update(self, node: Node) -> None:
        node.size = 1
        node.data = node.key
        node.height = 1
        if node.left:
            node.size += node.left.size
            node.data = self.op(node.left.data, node.data)
            node.height = max(node.left.height + 1, 1)
        if node.right:
            node.size += node.right.size
            node.data = self.op(node.data, node.right.data)
            node.height = max(node.height, node.right.height + 1)

    def _get_balance(self, node: Node) -> int:
        return (
            (0 if node.right is None else -node.right.height)
            if node.left is None
            else (node.left.height if node.right is None else node.left.height - node.right.height)
        )

    def _balance_left(self, node: Node) -> Node:
        self._propagate(node.left)
        if node.left.left is None or node.left.left.height + 2 == node.left.height:
            u = node.left.right
            self._propagate(u)
            node.left.right = u.left
            u.left = node.left
            node.left = u.right
            u.right = node
            self._update(u.left)
        else:
            u = node.left
            node.left = u.right
            u.right = node
        self._update(u.right)
        self._update(u)
        return u

    def _balance_right(self, node: Node) -> Node:
        self._propagate(node.right)
        if node.right.right is None or node.right.right.height + 2 == node.right.height:
            u = node.right.left
            self._propagate(u)
            node.right.left = u.right
            u.right = node.right
            node.right = u.left
            u.left = node
            self._update(u.right)
        else:
            u = node.right
            node.right = u.left
            u.left = node
        self._update(u.left)
        self._update(u)
        return u

    def _kth_elm(self, k: int) -> T:
        if k < 0:
            k += len(self)
        node = self.root
        while True:
            self._propagate(node)
            t = 0 if node.left is None else node.left.size
            if t == k:
                return node.key
            elif t < k:
                k -= t + 1
                node = node.right
            else:
                node = node.left

    def _merge_with_root(self, l: Node, root: Node, r: Node) -> Node:
        diff = (
            (0 if r is None else -r.height)
            if l is None
            else (l.height if r is None else l.height - r.height)
        )
        if diff > 1:
            self._propagate(l)
            l.right = self._merge_with_root(l.right, root, r)
            self._update(l)
            if -l.right.height if l.left is None else l.left.height - l.right.height == -2:
                return self._balance_right(l)
            return l
        elif diff < -1:
            self._propagate(r)
            r.left = self._merge_with_root(l, root, r.left)
            self._update(r)
            if r.left.height if r.right is None else r.left.height - r.right.height == 2:
                return self._balance_left(r)
            return r
        else:
            root.left = l
            root.right = r
            self._update(root)
            return root

    def _merge_node(self, l: Node, r: Node) -> Node:
        if l is None:
            return r
        if r is None:
            return l
        l, tmp = self._pop_max(l)
        return self._merge_with_root(l, tmp, r)

    def merge(self, other: "LazyAVLTree") -> None:
        self.root = self._merge_node(self.root, other.node)

    def _pop_max(self, node: Node) -> Tuple[Node, Node]:
        self._propagate(node)
        path = []
        mx = node
        while node.right:
            path.append(node)
            mx = node.right
            node = node.right
            self._propagate(node)
        path.append(node.left)
        for _ in range(len(path) - 1):
            node = path.pop()
            if node is None:
                path[-1].right = None
                self._update(path[-1])
                continue
            b = self._get_balance(node)
            path[-1].right = (
                self._balance_left(node)
                if b == 2
                else self._balance_right(node)
                if b == -2
                else node
            )
            self._update(path[-1])
        if path[0]:
            b = self._get_balance(path[0])
            path[0] = (
                self._balance_left(path[0])
                if b == 2
                else self._balance_right(path[0])
                if b == -2
                else path[0]
            )
        mx.left = None
        self._update(mx)
        return path[0], mx

    def _split_node(self, node: Node, k: int) -> Tuple[Node, Node]:
        if not node:
            return None, None
        self._propagate(node)
        tmp = k if node.left is None else k - node.left.size
        if tmp == 0:
            return node.left, self._merge_with_root(None, node, node.right)
        elif tmp < 0:
            s, t = self._split_node(node.left, k)
            return s, self._merge_with_root(t, node, node.right)
        else:
            s, t = self._split_node(node.right, tmp - 1)
            return self._merge_with_root(node.left, node, s), t

    def split(self, k: int) -> Tuple["LazyAVLTree", "LazyAVLTree"]:
        l, r = self._split_node(self.root, k)
        return LazyAVLTree(
            [], self.op, self.mapping, self.composition, self.e, self.id, l
        ), LazyAVLTree([], self.op, self.mapping, self.composition, self.e, self.id, r)

    def insert(self, k: int, key: T) -> None:
        s, t = self._split_node(self.root, k)
        self.root = self._merge_with_root(s, LazyAVLTree.Node(key, self.id), t)

    def pop(self, k: int) -> T:
        s, t = self._split_node(self.root, k + 1)
        s, tmp = self._pop_max(s)
        self.root = self._merge_node(s, t)
        return tmp.key

    def apply(self, l: int, r: int, f: F) -> None:
        if l >= r or (not self.root):
            return
        stack = [(self.root), (self.root, 0, len(self))]
        while stack:
            if isinstance(stack[-1], tuple):
                node, left, right = stack.pop()
                if right <= l or r <= left:
                    continue
                self._propagate(node)
                if l <= left and right < r:
                    node.key = self.mapping(f, node.key)
                    node.data = self.mapping(f, node.data)
                    node.lazy = f if node.lazy == self.id else self.composition(f, node.lazy)
                else:
                    lsize = node.left.size if node.left else 0
                    stack.append(node)
                    if node.left:
                        stack.append((node.left, left, left + lsize))
                    if l <= left + lsize < r:
                        node.key = self.mapping(f, node.key)
                    if node.right:
                        stack.append((node.right, left + lsize + 1, right))
            else:
                self._update(stack.pop())

    def all_apply(self, f: F) -> None:
        if not self.root:
            return
        self.root.key = self.mapping(f, self.root.key)
        self.root.data = self.mapping(f, self.root.data)
        self.root.lazy = f if self.root.lazy == self.id else self.composition(f, self.root.lazy)

    def reverse(self, l: int, r: int) -> None:
        if l >= r:
            return
        s, t = self._split_node(self.root, r)
        r, s = self._split_node(s, l)
        s.rev ^= 1
        self.root = self._merge_node(self._merge_node(r, s), t)

    def all_reverse(self) -> None:
        if self.root is None:
            return
        self.root.rev ^= 1

    def prod(self, l: int, r: int) -> T:
        if l >= r or (not self.root):
            return self.e

        def dfs(node: LazyAVLTree.Node, left: int, right: int) -> T:
            if right <= l or r <= left:
                return self.e
            self._propagate(node)
            if l <= left and right < r:
                return node.data
            lsize = node.left.size if node.left else 0
            res = self.e
            if node.left:
                res = dfs(node.left, left, left + lsize)
            if l <= left + lsize < r:
                res = self.op(res, node.key)
            if node.right:
                res = self.op(res, dfs(node.right, left + lsize + 1, right))
            return res

        return dfs(self.root, 0, len(self))

    def all_prod(self) -> T:
        return self.root.data if self.root else self.e

    def clear(self) -> None:
        self.root = None

    def tolist(self) -> List[T]:
        node = self.root
        stack = []
        a = []
        while stack or node:
            if node:
                self._propagate(node)
                stack.append(node)
                node = node.left
            else:
                node = stack.pop()
                a.append(node.key)
                node = node.right
        return a

    def __len__(self):
        return 0 if self.root is None else self.root.size

    def __iter__(self):
        self.__iter = 0
        return self

    def __next__(self):
        if self.__iter == len(self):
            raise StopIteration
        res = self[self.__iter]
        self.__iter += 1
        return res

    def __reversed__(self):
        for i in range(len(self)):
            yield self[-i - 1]

    def __bool__(self):
        return self.root is not None

    def __getitem__(self, k: int) -> T:
        return self._kth_elm(k)

    def __setitem__(self, k, key: T):
        if k < 0:
            k += len(self)
        node = self.root
        path = []
        while True:
            self._propagate(node)
            path.append(node)
            t = 0 if node.left is None else node.left.size
            if t == k:
                node.key = key
                break
            if t < k:
                k -= t + 1
                node = node.right
            else:
                node = node.left
        while path:
            self._update(path.pop())

    def __str__(self):
        return "[" + ", ".join(map(str, self.tolist())) + "]"

    def __repr__(self):
        return f"LazyAVLTree({self})"
