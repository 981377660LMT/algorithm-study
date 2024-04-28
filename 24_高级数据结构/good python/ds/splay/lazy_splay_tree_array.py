# from titan_pylib.data_structures.splay_tree.lazy_splay_tree_array import LazySplayTreeArray
from array import array
from typing import Generic, List, TypeVar, Tuple, Callable, Iterable, Optional, Union, Sequence
from __pypy__ import newlist_hint

T = TypeVar("T")
F = TypeVar("F")


class LazySplayTreeArrayData(Generic[T, F]):
    def __init__(
        self,
        op: Optional[Callable[[T, T], T]] = None,
        mapping: Optional[Callable[[F, T], T]] = None,
        composition: Optional[Callable[[F, F], F]] = None,
        e: T = None,
        id: F = None,
    ):
        self.op: Callable[[T, T], T] = (lambda s, t: e) if op is None else op
        self.mapping: Callable[[F, T], T] = (lambda f, s: e) if op is None else mapping
        self.composition: Callable[[F, F], F] = (lambda f, g: id) if op is None else composition
        self.e: T = e
        self.id: F = id
        self.keydata: List[T] = [e, e]
        self.lazy: List[F] = [id]
        self.arr: array[int] = array("I", bytes(16))
        # left:  arr[node<<2]
        # right: arr[node<<2|1]
        # size:  arr[node<<2|2]
        # rev:   arr[node<<2|3]
        self.end: int = 1

    def reserve(self, n: int) -> None:
        if n <= 0:
            return
        self.keydata += [self.e] * (2 * n)
        self.lazy += [self.id] * n
        self.arr += array("I", bytes(16 * n))


class LazySplayTreeArray(Generic[T, F]):
    def __init__(
        self,
        data: "LazySplayTreeArrayData[T, F]",
        n_or_a: Union[int, Iterable[T]] = 0,
        _root: int = 0,
    ) -> None:
        self.data = data
        self.root = _root
        if not n_or_a:
            return
        if isinstance(n_or_a, int):
            a = [data.e] * n_or_a
        elif not isinstance(n_or_a, Sequence):
            a = list(n_or_a)
        else:
            a = n_or_a
        if a:
            self._build(a)

    def _build(self, a: Sequence[T]) -> None:
        def rec(l: int, r: int) -> int:
            mid = (l + r) >> 1
            if l != mid:
                arr[mid << 2] = rec(l, mid)
            if mid + 1 != r:
                arr[mid << 2 | 1] = rec(mid + 1, r)
            self._update(mid)
            return mid

        n = len(a)
        keydata, arr = self.data.keydata, self.data.arr
        end = self.data.end
        self.data.reserve(n + end - len(keydata) // 2 + 1)
        self.data.end += n
        for i, e in enumerate(a):
            keydata[end + i << 1] = e
            keydata[end + i << 1 | 1] = e
        self.root = rec(end, n + end)

    def _make_node(self, key: T) -> int:
        data = self.data
        if data.end >= len(data.arr) // 4:
            data.keydata.append(key)
            data.keydata.append(key)
            data.lazy.append(data.id)
            data.arr.append(0)
            data.arr.append(0)
            data.arr.append(1)
            data.arr.append(0)
        else:
            data.keydata[data.end << 1] = key
            data.keydata[data.end << 1 | 1] = key
        data.end += 1
        return data.end - 1

    def _propagate(self, node: int) -> None:
        data = self.data
        arr = data.arr
        if arr[node << 2 | 3]:
            arr[node << 2], arr[node << 2 | 1] = arr[node << 2 | 1], arr[node << 2]
            arr[node << 2 | 3] = 0
            arr[arr[node << 2] << 2 | 3] ^= 1
            arr[arr[node << 2 | 1] << 2 | 3] ^= 1
        nlazy = data.lazy[node]
        if nlazy == data.id:
            return
        lnode, rnode = arr[node << 2], arr[node << 2 | 1]
        keydata, lazy = data.keydata, data.lazy
        lazy[node] = data.id
        if lnode:
            lazy[lnode] = data.composition(nlazy, lazy[lnode])
            lnode <<= 1
            keydata[lnode] = data.mapping(nlazy, keydata[lnode])
            keydata[lnode | 1] = data.mapping(nlazy, keydata[lnode | 1])
        if rnode:
            lazy[rnode] = data.composition(nlazy, lazy[rnode])
            rnode <<= 1
            keydata[rnode] = data.mapping(nlazy, keydata[rnode])
            keydata[rnode | 1] = data.mapping(nlazy, keydata[rnode | 1])

    def _update_triple(self, x: int, y: int, z: int) -> None:
        data = self.data
        keydata, arr = data.keydata, data.arr
        lx, rx = arr[x << 2], arr[x << 2 | 1]
        ly, ry = arr[y << 2], arr[y << 2 | 1]
        arr[z << 2 | 2] = arr[x << 2 | 2]
        arr[x << 2 | 2] = 1 + arr[lx << 2 | 2] + arr[rx << 2 | 2]
        arr[y << 2 | 2] = 1 + arr[ly << 2 | 2] + arr[ry << 2 | 2]
        keydata[z << 1 | 1] = keydata[x << 1 | 1]
        keydata[x << 1 | 1] = data.op(
            data.op(keydata[lx << 1 | 1], keydata[x << 1]), keydata[rx << 1 | 1]
        )
        keydata[y << 1 | 1] = data.op(
            data.op(keydata[ly << 1 | 1], keydata[y << 1]), keydata[ry << 1 | 1]
        )

    def _update_double(self, x: int, y: int) -> None:
        data = self.data
        keydata, arr = data.keydata, data.arr
        lx, rx = arr[x << 2], arr[x << 2 | 1]
        arr[y << 2 | 2] = arr[x << 2 | 2]
        arr[x << 2 | 2] = 1 + arr[lx << 2 | 2] + arr[rx << 2 | 2]
        keydata[y << 1 | 1] = keydata[x << 1 | 1]
        keydata[x << 1 | 1] = data.op(
            data.op(keydata[lx << 1 | 1], keydata[x << 1]), keydata[rx << 1 | 1]
        )

    def _update(self, node: int) -> None:
        data = self.data
        keydata, arr = data.keydata, data.arr
        lnode, rnode = arr[node << 2], arr[node << 2 | 1]
        arr[node << 2 | 2] = 1 + arr[lnode << 2 | 2] + arr[rnode << 2 | 2]
        keydata[node << 1 | 1] = data.op(
            data.op(keydata[lnode << 1 | 1], keydata[node << 1]), keydata[rnode << 1 | 1]
        )

    def _splay(self, path: List[int], d: int) -> None:
        arr = self.data.arr
        g = d & 1
        while len(path) > 1:
            pnode = path.pop()
            gnode = path.pop()
            f = d >> 1 & 1
            node = arr[pnode << 2 | g ^ 1]
            nnode = (pnode if g == f else node) << 2 | f
            arr[pnode << 2 | g ^ 1] = arr[node << 2 | g]
            arr[node << 2 | g] = pnode
            arr[gnode << 2 | f ^ 1] = arr[nnode]
            arr[nnode] = gnode
            self._update_triple(gnode, pnode, node)
            if not path:
                return
            d >>= 2
            g = d & 1
            arr[path[-1] << 2 | g ^ 1] = node
        pnode = path.pop()
        node = arr[pnode << 2 | g ^ 1]
        arr[pnode << 2 | g ^ 1] = arr[node << 2 | g]
        arr[node << 2 | g] = pnode
        self._update_double(pnode, node)

    def _kth_elm_splay(self, node: int, k: int) -> int:
        arr = self.data.arr
        if k < 0:
            k += arr[node << 2 | 2]
        d = 0
        path = []
        while True:
            self._propagate(node)
            t = arr[arr[node << 2] << 2 | 2]
            if t == k:
                if path:
                    self._splay(path, d)
                return node
            d = d << 1 | (t > k)
            path.append(node)
            node = arr[node << 2 | (t < k)]
            if t < k:
                k -= t + 1

    def _left_splay(self, node: int) -> int:
        if not node:
            return 0
        self._propagate(node)
        arr = self.data.arr
        if not arr[node << 2]:
            return node
        path = []
        while arr[node << 2]:
            path.append(node)
            node = arr[node << 2]
            self._propagate(node)
        self._splay(path, (1 << len(path)) - 1)
        return node

    def _right_splay(self, node: int) -> int:
        if not node:
            return 0
        self._propagate(node)
        arr = self.data.arr
        if not arr[node << 2 | 1]:
            return node
        path = []
        while arr[node << 2 | 1]:
            path.append(node)
            node = arr[node << 2 | 1]
            self._propagate(node)
        self._splay(path, 0)
        return node

    def reserve(self, n: int) -> None:
        self.data.reserve(n)

    def merge(self, other: "LazySplayTreeArray[T, F]") -> None:
        assert self.data is other.data
        if not other.root:
            return
        if not self.root:
            self.root = other.root
            return
        self.root = self._right_splay(self.root)
        self.data.arr[self.root << 2 | 1] = other.root
        self._update(self.root)

    def split(self, k: int) -> Tuple["LazySplayTreeArray[T, F]", "LazySplayTreeArray[T, F]"]:
        assert (
            -len(self) < k <= len(self)
        ), f"IndexError: LazySplayTreeArray.split({k}), len={len(self)}"
        if k < 0:
            k += len(self)
        if k >= self.data.arr[self.root << 2 | 2]:
            return self, LazySplayTreeArray(self.data, _root=0)
        self.root = self._kth_elm_splay(self.root, k)
        left = LazySplayTreeArray(self.data, _root=self.data.arr[self.root << 2])
        self.data.arr[self.root << 2] = 0
        self._update(self.root)
        return left, self

    def _internal_split(self, k: int) -> Tuple[int, int]:
        if k >= self.data.arr[self.root << 2 | 2]:
            return self.root, 0
        self.root = self._kth_elm_splay(self.root, k)
        left = self.data.arr[self.root << 2]
        self.data.arr[self.root << 2] = 0
        self._update(self.root)
        return left, self.root

    def reverse(self, l: int, r: int) -> None:
        assert (
            0 <= l <= r <= len(self)
        ), f"IndexError: LazySplayTreeArray.reverse({l}, {r}), len={len(self)}"
        if l == r:
            return
        data = self.data
        left, right = self._internal_split(r)
        if l:
            left = self._kth_elm_splay(left, l - 1)
        data.arr[(data.arr[left << 2 | 1] if l else left) << 2 | 3] ^= 1
        if right:
            data.arr[right << 2] = left
            self._update(right)
        self.root = right if right else left

    def all_reverse(self) -> None:
        self.data.arr[self.root << 2 | 3] ^= 1

    def apply(self, l: int, r: int, f: F) -> None:
        assert (
            0 <= l <= r <= len(self)
        ), f"IndexError: LazySplayTreeArray.apply({l}, {r}), len={len(self)}"
        data = self.data
        left, right = self._internal_split(r)
        keydata, lazy = data.keydata, data.lazy
        if l:
            left = self._kth_elm_splay(left, l - 1)
        node = data.arr[left << 2 | 1] if l else left
        keydata[node << 1] = data.mapping(f, keydata[node << 1])
        keydata[node << 1 | 1] = data.mapping(f, keydata[node << 1 | 1])
        lazy[node] = data.composition(f, lazy[node])
        if l:
            self._update(left)
        if right:
            data.arr[right << 2] = left
            self._update(right)
        self.root = right if right else left

    def all_apply(self, f: F) -> None:
        if not self.root:
            return
        data, node = self.data, self.root
        data.keydata[node << 1] = data.mapping(f, data.keydata[node << 1])
        data.keydata[node << 1 | 1] = data.mapping(f, data.keydata[node << 1 | 1])
        data.lazy[node] = data.composition(f, data.lazy[node])

    def prod(self, l: int, r: int) -> T:
        assert (
            0 <= l <= r <= len(self)
        ), f"IndexError: LazySplayTreeArray.prod({l}, {r}), len={len(self)}"
        data = self.data
        left, right = self._internal_split(r)
        if l:
            left = self._kth_elm_splay(left, l - 1)
        res = data.keydata[(data.arr[left << 2 | 1] if l else left) << 1 | 1]
        if right:
            data.arr[right << 2] = left
            self._update(right)
        self.root = right if right else left
        return res

    def all_prod(self) -> T:
        return self.data.keydata[self.root << 1 | 1]

    def insert(self, k: int, key: T) -> None:
        assert (
            -len(self) <= k <= len(self)
        ), f"IndexError: LazySplayTreeArray.insert({k}, {key}), len={len(self)}"
        if k < 0:
            k += len(self)
        data = self.data
        node = self._make_node(key)
        if not self.root:
            self._update(node)
            self.root = node
            return
        arr = data.arr
        if k == data.arr[self.root << 2 | 2]:
            arr[node << 2] = self._right_splay(self.root)
        else:
            node_ = self._kth_elm_splay(self.root, k)
            if arr[node_ << 2]:
                arr[node << 2] = arr[node_ << 2]
                arr[node_ << 2] = 0
                self._update(node_)
            arr[node << 2 | 1] = node_
        self._update(node)
        self.root = node

    def append(self, key: T) -> None:
        data = self.data
        node = self._right_splay(self.root)
        self.root = self._make_node(key)
        data.arr[self.root << 2] = node
        self._update(self.root)

    def appendleft(self, key: T) -> None:
        node = self._left_splay(self.root)
        self.root = self._make_node(key)
        self.data.arr[self.root << 2 | 1] = node
        self._update(self.root)

    def pop(self, k: int = -1) -> T:
        assert -len(self) <= k < len(self), f"IndexError: LazySplayTreeArray.pop({k})"
        data = self.data
        if k == -1:
            node = self._right_splay(self.root)
            self._propagate(node)
            self.root = data.arr[node << 2]
            return data.keydata[node << 1]
        self.root = self._kth_elm_splay(self.root, k)
        res = data.keydata[self.root << 1]
        if not data.arr[self.root << 2]:
            self.root = data.arr[self.root << 2 | 1]
        elif not data.arr[self.root << 2 | 1]:
            self.root = data.arr[self.root << 2]
        else:
            node = self._right_splay(data.arr[self.root << 2])
            data.arr[node << 2 | 1] = data.arr[self.root << 2 | 1]
            self.root = node
            self._update(self.root)
        return res

    def popleft(self) -> T:
        assert self, "IndexError: LazySplayTreeArray.popleft()"
        node = self._left_splay(self.root)
        self.root = self.data.arr[node << 2 | 1]
        return self.data.keydata[node << 1]

    def rotate(self, x: int) -> None:
        # 「末尾をを削除し先頭に挿入」をx回
        n = self.data.arr[self.root << 2 | 2]
        l, self = self.split(n - (x % n))
        self.merge(l)

    def tolist(self) -> List[T]:
        node = self.root
        arr, keydata = self.data.arr, self.data.keydata
        stack = newlist_hint(len(self))
        res = newlist_hint(len(self))
        while stack or node:
            if node:
                self._propagate(node)
                stack.append(node)
                node = arr[node << 2]
            else:
                node = stack.pop()
                res.append(keydata[node << 1])
                node = arr[node << 2 | 1]
        return res

    def clear(self) -> None:
        self.root = 0

    def __setitem__(self, k: int, key: T):
        assert -len(self) <= k < len(self), f"IndexError: LazySplayTreeArray.__setitem__({k})"
        self.root = self._kth_elm_splay(self.root, k)
        self.data.keydata[self.root << 1] = key
        self._update(self.root)

    def __getitem__(self, k: int) -> T:
        assert -len(self) <= k < len(self), f"IndexError: LazySplayTreeArray.__getitem__({k})"
        self.root = self._kth_elm_splay(self.root, k)
        return self.data.keydata[self.root << 1]

    def __iter__(self):
        self.__iter = 0
        return self

    def __next__(self):
        if self.__iter == self.data.arr[self.root << 2 | 2]:
            raise StopIteration
        res = self[self.__iter]
        self.__iter += 1
        return res

    def __reversed__(self):
        for i in range(len(self)):
            yield self[-i - 1]

    def __len__(self):
        return self.data.arr[self.root << 2 | 2]

    def __str__(self):
        return str(self.tolist())

    def __bool__(self):
        return self.root != 0

    def __repr__(self):
        return f"{self.__class__.__name__}({self})"
