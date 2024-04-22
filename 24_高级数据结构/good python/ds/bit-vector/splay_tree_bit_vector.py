# from titan_pylib.data_structures.bit_vector.splay_tree_bit_vector import SplayTreeBitVector
# from titan_pylib.data_structures.bit_vector.bit_vector_interface import BitVectorInterface
from abc import ABC, abstractmethod


class BitVectorInterface(ABC):
    @abstractmethod
    def access(self, k: int) -> int:
        raise NotImplementedError

    @abstractmethod
    def __getitem__(self, k: int) -> int:
        raise NotImplementedError

    @abstractmethod
    def rank0(self, r: int) -> int:
        raise NotImplementedError

    @abstractmethod
    def rank1(self, r: int) -> int:
        raise NotImplementedError

    @abstractmethod
    def rank(self, r: int, v: int) -> int:
        raise NotImplementedError

    @abstractmethod
    def select0(self, k: int) -> int:
        raise NotImplementedError

    @abstractmethod
    def select1(self, k: int) -> int:
        raise NotImplementedError

    @abstractmethod
    def select(self, k: int, v: int) -> int:
        raise NotImplementedError

    @abstractmethod
    def __len__(self) -> int:
        raise NotImplementedError

    @abstractmethod
    def __str__(self) -> str:
        raise NotImplementedError

    @abstractmethod
    def __repr__(self) -> str:
        raise NotImplementedError


from typing import Sequence, List
from array import array


class SplayTreeBitVector(BitVectorInterface):
    @staticmethod
    def _popcount(x: int) -> int:
        x = x - ((x >> 1) & 0x55555555)
        x = (x & 0x33333333) + ((x >> 2) & 0x33333333)
        x = x + (x >> 4) & 0x0F0F0F0F
        x += x >> 8
        x += x >> 16
        return x & 0x0000007F

    def __init__(self, a: Sequence[int] = []):
        self.root = 0
        self.bit_len = array("B", bytes(1))
        self.key = array("I", bytes(4))
        self.size = array("I", bytes(4))
        self.total = array("I", bytes(4))
        self.child = array("I", bytes(8))
        self.end = 1
        self.w = 32
        if a:
            self._build(a)

    def reserve(self, n: int) -> None:
        n = n // self.w + 1
        a = array("I", bytes(4 * n))
        self.bit_len += array("B", bytes(n))
        self.key += a
        self.size += a
        self.total += a
        self.child += array("I", bytes(8 * n))

    def _build(self, a: Sequence[int]) -> None:
        key, bit_len, child, size, total = self.key, self.bit_len, self.child, self.size, self.total
        _popcount = SplayTreeBitVector._popcount

        def rec(l: int, r: int) -> int:
            mid = (l + r) >> 1
            if l != mid:
                child[mid << 1] = rec(l, mid)
                size[mid] += size[child[mid << 1]]
                total[mid] += total[child[mid << 1]]
            if mid + 1 != r:
                child[mid << 1 | 1] = rec(mid + 1, r)
                size[mid] += size[child[mid << 1 | 1]]
                total[mid] += total[child[mid << 1 | 1]]
            return mid

        if not (hasattr(a, "__getitem__") and hasattr(a, "__len__")):
            a = list(a)
        n = len(a)
        end = self.end
        self.reserve(n)
        i = 0
        indx = end
        for i in range(0, n, self.w):
            j = 0
            v = 0
            while j < self.w and i + j < n:
                v <<= 1
                v |= a[i + j]
                j += 1
            key[indx] = v
            bit_len[indx] = j
            size[indx] = j
            total[indx] = _popcount(v)
            indx += 1
        self.end = indx
        self.root = rec(end, self.end)

    def _make_node(self, key: int, bit_len: int) -> int:
        end = self.end
        if end >= len(self.key):
            self.key.append(key)
            self.bit_len.append(bit_len)
            self.size.append(bit_len)
            self.total.append(SplayTreeBitVector._popcount(key))
            self.child.append(0)
            self.child.append(0)
        else:
            self.key[end] = key
            self.bit_len[end] = bit_len
            self.size[end] = bit_len
            self.total[end] = SplayTreeBitVector._popcount(key)
        self.end += 1
        return end

    def _update_triple(self, x: int, y: int, z: int) -> None:
        child, bit_len, size, total = self.child, self.bit_len, self.size, self.total
        lx, rx = child[x << 1], child[x << 1 | 1]
        ly, ry = child[y << 1], child[y << 1 | 1]
        size[z] = size[x]
        size[x] = bit_len[x] + size[lx] + size[rx]
        size[y] = bit_len[y] + size[ly] + size[ry]
        total[z] = total[x]
        total[x] = total[lx] + SplayTreeBitVector._popcount(self.key[x]) + total[rx]
        total[y] = total[ly] + SplayTreeBitVector._popcount(self.key[y]) + total[ry]

    def _update_double(self, x: int, y: int) -> None:
        child, bit_len, size, total = self.child, self.bit_len, self.size, self.total
        lx, rx = child[x << 1], child[x << 1 | 1]
        size[y] = size[x]
        size[x] = bit_len[x] + size[lx] + size[rx]
        total[y] = total[x]
        total[x] = total[lx] + SplayTreeBitVector._popcount(self.key[x]) + total[rx]

    def _update(self, node: int) -> None:
        lnode, rnode = self.child[node << 1], self.child[node << 1 | 1]
        self.size[node] = self.bit_len[node] + self.size[lnode] + self.size[rnode]
        self.total[node] = (
            SplayTreeBitVector._popcount(self.key[node]) + self.total[lnode] + self.total[rnode]
        )

    def _splay(self, path: List[int], d: int) -> None:
        child = self.child
        g = d & 1
        while len(path) > 1:
            pnode = path.pop()
            gnode = path.pop()
            f = d >> 1 & 1
            node = child[pnode << 1 | g ^ 1]
            nnode = (pnode if g == f else node) << 1 | f
            child[pnode << 1 | g ^ 1] = child[node << 1 | g]
            child[node << 1 | g] = pnode
            child[gnode << 1 | f ^ 1] = child[nnode]
            child[nnode] = gnode
            self._update_triple(gnode, pnode, node)
            if not path:
                return
            d >>= 2
            g = d & 1
            child[path[-1] << 1 | g ^ 1] = node
        pnode = path.pop()
        node = child[pnode << 1 | g ^ 1]
        child[pnode << 1 | g ^ 1] = child[node << 1 | g]
        child[node << 1 | g] = pnode
        self._update_double(pnode, node)

    def _kth_elm_splay(self, node: int, k: int) -> int:
        child, bit_len, size = self.child, self.bit_len, self.size
        d = 0
        path = []
        while True:
            t = size[child[node << 1]] + bit_len[node]
            if t - bit_len[node] <= k < t:
                if path:
                    self._splay(path, d)
                return node
            d = d << 1 | (t > k)
            path.append(node)
            node = child[node << 1 | (t <= k)]
            if t <= k:
                k -= t

    def _left_splay(self, node: int) -> int:
        if not node:
            return 0
        child = self.child
        if not child[node << 1]:
            return node
        path = []
        while child[node << 1]:
            path.append(node)
            node = child[node << 1]
        self._splay(path, (1 << len(path)) - 1)
        return node

    def _right_splay(self, node: int) -> int:
        if not node:
            return 0
        child = self.child
        if not child[node << 1 | 1]:
            return node
        path = []
        while child[node << 1 | 1]:
            path.append(node)
            node = child[node << 1 | 1]
        self._splay(path, 0)
        return node

    def insert(self, index: int, key: int) -> None:
        assert (
            0 <= index <= len(self)
        ), f"IndexError: SplayTreeBitVector.insert({index}, {key}), len={len(self)}"
        if not self.root:
            node = self._make_node(key, 1)
            self.root = node
            return
        bit_len, child, size, keys, total = (
            self.bit_len,
            self.child,
            self.size,
            self.key,
            self.total,
        )
        if index == size[self.root]:
            node = self._right_splay(self.root)
            if bit_len[node] == self.w:
                tmp = keys[node] << 1 | key
                new_node = self._make_node(tmp & 1, 1)
                keys[node] = tmp >> 1
                child[new_node << 1] = node
                self._update(node)
                size[new_node] += size[node]
                total[new_node] += total[node]
                self.root = new_node
            else:
                tmp = keys[node]
                bl = index - bit_len[node] - size[child[node << 1]]
                keys[node] = (((tmp >> bl) << 1 | key) << bl) | (tmp & ((1 << bl) - 1))
                bit_len[node] += 1
                size[node] += 1
                total[node] += key
                self.root = node
        else:
            node = self._kth_elm_splay(self.root, index)
            if bit_len[node] == self.w:
                index -= size[child[node << 1]]
                tmp = keys[node]
                bl = bit_len[node] - index
                tmp = (((tmp >> bl) << 1 | key) << bl) | (tmp & ((1 << bl) - 1))
                new_node = self._make_node(tmp >> self.w, 1)
                keys[node] = tmp & ((1 << self.w) - 1)
                self._update(node)
                if child[node << 1]:
                    child[new_node << 1] = child[node << 1]
                    child[node << 1] = 0
                    self._update(node)
                child[new_node << 1 | 1] = node
                self._update(new_node)
                self.root = new_node
            else:
                tmp = keys[node]
                bl = bit_len[node] - index + size[child[node << 1]]
                keys[node] = (((tmp >> bl) << 1 | key) << bl) | (tmp & ((1 << bl) - 1))
                bit_len[node] += 1
                size[node] += 1
                total[node] += key
                self.root = node

    def pop(self, index: int = -1) -> int:
        assert 0 <= index < len(self), f"IndexError: SplayTreeBitVector.pop({index})"
        root = self._kth_elm_splay(self.root, index)
        size, child, key, bit_len, total = self.size, self.child, self.key, self.bit_len, self.total
        index -= size[child[root << 1]]
        v = key[root]
        res = v >> (bit_len[root] - index - 1) & 1
        if bit_len[root] == 1:
            if not child[root << 1]:
                self.root = child[root << 1 | 1]
            elif not child[root << 1 | 1]:
                self.root = child[root << 1]
            else:
                node = self._right_splay(child[root << 1])
                child[node << 1 | 1] = child[root << 1 | 1]
                self._update(node)
                self.root = node
        else:
            key[root] = ((v >> (bit_len[root] - index)) << ((bit_len[root] - index - 1))) | (
                v & ((1 << (bit_len[root] - index - 1)) - 1)
            )
            bit_len[root] -= 1
            size[root] -= 1
            total[root] -= res
            self.root = root
        return res

    def _pref(self, r: int) -> int:
        assert 0 <= r <= len(self), f"IndexError: SplayTreeBitVector._pref({r}), len={len(self)}"
        if r == 0:
            return 0
        if r == len(self):
            return self.total[self.root]
        self.root = self._kth_elm_splay(self.root, r - 1)
        r -= self.size[self.child[self.root << 1]]
        return (
            self.total[self.root]
            - SplayTreeBitVector._popcount(
                self.key[self.root] & ((1 << (self.bit_len[self.root] - r)) - 1)
            )
            - self.total[self.child[self.root << 1 | 1]]
        )

    def __getitem__(self, k: int) -> int:
        assert 0 <= k < len(self), f"IndexError: SplayTreeBitVector.__getitem__({k})"
        self.root = self._kth_elm_splay(self.root, k)
        k -= self.size[self.child[self.root << 1]]
        return (self.key[self.root] >> (self.bit_len[self.root] - k - 1)) & 1

    def debug(self):
        print("### debug")
        print(f"{self.root=}")
        print(f"{self.key=}")
        print(f"{self.bit_len=}")
        print(f"{self.size=}")
        print(f"{self.total=}")
        print(f"{self.child=}")

    def __len__(self):
        return self.size[self.root]

    def tolist(self) -> List[int]:
        child, key, bit_len = self.child, self.key, self.bit_len
        a = []
        if not self.root:
            return a

        def rec(node):
            if child[node << 1]:
                rec(child[node << 1])
            for i in range(bit_len[node] - 1, -1, -1):
                a.append(key[node] >> i & 1)
            if child[node << 1 | 1]:
                rec(child[node << 1 | 1])

        rec(self.root)
        return a

    def __str__(self):
        return str(self.tolist())

    __repr__ = __str__

    def debug_acc(self) -> None:
        child = self.child
        key = self.key

        def rec(node):
            acc = self._popcount(key[node])
            if child[node << 1]:
                acc += rec(child[node << 1])
            if child[node << 1 | 1]:
                acc += rec(child[node << 1 | 1])
            if acc != self.total[node]:
                # self.debug()
                assert False, "acc Error"
            return acc

        rec(self.root)

    def access(self, k: int) -> int:
        return self.__getitem__(k)

    def rank0(self, r: int) -> int:
        # a[0, r) に含まれる 0 の個数
        return r - self._pref(r)

    def rank1(self, r: int) -> int:
        # a[0, r) に含まれる 1 の個数
        return self._pref(r)

    def rank(self, r: int, v: int) -> int:
        # a[0, r) に含まれる v の個数
        return self.rank1(r) if v else self.rank0(r)

    def select0(self, k: int) -> int:
        # k 番目の 0 のindex
        # O(log(N))
        if k < 0 or self.rank0(len(self)) <= k:
            return -1
        l, r = 0, len(self)
        while r - l > 1:
            m = (l + r) >> 1
            if m - self._pref(m) > k:
                r = m
            else:
                l = m
        return l

    def select1(self, k: int) -> int:
        # k 番目の 1 のindex
        # O(log(N))
        if k < 0 or self.rank1(len(self)) <= k:
            return -1
        l, r = 0, len(self)
        while r - l > 1:
            m = (l + r) >> 1
            if self._pref(m) > k:
                r = m
            else:
                l = m
        return l

    def select(self, k: int, v: int) -> int:
        # k 番目の v のindex
        # O(log(N))
        return self.select1(k) if v else self.select0(k)
