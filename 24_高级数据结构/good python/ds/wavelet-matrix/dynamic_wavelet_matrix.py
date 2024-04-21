# from titan_pylib.data_structures.wavelet_matrix.dynamic_wavelet_matrix import DynamicWaveletMatrix
# from titan_pylib.data_structures.bit_vector.avl_tree_bit_vector import AVLTreeBitVector
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
        print("debug_acc ok.")

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


# from titan_pylib.data_structures.wavelet_matrix.wavelet_matrix import WaveletMatrix
# from titan_pylib.data_structures.bit_vector.bit_vector import BitVector
from array import array


class BitVector(BitVectorInterface):
    """コンパクトな bit vector です。"""

    def __init__(self, n: int):
        """長さ ``n`` の ``BitVector`` です。

        bit を保持するのに ``array[I]`` を使用します。
        ``block_size= n / 32`` として、使用bitは ``32*block_size=2n bit`` です。

        累積和を保持するのに同様の ``array[I]`` を使用します。
        32bitごとの和を保存しています。同様に使用bitは ``2n bit`` です。
        """
        assert 0 <= n < 4294967295
        self.N = n
        self.block_size = (n + 31) >> 5
        b = bytes(4 * (self.block_size + 1))
        self.bit = array("I", b)
        self.acc = array("I", b)

    @staticmethod
    def _popcount(x: int) -> int:
        x = x - ((x >> 1) & 0x55555555)
        x = (x & 0x33333333) + ((x >> 2) & 0x33333333)
        x = x + (x >> 4) & 0x0F0F0F0F
        x += x >> 8
        x += x >> 16
        return x & 0x0000007F

    def set(self, k: int) -> None:
        """``k`` 番目の bit を ``1`` にします。
        :math:`O(1)` です。

        Args:
          k (int): インデックスです。
        """
        self.bit[k >> 5] |= 1 << (k & 31)

    def build(self) -> None:
        """構築します。
        **これ以降 ``set`` メソッドを使用してはいけません。**
        :math:`O(n)` です。
        """
        acc, bit = self.acc, self.bit
        for i in range(self.block_size):
            acc[i + 1] = acc[i] + BitVector._popcount(bit[i])

    def access(self, k: int) -> int:
        """``k`` 番目の bit を返します。
        :math:`O(1)` です。
        """
        return (self.bit[k >> 5] >> (k & 31)) & 1

    def __getitem__(self, k: int) -> int:
        return (self.bit[k >> 5] >> (k & 31)) & 1

    def rank0(self, r: int) -> int:
        """``a[0, r)`` に含まれる ``0`` の個数を返します。
        :math:`O(1)` です。
        """
        return r - (
            self.acc[r >> 5] + BitVector._popcount(self.bit[r >> 5] & ((1 << (r & 31)) - 1))
        )

    def rank1(self, r: int) -> int:
        """``a[0, r)`` に含まれる ``1`` の個数を返します。
        :math:`O(1)` です。
        """
        return self.acc[r >> 5] + BitVector._popcount(self.bit[r >> 5] & ((1 << (r & 31)) - 1))

    def rank(self, r: int, v: int) -> int:
        """``a[0, r)`` に含まれる ``v`` の個数を返します。
        :math:`O(1)` です。
        """
        return self.rank1(r) if v else self.rank0(r)

    def select0(self, k: int) -> int:
        """``k`` 番目の ``0`` のインデックスを返します。
        :math:`O(\\log{n})` です。
        """
        if k < 0 or self.rank0(self.N) <= k:
            return -1
        l, r = 0, self.block_size + 1
        while r - l > 1:
            m = (l + r) >> 1
            if m * 32 - self.acc[m] > k:
                r = m
            else:
                l = m
        indx = 32 * l
        k = k - (l * 32 - self.acc[l]) + self.rank0(indx)
        l, r = indx, indx + 32
        while r - l > 1:
            m = (l + r) >> 1
            if self.rank0(m) > k:
                r = m
            else:
                l = m
        return l

    def select1(self, k: int) -> int:
        """``k`` 番目の ``1`` のインデックスを返します。
        :math:`O(\\log{n})` です。
        """
        if k < 0 or self.rank1(self.N) <= k:
            return -1
        l, r = 0, self.block_size + 1
        while r - l > 1:
            m = (l + r) >> 1
            if self.acc[m] > k:
                r = m
            else:
                l = m
        indx = 32 * l
        k = k - self.acc[l] + self.rank1(indx)
        l, r = indx, indx + 32
        while r - l > 1:
            m = (l + r) >> 1
            if self.rank1(m) > k:
                r = m
            else:
                l = m
        return l

    def select(self, k: int, v: int) -> int:
        """``k`` 番目の ``v`` のインデックスを返します。
        :math:`O(\\log{n})` です。
        """
        return self.select1(k) if v else self.select0(k)

    def __len__(self):
        return self.N

    def __str__(self):
        return str([self.access(i) for i in range(self.N)])

    def __repr__(self):
        return f"{self.__class__.__name__}({self})"


from typing import Sequence, List, Tuple
from heapq import heappush, heappop
from array import array


class WaveletMatrix:
    """``WaveletMatrix`` です。
    静的であることに注意してください。

    以下の仕様の計算量には嘘があるかもしれません。import 元の ``BitVector`` の計算量も参考にしてください。

    - 参考:
      - `https://miti-7.hatenablog.com/entry/2018/04/28/152259 <https://miti-7.hatenablog.com/entry/2018/04/28/152259>`
      - `https://www.slideshare.net/pfi/ss-15916040 <https://www.slideshare.net/pfi/ss-15916040>`
      - `デwiki <https://scrapbox.io/data-structures/Wavelet_Matrix>`
    """

    def __init__(self, sigma: int, a: Sequence[int] = []):
        """``[0, sigma)`` の整数列を管理する ``WaveletMatrix`` を構築します。
        :math:`O(n\\log{\\sigma})` です。

        Args:
          sigma (int): 扱う整数の上限です。
          a (Sequence[int], optional): 構築する配列です。
        """
        self.sigma: int = sigma
        self.log: int = (sigma - 1).bit_length()
        self.mid: array[int] = array("I", bytes(4 * self.log))
        self.size: int = len(a)
        self.v: List[BitVector] = [BitVector(self.size) for _ in range(self.log)]
        self._build(a)

    def _build(self, a: Sequence[int]) -> None:
        # 列 a から wm を構築する
        for bit in range(self.log - 1, -1, -1):
            # bit目の0/1に応じてvを構築 + aを安定ソート
            v = self.v[bit]
            zero, one = [], []
            for i, e in enumerate(a):
                if e >> bit & 1:
                    v.set(i)
                    one.append(e)
                else:
                    zero.append(e)
            v.build()
            self.mid[bit] = len(zero)  # 境界をmid[bit]に保持
            a = zero + one

    def access(self, k: int) -> int:
        """k番目の値を返します。
        :math:`O(\\log{\\sigma})` です。

        Args:
          k (int): インデックスです。
        """
        assert (
            -self.size <= k < self.size
        ), f"IndexError: {self.__class__.__name__}.access({k}), size={self.size}"
        if k < 0:
            k += self.size
        s = 0  # 答え
        for bit in range(self.log - 1, -1, -1):
            if self.v[bit].access(k):
                # k番目が立ってたら、
                # kまでの1とすべての0が次のk
                s |= 1 << bit
                k = self.v[bit].rank1(k) + self.mid[bit]
            else:
                # kまでの0が次のk
                k = self.v[bit].rank0(k)
        return s

    def __getitem__(self, k: int) -> int:
        assert (
            -self.size <= k < self.size
        ), f"IndexError: {self.__class__.__name__}[{k}], size={self.size}"
        return self.access(k)

    def rank(self, r: int, x: int) -> int:
        """``a[0, r)`` に含まれる ``x`` の個数を返します。
        :math:`O(\\log{\\sigma})` です。
        """
        assert (
            0 <= r <= self.size
        ), f"IndexError: {self.__class__.__name__}.rank(), r={r}, size={self.size}"
        assert (
            0 <= x < 1 << self.log
        ), f"ValueError: {self.__class__.__name__}.rank(), x={x}, LIM={1<<self.log}"
        l = 0
        mid = self.mid
        for bit in range(self.log - 1, -1, -1):
            # 位置 r より左に x が何個あるか
            # x の bit 目で場合分け
            if x >> bit & 1:
                # 立ってたら、次のl, rは以下
                l = self.v[bit].rank1(l) + mid[bit]
                r = self.v[bit].rank1(r) + mid[bit]
            else:
                # そうでなければ次のl, rは以下
                l = self.v[bit].rank0(l)
                r = self.v[bit].rank0(r)
        return r - l

    def select(self, k: int, x: int) -> int:
        """``k`` 番目の ``v`` のインデックスを返します。
        :math:`O(\\log{\\sigma})` です。
        """
        assert (
            0 <= k < self.size
        ), f"IndexError: {self.__class__.__name__}.select({k}, {x}), k={k}, size={self.size}"
        assert (
            0 <= x < 1 << self.log
        ), f"ValueError: {self.__class__.__name__}.select({k}, {x}), x={x}, LIM={1<<self.log}"
        # x の開始位置 s を探す
        s = 0
        for bit in range(self.log - 1, -1, -1):
            if x >> bit & 1:
                s = self.v[bit].rank0(self.size) + self.v[bit].rank1(s)
            else:
                s = self.v[bit].rank0(s)
        s += k  # s から k 進んだ位置が、元の列で何番目か調べる
        for bit in range(self.log):
            if x >> bit & 1:
                s = self.v[bit].select1(s - self.v[bit].rank0(self.size))
            else:
                s = self.v[bit].select0(s)
        return s

    def kth_smallest(self, l: int, r: int, k: int) -> int:
        """``a[l, r)`` の中で k 番目に **小さい** 値を返します。
        :math:`O(\\log{\\sigma})` です。
        """
        assert (
            0 <= l <= r <= self.size
        ), f"IndexError: {self.__class__.__name__}.kth_smallest({l}, {r}, {k}), size={self.size}"
        assert (
            0 <= k < r - l
        ), f"IndexError: {self.__class__.__name__}.kth_smallest({l}, {r}, {k}), wrong k"
        s = 0
        mid = self.mid
        for bit in range(self.log - 1, -1, -1):
            r0, l0 = self.v[bit].rank0(r), self.v[bit].rank0(l)
            cnt = r0 - l0  # 区間内の 0 の個数
            if cnt <= k:  # 0 が k 以下のとき、 k 番目は 1
                s |= 1 << bit
                k -= cnt
                # この 1 が次の bit 列でどこに行くか
                l = l - l0 + mid[bit]
                r = r - r0 + mid[bit]
            else:
                # この 0 が次の bit 列でどこに行くか
                l = l0
                r = r0
        return s

    quantile = kth_smallest

    def kth_largest(self, l: int, r: int, k: int) -> int:
        """``a[l, r)`` の中で k 番目に **大きい値** を返します。
        :math:`O(\\log{\\sigma})` です。
        """
        assert (
            0 <= l <= r <= self.size
        ), f"IndexError: {self.__class__.__name__}.kth_largest({l}, {r}, {k}), size={self.size}"
        assert (
            0 <= k < r - l
        ), f"IndexError: {self.__class__.__name__}.kth_largest({l}, {r}, {k}), wrong k"
        return self.kth_smallest(l, r, r - l - k - 1)

    def topk(self, l: int, r: int, k: int) -> List[Tuple[int, int]]:
        """``a[l, r)`` の中で、要素を出現回数が多い順にその頻度とともに ``k`` 個返します。
        :math:`O(\\min(r-l, \\sigam) \\log(\\sigam))` です。

        Note:
          :math:`\\sigma` が大きい場合、計算量に注意です。

        Returns:
          List[Tuple[int, int]]: ``(要素, 頻度)`` を要素とする配列です。
        """
        assert (
            0 <= l <= r <= self.size
        ), f"IndexError: {self.__class__.__name__}.topk({l}, {r}, {k}), size={self.size}"
        assert 0 <= k < r - l, f"IndexError: {self.__class__.__name__}.topk({l}, {r}, {k}), wrong k"
        # heap[-length, x, l, bit]
        hq: List[Tuple[int, int, int, int]] = [(-(r - l), 0, l, self.log - 1)]
        ans = []
        while hq:
            length, x, l, bit = heappop(hq)
            length = -length
            if bit == -1:
                ans.append((x, length))
                k -= 1
                if k == 0:
                    break
            else:
                r = l + length
                l0 = self.v[bit].rank0(l)
                r0 = self.v[bit].rank0(r)
                if l0 < r0:
                    heappush(hq, (-(r0 - l0), x, l0, bit - 1))
                l1 = self.v[bit].rank1(l) + self.mid[bit]
                r1 = self.v[bit].rank1(r) + self.mid[bit]
                if l1 < r1:
                    heappush(hq, (-(r1 - l1), x | (1 << bit), l1, bit - 1))
        return ans

    def sum(self, l: int, r: int) -> int:
        """``topk`` メソッドを用いて ``a[l, r)`` の総和を返します。

        計算量に注意です。
        """
        assert False, "Yabai Keisanryo Error"
        return sum(k * v for k, v in self.topk(l, r, r - l))

    def _range_freq(self, l: int, r: int, x: int) -> int:
        """a[l, r) で x 未満の要素の数を返す"""
        ans = 0
        for bit in range(self.log - 1, -1, -1):
            l0, r0 = self.v[bit].rank0(l), self.v[bit].rank0(r)
            if x >> bit & 1:
                # bit が立ってたら、区間の 0 の個数を答えに加算し、新たな区間は 1 のみ
                ans += r0 - l0
                # 1 が次の bit 列でどこに行くか
                l += self.mid[bit] - l0
                r += self.mid[bit] - r0
            else:
                # 0 が次の bit 列でどこに行くか
                l, r = l0, r0
        return ans

    def range_freq(self, l: int, r: int, x: int, y: int) -> int:
        """``a[l, r)`` に含まれる、 ``x`` 以上 ``y`` 未満である要素の個数を返します。
        :math:`O(\\log{\\sigma})` です。
        """
        assert (
            0 <= l <= r <= self.size
        ), f"IndexError: {self.__class__.__name__}.range_freq({l}, {r}, {x}, {y})"
        return self._range_freq(l, r, y) - self._range_freq(l, r, x)

    def prev_value(self, l: int, r: int, x: int) -> int:
        """``a[l, r)`` で、``x`` 以上 ``y`` 未満であるような要素のうち最大の要素を返します。
        :math:`O(\\log{\\sigma})` です。
        """
        assert (
            0 <= l <= r <= self.size
        ), f"IndexError: {self.__class__.__name__}.prev_value({l}, {r}, {x})"
        return self.kth_smallest(l, r, self._range_freq(l, r, x) - 1)

    def next_value(self, l: int, r: int, x: int) -> int:
        """``a[l, r)`` で、``x`` 以上 ``y`` 未満であるような要素のうち最小の要素を返します。
        :math:`O(\\log{\\sigma})` です。
        """
        assert (
            0 <= l <= r <= self.size
        ), f"IndexError: {self.__class__.__name__}.next_value({l}, {r}, {x})"
        return self.kth_smallest(l, r, self._range_freq(l, r, x))

    def range_count(self, l: int, r: int, x: int) -> int:
        """``a[l, r)`` に含まれる ``x`` の個数を返します。
        ``wm.rank(r, x) - wm.rank(l, x)`` と等価です。
        :math:`O(\\log{\\sigma})` です。
        """
        assert (
            0 <= l <= r <= self.size
        ), f"IndexError: {self.__class__.__name__}.range_count({l}, {r}, {x})"
        return self.rank(r, x) - self.rank(l, x)

    def __len__(self):
        return self.size

    def __str__(self):
        return f"{self.__class__.__name__}({[self.access(i) for i in range(self.size)]})"

    __repr__ = __str__


from typing import Sequence, List
from array import array


class DynamicWaveletMatrix(WaveletMatrix):
    """動的ウェーブレット行列です。

    (静的)ウェーブレット行列でできる操作に加えて ``insert / pop / update`` 等ができます。
      - ``BitVector`` を平衡二分木にしています(``AVLTreeBitVector``)。あらゆる操作に平衡二分木の log がつきます。これヤバくね

    :math:`O(n\\log{(\\sigma)})` です。
    """

    def __init__(self, sigma: int, a: Sequence[int] = []) -> None:
        self.sigma: int = sigma
        self.log: int = (sigma - 1).bit_length()
        self.v: List[AVLTreeBitVector] = [AVLTreeBitVector()] * self.log
        self.mid: array[int] = array("I", bytes(4 * self.log))
        self.size: int = len(a)
        self._build(a)

    def _build(self, a: Sequence[int]) -> None:
        v = array("B", bytes(self.size))
        for bit in range(self.log - 1, -1, -1):
            # bit目の0/1に応じてvを構築 + aを安定ソート
            zero, one = [], []
            for i, e in enumerate(a):
                if e >> bit & 1:
                    v[i] = 1
                    one.append(e)
                else:
                    v[i] = 0
                    zero.append(e)
            self.mid[bit] = len(zero)  # 境界をmid[bit]に保持
            self.v[bit] = AVLTreeBitVector(v)
            a = zero + one

    def reserve(self, n: int) -> None:
        """``n`` 要素分のメモリを確保します。
        :math:`O(n)` です。
        """
        assert n >= 0, f"ValueError: {self.__class__.__name__}.reserve({n})"
        for v in self.v:
            v.reserve(n)

    def insert(self, k: int, x: int) -> None:
        """位置 ``k`` に ``x`` を挿入します。
        :math:`O(\\log{(n)}\\log{(\\sigma)})` です。
        """
        assert (
            0 <= k <= self.size
        ), f"IndexError: {self.__class__.__name__}.insert({k}, {x}), n={self.size}"
        assert (
            0 <= x < 1 << self.log
        ), f"ValueError: {self.__class__.__name__}.insert({k}, {x}), LIM={1<<self.log}"
        mid = self.mid
        for bit in range(self.log - 1, -1, -1):
            v = self.v[bit]
            # if x >> bit & 1:
            #   v.insert(k, 1)
            #   k = v.rank1(k) + mid[bit]
            # else:
            #   v.insert(k, 0)
            #   mid[bit] += 1
            #   k = v.rank0(k)
            if x >> bit & 1:
                s = v._insert_and_rank1(k, 1)
                k = s + mid[bit]
            else:
                s = v._insert_and_rank1(k, 0)
                k -= s
                mid[bit] += 1
        self.size += 1

    def pop(self, k: int) -> int:
        """位置 ``k`` の要素を削除し、その値を返します。
        :math:`O(\\log{(n)}\\log{(\\sigma)})` です。
        """
        assert 0 <= k < self.size, f"IndexError: {self.__class__.__name__}.pop({k}), n={self.size}"
        mid = self.mid
        ans = 0
        for bit in range(self.log - 1, -1, -1):
            v = self.v[bit]
            # K = k
            # if v.access(k):
            #   ans |= 1 << bit
            #   k = v.rank1(k) + mid[bit]
            # else:
            #   mid[bit] -= 1
            #   k = v.rank0(k)
            # v.pop(K)
            sb = v._access_pop_and_rank1(k)
            s = sb >> 1
            if sb & 1:
                ans |= 1 << bit
                k = s + mid[bit]
            else:
                mid[bit] -= 1
                k -= s
        self.size -= 1
        return ans

    def update(self, k: int, x: int) -> None:
        """位置 ``k`` の要素を ``x`` に更新します。
        :math:`O(\\log{(n)}\\log{(\\sigma)})` です。
        """
        assert (
            0 <= k < self.size
        ), f"IndexError: {self.__class__.__name__}.update({k}, {x}), n={self.size}"
        assert (
            0 <= x < 1 << self.log
        ), f"ValueError: {self.__class__.__name__}.update({k}, {x}), LIM={1<<self.log}"
        self.pop(k)
        self.insert(k, x)

    def __setitem__(self, k: int, x: int):
        assert (
            0 <= k < self.size
        ), f"IndexError: {self.__class__.__name__}[{k}] = {x}, n={self.size}"
        assert (
            0 <= x < 1 << self.log
        ), f"ValueError: {self.__class__.__name__}[{k}] = {x}, LIM={1<<self.log}"
        self.update(k, x)

    def __str__(self):
        return f"{self.__class__.__name__}({[self[i] for i in range(self.size)]})"
