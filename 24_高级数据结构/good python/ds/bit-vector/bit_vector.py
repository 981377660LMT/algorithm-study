# from titan_pylib.data_structures.bit_vector.bit_vector import BitVector
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
