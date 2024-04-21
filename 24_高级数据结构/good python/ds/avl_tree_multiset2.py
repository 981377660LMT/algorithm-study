# 多重集合としての AVL 木です。 配列を用いてノードを表現しています。 size を持たないので軽めです。

# from titan_pylib.data_structures.avl_tree.avl_tree_multiset2 import AVLTreeMultiset2
# from titan_pylib.my_class.supports_less_than import SupportsLessThan
from typing import Protocol


class SupportsLessThan(Protocol):
    def __lt__(self, other) -> bool:
        ...


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


from typing import Generic, Iterable, Tuple, TypeVar, List, Optional
from array import array

T = TypeVar("T", bound=SupportsLessThan)


class AVLTreeMultiset2(Generic[T]):
    """
    多重集合としての AVL 木です。
    配列を用いてノードを表現しています。
    size を持たないので軽めです。
    """

    def __init__(self, a: Iterable[T] = []):
        self.root = 0
        self._len = 0
        self.key = [0]
        self.val = [0]
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
            self.left.append(0)
            self.right.append(0)
            self.balance.append(0)
        else:
            self.key[end] = key
            self.val[end] = val
        self.end += 1
        return end

    def reserve(self, n: int) -> None:
        a = [0] * n
        self.key += a
        self.val += a
        a = array("I", bytes(4 * n))
        self.left += a
        self.right += a
        self.balance += array("b", bytes(n))

    def _build(self, a: List[T]) -> None:
        left, right, balance = self.left, self.right, self.balance

        def sort(l: int, r: int) -> Tuple[int, int]:
            mid = (l + r) >> 1
            node = mid
            hl, hr = 0, 0
            if l != mid:
                left[node], hl = sort(l, mid)
            if mid + 1 != r:
                right[node], hr = sort(mid + 1, r)
            balance[node] = hl - hr
            return node, max(hl, hr) + 1

        self._len = len(a)
        if not all(a[i] < a[i + 1] for i in range(len(a) - 1)):
            a = sorted(a)
        x, y = BSTMultisetArrayBase[AVLTreeMultiset2, T]._rle(a)
        n = len(x)
        end = self.end
        self.end += n
        self.reserve(n)
        self.key[end : end + n] = x
        self.val[end : end + n] = y
        self.root = sort(end, n + end)[0]

    def _rotate_L(self, node: int) -> int:
        left, right, balance = self.left, self.right, self.balance
        u = left[node]
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
        left, right, balance = self.left, self.right, self.balance
        u = right[node]
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
        left, right = self.left, self.right
        B = left[node]
        E = right[B]
        right[B] = left[E]
        left[E] = B
        left[node] = right[E]
        right[E] = node
        self._update_balance(E)
        return E

    def _rotate_RL(self, node: int) -> int:
        left, right = self.left, self.right
        C = right[node]
        D = left[C]
        left[C] = right[D]
        right[D] = C
        right[node] = left[D]
        left[D] = node
        self._update_balance(D)
        return D

    def _discard(self, node: int, path: List[int], di: int) -> bool:
        left, right, keys, vals, balance = self.left, self.right, self.key, self.val, self.balance
        if left[node] and right[node]:
            path.append(node)
            di <<= 1
            di |= 1
            lmax = left[node]
            while right[lmax]:
                path.append(lmax)
                di <<= 1
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
            di >>= 1
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
        return True

    def discard(self, key: T, val: int = 1) -> bool:
        keys, vals, left, right = self.key, self.val, self.left, self.right
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
        self._len -= min(val, vals[node])
        if val > vals[node]:
            val = vals[node] - 1
            vals[node] -= val
        if vals[node] == 1:
            self._discard(node, path, di)
        else:
            vals[node] -= val
        return True

    def discard_all(self, key: T) -> None:
        self.discard(key, self.count(key))

    def remove(self, key: T, val: int = 1) -> None:
        if self.discard(key, val):
            return
        raise KeyError(key)

    def add(self, key: T, val: int = 1) -> None:
        self._len += val
        if self.root == 0:
            self.root = self._make_node(key, val)
            return
        left, right, keys, balance = self.left, self.right, self.key, self.balance
        node = self.root
        di = 0
        path = []
        while node:
            if key == keys[node]:
                self.val[node] += val
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

    def count(self, key: T) -> int:
        return BSTMultisetArrayBase[AVLTreeMultiset2, T].count(self, key)

    def le(self, key: T) -> Optional[T]:
        return BSTMultisetArrayBase[AVLTreeMultiset2, T].le(self, key)

    def lt(self, key: T) -> Optional[T]:
        return BSTMultisetArrayBase[AVLTreeMultiset2, T].lt(self, key)

    def ge(self, key: T) -> Optional[T]:
        return BSTMultisetArrayBase[AVLTreeMultiset2, T].ge(self, key)

    def gt(self, key: T) -> Optional[T]:
        return BSTMultisetArrayBase[AVLTreeMultiset2, T].gt(self, key)

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

    def pop_min(self) -> T:
        left, vals, keys = self.left, self.val, self.key
        self._len -= 1
        node = self.root
        path = []
        while left[node]:
            path.append(node)
            node = left[node]
        x = keys[node]
        if vals[node] == 1:
            self._discard(node, path, (1 << len(path)) - 1)
        else:
            vals[node] -= 1
        return x

    def pop_max(self) -> T:
        right, vals, keys = self.right, self.val, self.key
        self._len -= 1
        node = self.root
        path = []
        while right[node]:
            path.append(node)
            node = right[node]
        x = keys[node]
        if vals[node] == 1:
            self._discard(node, path, 0)
        else:
            vals[node] -= 1
        return x

    def clear(self) -> None:
        self.root = 0

    def tolist(self) -> List[T]:
        return BSTMultisetArrayBase[AVLTreeMultiset2, T].tolist(self)

    def tolist_items(self) -> List[Tuple[T, int]]:
        left, right, keys, vals = self.left, self.right, self.key, self.val
        node = self.root
        stack: List[int] = []
        a: List[Tuple[T, int]] = []
        while stack or node:
            if node:
                stack.append(node)
                node = left[node]
            else:
                node = stack.pop()
                a.append((keys[node], vals[node]))
                node = right[node]
        return a

    def __contains__(self, key: T):
        return BSTMultisetArrayBase[AVLTreeMultiset2, T].contains(self, key)

    def __len__(self):
        return self._len

    def __bool__(self):
        return self.root != 0

    def __str__(self):
        return "{" + ", ".join(map(str, self.tolist())) + "}"

    def __repr__(self):
        return f"{self.__class__.__name__}({self.tolist()})"
