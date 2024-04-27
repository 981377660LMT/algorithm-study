# from titan_pylib.data_structures.avl_tree.avl_tree_multiset import AVLTreeMultiset
# from titan_pylib.my_class.ordered_multiset_interface import OrderedMultisetInterface
# from titan_pylib.my_class.supports_less_than import SupportsLessThan


# 多重集合としての AVL 木です。 配列を用いてノードを表現しています。 size を持ちます。
# 仅用于保存

from typing import Protocol


class SupportsLessThan(Protocol):
    def __lt__(self, other) -> bool:
        ...


from abc import ABC, abstractmethod
from typing import Iterable, Optional, Iterator, TypeVar, Generic, List

T = TypeVar("T", bound=SupportsLessThan)


class OrderedMultisetInterface(ABC, Generic[T]):
    @abstractmethod
    def __init__(self, a: Iterable[T]) -> None:
        raise NotImplementedError

    @abstractmethod
    def add(self, key: T, cnt: int) -> None:
        raise NotImplementedError

    @abstractmethod
    def discard(self, key: T, cnt: int) -> bool:
        raise NotImplementedError

    @abstractmethod
    def discard_all(self, key: T) -> bool:
        raise NotImplementedError

    @abstractmethod
    def count(self, key: T) -> int:
        raise NotImplementedError

    @abstractmethod
    def remove(self, key: T, cnt: int) -> None:
        raise NotImplementedError

    @abstractmethod
    def le(self, key: T) -> Optional[T]:
        raise NotImplementedError

    @abstractmethod
    def lt(self, key: T) -> Optional[T]:
        raise NotImplementedError

    @abstractmethod
    def ge(self, key: T) -> Optional[T]:
        raise NotImplementedError

    @abstractmethod
    def gt(self, key: T) -> Optional[T]:
        raise NotImplementedError

    @abstractmethod
    def get_max(self) -> Optional[T]:
        raise NotImplementedError

    @abstractmethod
    def get_min(self) -> Optional[T]:
        raise NotImplementedError

    @abstractmethod
    def pop_max(self) -> T:
        raise NotImplementedError

    @abstractmethod
    def pop_min(self) -> T:
        raise NotImplementedError

    @abstractmethod
    def clear(self) -> None:
        raise NotImplementedError

    @abstractmethod
    def tolist(self) -> List[T]:
        raise NotImplementedError

    @abstractmethod
    def __iter__(self) -> Iterator:
        raise NotImplementedError

    @abstractmethod
    def __next__(self) -> T:
        raise NotImplementedError

    @abstractmethod
    def __contains__(self, key: T) -> bool:
        raise NotImplementedError

    @abstractmethod
    def __len__(self) -> int:
        raise NotImplementedError

    @abstractmethod
    def __bool__(self) -> bool:
        raise NotImplementedError

    @abstractmethod
    def __str__(self) -> str:
        raise NotImplementedError

    @abstractmethod
    def __repr__(self) -> str:
        raise NotImplementedError


# from titan_pylib.data_structures.bst_base.bst_multiset_array_base import BSTMultisetArrayBase
from __pypy__ import newlist_hint
from typing import List, Tuple, TypeVar, Generic, Optional

T = TypeVar("T")
BST = TypeVar("BST")
# protcolで、key,val,left,right を規定


class BSTMultisetArrayBase(Generic[BST, T]):
    @staticmethod
    def _rle(a: List[T]) -> Tuple[List[T], List[int]]:
        keys, vals = [a[0]], [1]
        for i, elm in enumerate(a):
            if i == 0:
                continue
            if elm == keys[-1]:
                vals[-1] += 1
                continue
            keys.append(elm)
            vals.append(1)
        return keys, vals

    @staticmethod
    def count(bst: BST, key: T) -> int:
        keys, left, right = bst.key, bst.left, bst.right
        node = bst.root
        while node:
            if keys[node] == key:
                return bst.val[node]
            node = left[node] if key < keys[node] else right[node]
        return 0

    @staticmethod
    def le(bst: BST, key: T) -> Optional[T]:
        keys, left, right = bst.key, bst.left, bst.right
        res = None
        node = bst.root
        while node:
            if key == keys[node]:
                res = key
                break
            if key < keys[node]:
                node = left[node]
            else:
                res = keys[node]
                node = right[node]
        return res

    @staticmethod
    def lt(bst: BST, key: T) -> Optional[T]:
        keys, left, right = bst.key, bst.left, bst.right
        res = None
        node = bst.root
        while node:
            if key <= keys[node]:
                node = left[node]
            else:
                res = keys[node]
                node = right[node]
        return res

    @staticmethod
    def ge(bst: BST, key: T) -> Optional[T]:
        keys, left, right = bst.key, bst.left, bst.right
        res = None
        node = bst.root
        while node:
            if key == keys[node]:
                res = key
                break
            if key < keys[node]:
                res = keys[node]
                node = left[node]
            else:
                node = right[node]
        return res

    @staticmethod
    def gt(bst: BST, key: T) -> Optional[T]:
        keys, left, right = bst.key, bst.left, bst.right
        res = None
        node = bst.root
        while node:
            if key < keys[node]:
                res = keys[node]
                node = left[node]
            else:
                node = right[node]
        return res

    @staticmethod
    def index(bst: BST, key: T) -> int:
        keys, left, right, vals, valsize = bst.key, bst.left, bst.right, bst.val, bst.valsize
        k = 0
        node = bst.root
        while node:
            if key == keys[node]:
                if left[node]:
                    k += valsize[left[node]]
                break
            if key < keys[node]:
                node = left[node]
            else:
                k += valsize[left[node]] + vals[node]
                node = right[node]
        return k

    @staticmethod
    def index_right(bst: BST, key: T) -> int:
        keys, left, right, vals, valsize = bst.key, bst.left, bst.right, bst.val, bst.valsize
        k = 0
        node = bst.root
        while node:
            if key == keys[node]:
                k += valsize[left[node]] + vals[node]
                break
            if key < keys[node]:
                node = left[node]
            else:
                k += valsize[left[node]] + vals[node]
                node = right[node]
        return k

    @staticmethod
    def _kth_elm(bst: BST, k: int) -> Tuple[T, int]:
        left, right, vals, valsize = bst.left, bst.right, bst.val, bst.valsize
        if k < 0:
            k += len(bst)
        node = bst.root
        while True:
            t = vals[node] + valsize[left[node]]
            if t - vals[node] <= k < t:
                return bst.key[node], vals[node]
            if t > k:
                node = left[node]
            else:
                node = right[node]
                k -= t

    @staticmethod
    def contains(bst: BST, key: T) -> bool:
        left, right, keys = bst.left, bst.right, bst.key
        node = bst.root
        while node:
            if keys[node] == key:
                return True
            node = left[node] if key < keys[node] else right[node]
        return False

    @staticmethod
    def tolist(bst: BST) -> List[T]:
        left, right, keys, vals = bst.left, bst.right, bst.key, bst.val
        node = bst.root
        stack, a = [], newlist_hint(len(bst))
        while stack or node:
            if node:
                stack.append(node)
                node = left[node]
            else:
                node = stack.pop()
                key = keys[node]
                for _ in range(vals[node]):
                    a.append(key)
                node = right[node]
        return a


from typing import Generic, Iterable, Iterator, Tuple, TypeVar, List, Optional
from array import array

T = TypeVar("T", bound=SupportsLessThan)


class AVLTreeMultiset(OrderedMultisetInterface, Generic[T]):
    """
    多重集合としての AVL 木です。
    配列を用いてノードを表現しています。
    size を持ちます。
    """

    def __init__(self, a: Iterable[T] = []):
        self.root = 0
        self.key = [0]
        self.val = [0]
        self.valsize = [0]
        self.size = array("I", bytes(4))
        self.left = array("I", bytes(4))
        self.right = array("I", bytes(4))
        self.balance = array("b", bytes(1))
        self.end = 1
        if not isinstance(a, list):
            a = list(a)
        if a:
            self._build(a)

    def _make_node(self, key: T, val: int) -> int:
        end = self.end
        if end >= len(self.key):
            self.key.append(key)
            self.val.append(val)
            self.valsize.append(val)
            self.size.append(1)
            self.left.append(0)
            self.right.append(0)
            self.balance.append(0)
        else:
            self.key[end] = key
            self.val[end] = val
            self.valsize[end] = val
        self.end += 1
        return end

    def reserve(self, n: int) -> None:
        a = [0] * n
        self.key += a
        self.val += a
        self.valsize += a
        a = array("I", bytes(4 * n))
        self.left += a
        self.right += a
        self.size += array("I", [1] * n)
        self.balance += array("b", bytes(n))

    def _build(self, a: List[T]) -> None:
        left, right, size, valsize, balance = (
            self.left,
            self.right,
            self.size,
            self.valsize,
            self.balance,
        )

        def sort(l: int, r: int) -> Tuple[int, int]:
            mid = (l + r) >> 1
            node = mid
            hl, hr = 0, 0
            if l != mid:
                left[node], hl = sort(l, mid)
                size[node] += size[left[node]]
                valsize[node] += valsize[left[node]]
            if mid + 1 != r:
                right[node], hr = sort(mid + 1, r)
                size[node] += size[right[node]]
                valsize[node] += valsize[right[node]]
            balance[node] = hl - hr
            return node, max(hl, hr) + 1

        if not all(a[i] <= a[i + 1] for i in range(len(a) - 1)):
            a = sorted(a)
        if not a:
            return
        x, y = BSTMultisetArrayBase[AVLTreeMultiset, T]._rle(a)
        n = len(x)
        end = self.end
        self.end += n
        self.reserve(n)
        self.key[end : end + n] = x
        self.val[end : end + n] = y
        self.valsize[end : end + n] = y
        self.root = sort(end, n + end)[0]

    def _rotate_L(self, node: int) -> int:
        left, right, size, valsize, balance = (
            self.left,
            self.right,
            self.size,
            self.valsize,
            self.balance,
        )
        u = left[node]
        size[u] = size[node]
        valsize[u] = valsize[node]
        if left[u] == 0:
            size[node] -= 1
            valsize[node] -= self.val[u]
        else:
            size[node] -= size[left[u]] + 1
            valsize[node] -= valsize[left[u]] + self.val[u]
        left[node] = right[u]
        right[u] = node
        if balance[u] == 1:
            balance[u] = 0
            balance[node] = 0
        else:
            balance[u] = -1
            balance[node] = 1
        return u

    def _rotate_R(self, node: int) -> int:
        left, right, size, valsize, balance = (
            self.left,
            self.right,
            self.size,
            self.valsize,
            self.balance,
        )
        u = right[node]
        size[u] = size[node]
        valsize[u] = valsize[node]
        if right[u] == 0:
            size[node] -= 1
            valsize[node] -= self.val[u]
        else:
            size[node] -= size[right[u]] + 1
            valsize[node] -= valsize[right[u]] + self.val[u]
        right[node] = left[u]
        left[u] = node
        if balance[u] == -1:
            balance[u] = 0
            balance[node] = 0
        else:
            balance[u] = 1
            balance[node] = -1
        return u

    def _update_balance(self, node: int) -> None:
        left, right, balance = self.left, self.right, self.balance
        if balance[node] == 1:
            balance[right[node]] = -1
            balance[left[node]] = 0
        elif balance[node] == -1:
            balance[right[node]] = 0
            balance[left[node]] = 1
        else:
            balance[right[node]] = 0
            balance[left[node]] = 0
        balance[node] = 0

    def _rotate_LR(self, node: int) -> int:
        left, right, size, valsize = self.left, self.right, self.size, self.valsize
        B = left[node]
        E = right[B]
        size[E] = size[node]
        valsize[E] = valsize[node]
        if right[E] == 0:
            size[node] -= size[B]
            valsize[node] -= valsize[B]
            size[B] -= 1
            valsize[B] -= self.val[E]
        else:
            size[node] -= size[B] - size[right[E]]
            valsize[node] -= valsize[B] - valsize[right[E]]
            size[B] -= size[right[E]] + 1
            valsize[B] -= valsize[right[E]] + self.val[E]
        right[B] = left[E]
        left[E] = B
        left[node] = right[E]
        right[E] = node
        self._update_balance(E)
        return E

    def _rotate_RL(self, node: int) -> int:
        left, right, size, valsize = self.left, self.right, self.size, self.valsize
        C = right[node]
        D = left[C]
        size[D] = size[node]
        valsize[D] = valsize[node]
        if left[D] == 0:
            size[node] -= size[C]
            valsize[node] -= valsize[C]
            size[C] -= 1
            valsize[C] -= self.val[D]
        else:
            size[node] -= size[C] - size[left[D]]
            valsize[node] -= valsize[C] - valsize[left[D]]
            size[C] -= size[left[D]] + 1
            valsize[C] -= valsize[left[D]] + self.val[D]
        left[C] = right[D]
        right[D] = C
        right[node] = left[D]
        left[D] = node
        self._update_balance(D)
        return D

    def _kth_elm(self, k: int) -> Tuple[T, int]:
        return BSTMultisetArrayBase[AVLTreeMultiset, T]._kth_elm(self, k)

    def _kth_elm_tree(self, k: int) -> Tuple[T, int]:
        left, right, vals, size = self.left, self.right, self.val, self.size
        if k < 0:
            k += self.len_elm()
        assert 0 <= k < self.len_elm()
        node = self.root
        while True:
            t = size[left[node]]
            if t == k:
                return self.key[node], vals[node]
            if t > k:
                node = left[node]
            else:
                node = right[node]
                k -= t + 1

    def _discard(self, node: int, path: List[int], di: int) -> bool:
        left, right, keys, vals = self.left, self.right, self.key, self.val
        balance, size, valsize = self.balance, self.size, self.valsize
        fdi = 0
        if left[node] and right[node]:
            path.append(node)
            di <<= 1
            di |= 1
            lmax = left[node]
            while right[lmax]:
                path.append(lmax)
                di <<= 1
                fdi <<= 1
                fdi |= 1
                lmax = right[lmax]
            lmax_val = vals[lmax]
            keys[node] = keys[lmax]
            vals[node] = lmax_val
            node = lmax
        cnode = right[node] if left[node] == 0 else left[node]
        if path:
            if di & 1:
                left[path[-1]] = cnode
            else:
                right[path[-1]] = cnode
        else:
            self.root = cnode
            return True
        while path:
            new_node = 0
            pnode = path.pop()
            balance[pnode] -= 1 if di & 1 else -1
            size[pnode] -= 1
            valsize[pnode] -= lmax_val if fdi & 1 else 1
            di >>= 1
            fdi >>= 1
            if balance[pnode] == 2:
                new_node = (
                    self._rotate_LR(pnode) if balance[left[pnode]] < 0 else self._rotate_L(pnode)
                )
            elif balance[pnode] == -2:
                new_node = (
                    self._rotate_RL(pnode) if balance[right[pnode]] > 0 else self._rotate_R(pnode)
                )
            elif balance[pnode] != 0:
                break
            if new_node:
                if not path:
                    self.root = new_node
                    return
                if di & 1:
                    left[path[-1]] = new_node
                else:
                    right[path[-1]] = new_node
                if balance[new_node] != 0:
                    break
        while path:
            pnode = path.pop()
            size[pnode] -= 1
            valsize[pnode] -= lmax_val if fdi & 1 else 1
            fdi >>= 1
        return True

    def discard(self, key: T, val: int = 1) -> bool:
        keys, vals, left, right, valsize = self.key, self.val, self.left, self.right, self.valsize
        path = []
        di = 0
        node = self.root
        while node:
            if key == keys[node]:
                break
            path.append(node)
            di <<= 1
            if key < keys[node]:
                di |= 1
                node = left[node]
            else:
                node = right[node]
        else:
            return False
        if val > vals[node]:
            val = vals[node] - 1
            vals[node] -= val
            valsize[node] -= val
            for p in path:
                valsize[p] -= val
        if vals[node] == 1:
            self._discard(node, path, di)
        else:
            vals[node] -= val
            valsize[node] -= val
            for p in path:
                valsize[p] -= val
        return True

    def discard_all(self, key: T) -> None:
        self.discard(key, self.count(key))

    def remove(self, key: T, val: int = 1) -> None:
        if self.discard(key, val):
            return
        raise KeyError(key)

    def add(self, key: T, val: int = 1) -> None:
        if self.root == 0:
            self.root = self._make_node(key, val)
            return
        left, right, keys, valsize = self.left, self.right, self.key, self.valsize
        size, balance = self.size, self.balance
        node = self.root
        di = 0
        path = []
        while node:
            if key == keys[node]:
                self.val[node] += val
                valsize[node] += val
                for p in path:
                    valsize[p] += val
                return
            path.append(node)
            di <<= 1
            if key < keys[node]:
                di |= 1
                node = left[node]
            else:
                node = right[node]
        if di & 1:
            left[path[-1]] = self._make_node(key, val)
        else:
            right[path[-1]] = self._make_node(key, val)
        new_node = 0
        while path:
            node = path.pop()
            size[node] += 1
            valsize[node] += val
            balance[node] += 1 if di & 1 else -1
            di >>= 1
            if balance[node] == 0:
                break
            if balance[node] == 2:
                new_node = (
                    self._rotate_LR(node) if balance[left[node]] < 0 else self._rotate_L(node)
                )
                break
            elif balance[node] == -2:
                new_node = (
                    self._rotate_RL(node) if balance[right[node]] > 0 else self._rotate_R(node)
                )
                break
        if new_node:
            if path:
                if di & 1:
                    left[path[-1]] = new_node
                else:
                    right[path[-1]] = new_node
            else:
                self.root = new_node
        for p in path:
            size[p] += 1
            valsize[p] += val

    def count(self, key: T) -> int:
        return BSTMultisetArrayBase[AVLTreeMultiset, T].count(self, key)

    def le(self, key: T) -> Optional[T]:
        return BSTMultisetArrayBase[AVLTreeMultiset, T].le(self, key)

    def lt(self, key: T) -> Optional[T]:
        return BSTMultisetArrayBase[AVLTreeMultiset, T].lt(self, key)

    def ge(self, key: T) -> Optional[T]:
        return BSTMultisetArrayBase[AVLTreeMultiset, T].ge(self, key)

    def gt(self, key: T) -> Optional[T]:
        return BSTMultisetArrayBase[AVLTreeMultiset, T].gt(self, key)

    def index(self, key: T) -> int:
        return BSTMultisetArrayBase[AVLTreeMultiset, T].index(self, key)

    def index_right(self, key: T) -> int:
        return BSTMultisetArrayBase[AVLTreeMultiset, T].index_right(self, key)

    def index_keys(self, key: T) -> int:
        keys, left, right, vals, size = self.key, self.left, self.right, self.val, self.size
        k = 0
        node = self.root
        while node:
            if key == keys[node]:
                if left[node]:
                    k += size[left[node]]
                break
            if key < keys[node]:
                node = left[node]
            else:
                k += size[left[node]] + vals[node]
                node = right[node]
        return k

    def index_right_keys(self, key: T) -> int:
        keys, left, right, vals, size = self.key, self.left, self.right, self.val, self.size
        k = 0
        node = self.root
        while node:
            if key == keys[node]:
                k += size[left[node]] + vals[node]
                break
            if key < keys[node]:
                node = left[node]
            else:
                k += size[left[node]] + vals[node]
                node = right[node]
        return k

    def get_min(self) -> Optional[T]:
        if self.root == 0:
            return
        left = self.left
        node = self.root
        while left[node]:
            node = left[node]
        return self.key[node]

    def get_max(self) -> Optional[T]:
        if self.root == 0:
            return
        right = self.right
        node = self.root
        while right[node]:
            node = right[node]
        return self.key[node]

    def pop(self, k: int = -1) -> T:
        left, right, vals, valsize = self.left, self.right, self.val, self.valsize
        keys = self.key
        node = self.root
        if k < 0:
            k += valsize[node]
        path = []
        if k == valsize[node] - 1:
            while right[node]:
                path.append(node)
                node = right[node]
            x = keys[node]
            if vals[node] == 1:
                self._discard(node, path, 0)
            else:
                vals[node] -= 1
                valsize[node] -= 1
                for p in path:
                    valsize[p] -= 1
            return x
        di = 0
        while True:
            t = vals[node] + valsize[left[node]]
            if t - vals[node] <= k < t:
                x = keys[node]
                break
            path.append(node)
            di <<= 1
            if t > k:
                di |= 1
                node = left[node]
            else:
                node = right[node]
                k -= t
        if vals[node] == 1:
            self._discard(node, path, di)
        else:
            vals[node] -= 1
            valsize[node] -= 1
            for p in path:
                valsize[p] -= 1
        return x

    def pop_max(self) -> T:
        assert self
        return self.pop()

    def pop_min(self) -> T:
        assert self
        return self.pop(0)

    def items(self) -> Iterator[Tuple[T, int]]:
        for i in range(self.len_elm()):
            yield self._kth_elm_tree(i)

    def keys(self) -> Iterator[T]:
        for i in range(self.len_elm()):
            yield self._kth_elm_tree(i)[0]

    def values(self) -> Iterator[int]:
        for i in range(self.len_elm()):
            yield self._kth_elm_tree(i)[1]

    def len_elm(self) -> int:
        return self.size[self.root]

    def show(self) -> None:
        print("{" + ", ".join(map(lambda x: f"{x[0]}: {x[1]}", self.tolist_items())) + "}")

    def clear(self) -> None:
        self.root = 0

    def get_elm(self, k: int) -> T:
        return self._kth_elm_tree(k)[0]

    def tolist(self) -> List[T]:
        return BSTMultisetArrayBase[AVLTreeMultiset, T].tolist(self)

    def tolist_items(self) -> List[Tuple[T, int]]:
        left, right, keys, vals = self.left, self.right, self.key, self.val
        node = self.root
        stack, a = [], []
        while stack or node:
            if node:
                stack.append(node)
                node = left[node]
            else:
                node = stack.pop()
                a.append((keys[node], vals[node]))
                node = right[node]
        return a

    def __getitem__(self, k: int) -> T:
        return self._kth_elm(k)[0]

    def __contains__(self, key: T):
        return BSTMultisetArrayBase[AVLTreeMultiset, T].contains(self, key)

    def __iter__(self):
        self.__iter = 0
        return self

    def __next__(self):
        if self.__iter == len(self):
            raise StopIteration
        res = self._kth_elm(self.__iter)
        self.__iter += 1
        return res

    def __reversed__(self):
        for i in range(len(self)):
            yield self._kth_elm(-i - 1)[0]

    def __len__(self):
        return self.valsize[self.root]

    def __bool__(self):
        return self.root != 0

    def __str__(self):
        return "{" + ", ".join(map(str, self.tolist())) + "}"

    def __repr__(self):
        return f"{self.__class__.__name__}({self.tolist()})"
