# from titan_pylib.data_structures.bit_vector.avl_tree_bit_vector import AVLTreeBitVector
# from titan_pylib.data_structures.bit_vector.bit_vector_interface import BitVectorInterface
from abc import ABC, abstractmethod
from random import randint


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


from array import array
from typing import Iterable, List, Final, Sequence

titan_pylib_AVLTreeBitVector_W: Final[int] = 31


class AVLTreeBitVector(BitVectorInterface):
    """AVL木で書かれたビットベクトルです。簡潔でもなんでもありません。

    bit列を管理するわけですが、各節点は 1~32 bit を持つようにしています。
    これにより、最大 32 倍高速化が行えます。(16~32bitとするといいんだろうけど)
    """

    @staticmethod
    def _popcount(x: int) -> int:
        x = x - ((x >> 1) & 0x55555555)
        x = (x & 0x33333333) + ((x >> 2) & 0x33333333)
        x = x + (x >> 4) & 0x0F0F0F0F
        x += x >> 8
        x += x >> 16
        return x & 0x0000007F

    def __init__(self, a: Iterable[int] = []):
        """
        :math:`O(n)` です。

        Args:
          a (Iterable[int], optional): 構築元の配列です。
        """
        self.root = 0
        self.bit_len = array("B", bytes(1))
        self.key = array("I", bytes(4))
        self.size = array("I", bytes(4))
        self.total = array("I", bytes(4))
        self.left = array("I", bytes(4))
        self.right = array("I", bytes(4))
        self.balance = array("b", bytes(1))
        self.end = 1
        if a:
            self._build(a)

    def reserve(self, n: int) -> None:
        """``n`` 要素分のメモリを確保します。
        :math:`O(n)` です。
        """
        n = n // titan_pylib_AVLTreeBitVector_W + 1
        a = array("I", bytes(4 * n))
        self.bit_len += array("B", bytes(n))
        self.key += a
        self.size += a
        self.total += a
        self.left += a
        self.right += a
        self.balance += array("b", bytes(n))

    def _build(self, a: Iterable[int]) -> None:
        key, bit_len, left, right, size, balance, total = (
            self.key,
            self.bit_len,
            self.left,
            self.right,
            self.size,
            self.balance,
            self.total,
        )
        _popcount = AVLTreeBitVector._popcount

        def rec(lr: int) -> int:
            l, r = lr >> bit, lr & msk
            mid = (l + r) >> 1
            hl, hr = 0, 0
            if l != mid:
                le = rec(l << bit | mid)
                left[mid], hl = le >> bit, le & msk
                size[mid] += size[left[mid]]
                total[mid] += total[left[mid]]
            if mid + 1 != r:
                ri = rec((mid + 1) << bit | r)
                right[mid], hr = ri >> bit, ri & msk
                size[mid] += size[right[mid]]
                total[mid] += total[right[mid]]
            balance[mid] = hl - hr
            return mid << bit | (max(hl, hr) + 1)

        if not isinstance(a, Sequence):
            a = list(a)
        n = len(a)
        bit = n.bit_length() + 2
        msk = (1 << bit) - 1
        end = self.end
        self.reserve(n)
        i = 0
        indx = end
        for i in range(0, n, titan_pylib_AVLTreeBitVector_W):
            j = 0
            v = 0
            while j < titan_pylib_AVLTreeBitVector_W and i + j < n:
                v <<= 1
                v |= a[i + j]
                j += 1
            key[indx] = v
            bit_len[indx] = j
            size[indx] = j
            total[indx] = _popcount(v)
            indx += 1
        self.end = indx
        self.root = rec(end << bit | self.end) >> bit

    def _rotate_L(self, node: int) -> int:
        left, right, size, balance, total = (
            self.left,
            self.right,
            self.size,
            self.balance,
            self.total,
        )
        u = left[node]
        size[u] = size[node]
        total[u] = total[node]
        size[node] -= size[left[u]] + self.bit_len[u]
        total[node] -= total[left[u]] + AVLTreeBitVector._popcount(self.key[u])
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
        left, right, size, balance, total = (
            self.left,
            self.right,
            self.size,
            self.balance,
            self.total,
        )
        u = right[node]
        size[u] = size[node]
        total[u] = total[node]
        size[node] -= size[right[u]] + self.bit_len[u]
        total[node] -= total[right[u]] + AVLTreeBitVector._popcount(self.key[u])
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
        balance = self.balance
        if balance[node] == 1:
            balance[self.right[node]] = -1
            balance[self.left[node]] = 0
        elif balance[node] == -1:
            balance[self.right[node]] = 0
            balance[self.left[node]] = 1
        else:
            balance[self.right[node]] = 0
            balance[self.left[node]] = 0
        balance[node] = 0

    def _rotate_LR(self, node: int) -> int:
        left, right, size, total = self.left, self.right, self.size, self.total
        B = left[node]
        E = right[B]
        size[E] = size[node]
        size[node] -= size[B] - size[right[E]]
        size[B] -= size[right[E]] + self.bit_len[E]
        total[E] = total[node]
        total[node] -= total[B] - total[right[E]]
        total[B] -= total[right[E]] + AVLTreeBitVector._popcount(self.key[E])
        right[B] = left[E]
        left[E] = B
        left[node] = right[E]
        right[E] = node
        self._update_balance(E)
        return E

    def _rotate_RL(self, node: int) -> int:
        left, right, size, total = self.left, self.right, self.size, self.total
        C = right[node]
        D = left[C]
        size[D] = size[node]
        size[node] -= size[C] - size[left[D]]
        size[C] -= size[left[D]] + self.bit_len[D]
        total[D] = total[node]
        total[node] -= total[C] - total[left[D]]
        total[C] -= total[left[D]] + AVLTreeBitVector._popcount(self.key[D])
        left[C] = right[D]
        right[D] = C
        right[node] = left[D]
        left[D] = node
        self._update_balance(D)
        return D

    def _pref(self, r: int) -> int:
        left, right, bit_len, size, key, total = (
            self.left,
            self.right,
            self.bit_len,
            self.size,
            self.key,
            self.total,
        )
        node = self.root
        s = 0
        while r > 0:
            t = size[left[node]] + bit_len[node]
            if t - bit_len[node] < r <= t:
                r -= size[left[node]]
                s += total[left[node]] + AVLTreeBitVector._popcount(
                    key[node] >> (bit_len[node] - r)
                )
                break
            if t > r:
                node = left[node]
            else:
                s += total[left[node]] + AVLTreeBitVector._popcount(key[node])
                node = right[node]
                r -= t
        return s

    def _make_node(self, key: int, bit_len: int) -> int:
        end = self.end
        if end >= len(self.key):
            self.key.append(key)
            self.bit_len.append(bit_len)
            self.size.append(bit_len)
            self.total.append(AVLTreeBitVector._popcount(key))
            self.left.append(0)
            self.right.append(0)
            self.balance.append(0)
        else:
            self.key[end] = key
            self.bit_len[end] = bit_len
            self.size[end] = bit_len
            self.total[end] = AVLTreeBitVector._popcount(key)
        self.end += 1
        return end

    def insert(self, k: int, key: int) -> None:
        """``k`` 番目に ``v`` を挿入します。
        :math:`O(\\log{n})` です。

        Args:
          k (int): 挿入位置のインデックスです。
          key (int): 挿入する値です。 ``0`` または ``1`` である必要があります。
        """
        if self.root == 0:
            self.root = self._make_node(key, 1)
            return
        left, right, size, bit_len, balance, keys, total = (
            self.left,
            self.right,
            self.size,
            self.bit_len,
            self.balance,
            self.key,
            self.total,
        )
        node = self.root
        path = []
        d = 0
        while node:
            t = size[left[node]] + bit_len[node]
            if t - bit_len[node] <= k <= t:
                break
            d <<= 1
            size[node] += 1
            total[node] += key
            path.append(node)
            node = left[node] if t > k else right[node]
            if t > k:
                d |= 1
            else:
                k -= t
        k -= size[left[node]]
        if bit_len[node] < titan_pylib_AVLTreeBitVector_W:
            v = keys[node]
            bl = bit_len[node] - k
            keys[node] = (((v >> bl) << 1 | key) << bl) | (v & ((1 << bl) - 1))
            bit_len[node] += 1
            size[node] += 1
            total[node] += key
            return
        path.append(node)
        size[node] += 1
        total[node] += key
        v = keys[node]
        bl = titan_pylib_AVLTreeBitVector_W - k
        v = (((v >> bl) << 1 | key) << bl) | (v & ((1 << bl) - 1))
        left_key = v >> titan_pylib_AVLTreeBitVector_W
        left_key_popcount = left_key & 1
        keys[node] = v & ((1 << titan_pylib_AVLTreeBitVector_W) - 1)
        node = left[node]
        d <<= 1
        d |= 1
        if not node:
            if bit_len[path[-1]] < titan_pylib_AVLTreeBitVector_W:
                bit_len[path[-1]] += 1
                keys[path[-1]] = (keys[path[-1]] << 1) | left_key
                return
            else:
                left[path[-1]] = self._make_node(left_key, 1)
        else:
            path.append(node)
            size[node] += 1
            total[node] += left_key_popcount
            d <<= 1
            while right[node]:
                node = right[node]
                path.append(node)
                size[node] += 1
                total[node] += left_key_popcount
                d <<= 1
            if bit_len[node] < titan_pylib_AVLTreeBitVector_W:
                bit_len[node] += 1
                keys[node] = (keys[node] << 1) | left_key
                return
            else:
                right[node] = self._make_node(left_key, 1)
        new_node = 0
        while path:
            node = path.pop()
            balance[node] += 1 if d & 1 else -1
            d >>= 1
            if balance[node] == 0:
                break
            if balance[node] == 2:
                new_node = (
                    self._rotate_LR(node) if balance[left[node]] == -1 else self._rotate_L(node)
                )
                break
            elif balance[node] == -2:
                new_node = (
                    self._rotate_RL(node) if balance[right[node]] == 1 else self._rotate_R(node)
                )
                break
        if new_node:
            if path:
                if d & 1:
                    left[path[-1]] = new_node
                else:
                    right[path[-1]] = new_node
            else:
                self.root = new_node

    def _pop_under(self, path: List[int], d: int, node: int, res: int) -> None:
        left, right, size, bit_len, balance, keys, total = (
            self.left,
            self.right,
            self.size,
            self.bit_len,
            self.balance,
            self.key,
            self.total,
        )
        fd, lmax_total, lmax_bit_len = 0, 0, 0
        if left[node] and right[node]:
            path.append(node)
            d <<= 1
            d |= 1
            lmax = left[node]
            while right[lmax]:
                path.append(lmax)
                d <<= 1
                fd <<= 1
                fd |= 1
                lmax = right[lmax]
            lmax_total = AVLTreeBitVector._popcount(keys[lmax])
            lmax_bit_len = bit_len[lmax]
            keys[node] = keys[lmax]
            bit_len[node] = lmax_bit_len
            node = lmax
        cnode = right[node] if left[node] == 0 else left[node]
        if path:
            if d & 1:
                left[path[-1]] = cnode
            else:
                right[path[-1]] = cnode
        else:
            self.root = cnode
            return
        while path:
            new_node = 0
            node = path.pop()
            balance[node] -= 1 if d & 1 else -1
            size[node] -= lmax_bit_len if fd & 1 else 1
            total[node] -= lmax_total if fd & 1 else res
            d >>= 1
            fd >>= 1
            if balance[node] == 2:
                new_node = (
                    self._rotate_LR(node) if balance[left[node]] < 0 else self._rotate_L(node)
                )
            elif balance[node] == -2:
                new_node = (
                    self._rotate_RL(node) if balance[right[node]] > 0 else self._rotate_R(node)
                )
            elif balance[node] != 0:
                break
            if new_node:
                if not path:
                    self.root = new_node
                    return
                if d & 1:
                    left[path[-1]] = new_node
                else:
                    right[path[-1]] = new_node
                if balance[new_node] != 0:
                    break
        while path:
            node = path.pop()
            size[node] -= lmax_bit_len if fd & 1 else 1
            total[node] -= lmax_total if fd & 1 else res
            fd >>= 1

    def pop(self, k: int) -> int:
        """``k`` 番目の要素を削除し、その値を返します。
        :math:`O(\\log{n})` です。

        Args:
          k (int): 削除位置のインデックスです。
        """
        assert 0 <= k < len(self)
        left, right, size = self.left, self.right, self.size
        bit_len, keys, total = self.bit_len, self.key, self.total
        node = self.root
        d = 0
        path = []
        while node:
            t = size[left[node]] + bit_len[node]
            if t - bit_len[node] <= k < t:
                break
            path.append(node)
            node = left[node] if t > k else right[node]
            d <<= 1
            if t > k:
                d |= 1
            else:
                k -= t
        k -= size[left[node]]
        v = keys[node]
        res = v >> (bit_len[node] - k - 1) & 1
        if bit_len[node] == 1:
            self._pop_under(path, d, node, res)
            return res
        keys[node] = ((v >> (bit_len[node] - k)) << ((bit_len[node] - k - 1))) | (
            v & ((1 << (bit_len[node] - k - 1)) - 1)
        )
        bit_len[node] -= 1
        size[node] -= 1
        total[node] -= res
        for p in path:
            size[p] -= 1
            total[p] -= res
        return res

    def set(self, k: int, v: int) -> None:
        """``k`` 番目の値を ``v`` に更新します。
        :math:`O(\\log{n})` です。

        Args:
          k (int): 更新位置のインデックスです。
          key (int): 更新する値です。 ``0`` または ``1`` である必要があります。
        """
        self.__setitem__(k, v)

    def tolist(self) -> List[int]:
        """リストにして返します。
        :math:`O(n)` です。
        """
        left, right, key, bit_len = self.left, self.right, self.key, self.bit_len
        a = []
        if not self.root:
            return a

        def rec(node):
            if left[node]:
                rec(left[node])
            for i in range(bit_len[node] - 1, -1, -1):
                a.append(key[node] >> i & 1)
            if right[node]:
                rec(right[node])

        rec(self.root)
        return a

    def _debug_acc(self) -> None:
        """デバッグ用のメソッドです。
        key,totalをチェックします。
        """
        left, right = self.left, self.right
        key = self.key

        def rec(node):
            acc = self._popcount(key[node])
            if left[node]:
                acc += rec(left[node])
            if right[node]:
                acc += rec(right[node])
            if acc != self.total[node]:
                # self.debug()
                assert False, "acc Error"
            return acc

        rec(self.root)

    def access(self, k: int) -> int:
        """``k`` 番目の値を返します。
        :math:`O(\\log{n})` です。

        Args:
          k (int): 取得位置のインデックスです。
        """
        return self.__getitem__(k)

    def rank0(self, r: int) -> int:
        """``a[0, r)`` に含まれる ``0`` の個数を返します。
        :math:`O(\\log{n})` です。
        """
        return r - self._pref(r)

    def rank1(self, r: int) -> int:
        """``a[0, r)`` に含まれる ``1`` の個数を返します。
        :math:`O(\\log{n})` です。
        """
        return self._pref(r)

    def rank(self, r: int, v: int) -> int:
        """``a[0, r)`` に含まれる ``v`` の個数を返します。
        :math:`O(\\log{n})` です。
        """
        return self.rank1(r) if v else self.rank0(r)

    def select0(self, k: int) -> int:
        """``k`` 番目の ``0`` のインデックスを返します。
        :math:`O(\\log{n}^2)` です。
        """
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
        """``k`` 番目の ``1`` のインデックスを返します。
        :math:`O(\\log{n}^2)` です。
        """
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
        """``k`` 番目の ``v`` のインデックスを返します。
        :math:`O(\\log{n}^2)` です。
        """
        return self.select1(k) if v else self.select0(k)

    def _insert_and_rank1(self, k: int, key: int) -> int:
        if self.root == 0:
            self.root = self._make_node(key, 1)
            return 0
        left, right, size, bit_len, balance, keys, total = (
            self.left,
            self.right,
            self.size,
            self.bit_len,
            self.balance,
            self.key,
            self.total,
        )
        node = self.root
        s = 0
        path = []
        d = 0
        while node:
            t = size[left[node]] + bit_len[node]
            if t - bit_len[node] <= k <= t:
                break
            if t <= k:
                s += total[left[node]] + AVLTreeBitVector._popcount(keys[node])
            d <<= 1
            size[node] += 1
            total[node] += key
            path.append(node)
            node = left[node] if t > k else right[node]
            if t > k:
                d |= 1
            else:
                k -= t
        k -= size[left[node]]
        s += total[left[node]] + AVLTreeBitVector._popcount(keys[node] >> (bit_len[node] - k))
        if bit_len[node] < titan_pylib_AVLTreeBitVector_W:
            v = keys[node]
            bl = bit_len[node] - k
            keys[node] = (((v >> bl) << 1 | key) << bl) | (v & ((1 << bl) - 1))
            bit_len[node] += 1
            size[node] += 1
            total[node] += key
            return s
        path.append(node)
        size[node] += 1
        total[node] += key
        v = keys[node]
        bl = titan_pylib_AVLTreeBitVector_W - k
        v = (((v >> bl) << 1 | key) << bl) | (v & ((1 << bl) - 1))
        left_key = v >> titan_pylib_AVLTreeBitVector_W
        left_key_popcount = left_key & 1
        keys[node] = v & ((1 << titan_pylib_AVLTreeBitVector_W) - 1)
        node = left[node]
        d <<= 1
        d |= 1
        if not node:
            if bit_len[path[-1]] < titan_pylib_AVLTreeBitVector_W:
                bit_len[path[-1]] += 1
                keys[path[-1]] = (keys[path[-1]] << 1) | left_key
                return s
            else:
                left[path[-1]] = self._make_node(left_key, 1)
        else:
            path.append(node)
            size[node] += 1
            total[node] += left_key_popcount
            d <<= 1
            while right[node]:
                node = right[node]
                path.append(node)
                size[node] += 1
                total[node] += left_key_popcount
                d <<= 1
            if bit_len[node] < titan_pylib_AVLTreeBitVector_W:
                bit_len[node] += 1
                keys[node] = (keys[node] << 1) | left_key
                return s
            else:
                right[node] = self._make_node(left_key, 1)
        new_node = 0
        while path:
            node = path.pop()
            balance[node] += 1 if d & 1 else -1
            d >>= 1
            if balance[node] == 0:
                break
            if balance[node] == 2:
                new_node = (
                    self._rotate_LR(node) if balance[left[node]] == -1 else self._rotate_L(node)
                )
                break
            elif balance[node] == -2:
                new_node = (
                    self._rotate_RL(node) if balance[right[node]] == 1 else self._rotate_R(node)
                )
                break
        if new_node:
            if path:
                if d & 1:
                    left[path[-1]] = new_node
                else:
                    right[path[-1]] = new_node
            else:
                self.root = new_node
        return s

    def _access_pop_and_rank1(self, k: int) -> int:
        assert 0 <= k < len(self)
        left, right, size = self.left, self.right, self.size
        bit_len, keys, total = self.bit_len, self.key, self.total
        s = 0
        node = self.root
        d = 0
        path = []
        while node:
            t = size[left[node]] + bit_len[node]
            if t - bit_len[node] <= k < t:
                break
            if t <= k:
                s += total[left[node]] + AVLTreeBitVector._popcount(keys[node])
            path.append(node)
            node = left[node] if t > k else right[node]
            d <<= 1
            if t > k:
                d |= 1
            else:
                k -= t
        k -= size[left[node]]
        s += total[left[node]] + AVLTreeBitVector._popcount(keys[node] >> (bit_len[node] - k))
        v = keys[node]
        res = v >> (bit_len[node] - k - 1) & 1
        if bit_len[node] == 1:
            self._pop_under(path, d, node, res)
            return s << 1 | res
        keys[node] = ((v >> (bit_len[node] - k)) << ((bit_len[node] - k - 1))) | (
            v & ((1 << (bit_len[node] - k - 1)) - 1)
        )
        bit_len[node] -= 1
        size[node] -= 1
        total[node] -= res
        for p in path:
            size[p] -= 1
            total[p] -= res
        return s << 1 | res

    def __getitem__(self, k: int) -> int:
        """``k`` 番目の要素を返します。
        :math:`O(\\log{n})` です。
        """
        assert 0 <= k < len(self)
        left, right, bit_len, size, key = self.left, self.right, self.bit_len, self.size, self.key
        node = self.root
        while True:
            t = size[left[node]] + bit_len[node]
            if t - bit_len[node] <= k < t:
                k -= size[left[node]]
                return key[node] >> (bit_len[node] - k - 1) & 1
            if t > k:
                node = left[node]
            else:
                node = right[node]
                k -= t

    def __setitem__(self, k: int, v: int) -> None:
        """``k`` 番目の要素を ``v`` に更新します。
        :math:`O(\\log{n})` です。
        """
        left, right, bit_len, size, key, total = (
            self.left,
            self.right,
            self.bit_len,
            self.size,
            self.key,
            self.total,
        )
        assert v == 0 or v == 1, "ValueError"
        node = self.root
        path = []
        while True:
            t = size[left[node]] + bit_len[node]
            path.append(node)
            if t - bit_len[node] <= k < t:
                k -= size[left[node]]
                if v:
                    key[node] |= 1 << k
                else:
                    key[node] &= ~(1 << k)
                break
            elif t > k:
                node = left[node]
            else:
                node = right[node]
                k -= t
        while path:
            node = path.pop()
            total[node] = (
                AVLTreeBitVector._popcount(key[node]) + total[left[node]] + total[right[node]]
            )

    def __str__(self):
        return str(self.tolist())

    def __len__(self):
        return self.size[self.root]

    def __repr__(self):
        return f"{self.__class__.__name__}({self})"


if __name__ == "__main__":
    wm = AVLTreeBitVector([0] * 5)
    # wm.insert(0, 1)
    print(wm)
    # wm.pop(2)
    # print(wm)
    wm.set(1, 1)
    print(wm)

    print(all(a == b for a, b in zip(wm.tolist(), [wm[i] for i in range(len(wm))])))

    for _ in range(10):
        n = randint(1, 100) + 100
        a = [randint(0, 1) for _ in range(n)]
        wm = AVLTreeBitVector(a)

        def countBf(nums):
            return sum(nums)

        def kthBf(nums, k, v):
            cnt = 0
            for i in range(len(nums)):
                if nums[i] == v:
                    if cnt == k:
                        return i
                    cnt += 1
            return -1

        def insertBf(nums, k, v):
            nums.insert(k, v)

        def popBf(nums, k):
            return nums.pop(k)

        for _ in range(100):
            # count
            for i in range(len(a) + 1):
                assert wm.rank(i, 0) == i - countBf(a[:i])
                assert wm.rank(i, 1) == countBf(a[:i])
                assert wm.select(i, 0) == kthBf(a, i, 0)
                assert wm.select(i, 1) == kthBf(a, i, 1)

            # insert
            insertIndex = randint(0, len(a))
            insertValue = randint(0, 1)
            insertBf(a, insertIndex, insertValue)
            wm.insert(insertIndex, insertValue)

            # set
            # setIndex = randint(0, len(a) - 1)
            # setValue = randint(0, 1)
            # a[setIndex] = setValue
            # wm.set(setIndex, setValue)

            # pop
            popIndex = randint(0, len(a) - 1)
            assert popBf(a, popIndex) == wm.pop(popIndex)

            # len
            assert len(a) == len(wm)

            # get
            for i in range(len(a)):
                assert a[i] == wm[i]

            # tolist
            assert a == wm.tolist()

            wm._debug_acc()

    print("ok")
