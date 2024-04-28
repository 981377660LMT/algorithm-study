# from titan_pylib.string.hash_string import HashString
# ref: https://qiita.com/keymoon/items/11fac5627672a6d6a9f6
# from titan_pylib.data_structures.segment_tree.segment_tree import SegmentTree
# from titan_pylib.data_structures.segment_tree.segment_tree_interface import SegmentTreeInterface
from abc import ABC, abstractmethod
from typing import TypeVar, Generic, Union, Iterable, Callable, List

T = TypeVar("T")


class SegmentTreeInterface(ABC, Generic[T]):
    @abstractmethod
    def __init__(self, n_or_a: Union[int, Iterable[T]], op: Callable[[T, T], T], e: T):
        raise NotImplementedError

    @abstractmethod
    def set(self, k: int, v: T) -> None:
        raise NotImplementedError

    @abstractmethod
    def get(self, k: int) -> T:
        raise NotImplementedError

    @abstractmethod
    def prod(self, l: int, r: int) -> T:
        raise NotImplementedError

    @abstractmethod
    def all_prod(self) -> T:
        raise NotImplementedError

    @abstractmethod
    def max_right(self, l: int, f: Callable[[T], bool]) -> int:
        raise NotImplementedError

    @abstractmethod
    def min_left(self, r: int, f: Callable[[T], bool]) -> int:
        raise NotImplementedError

    @abstractmethod
    def tolist(self) -> List[T]:
        raise NotImplementedError

    @abstractmethod
    def __getitem__(self, k: int) -> T:
        raise NotImplementedError

    @abstractmethod
    def __setitem__(self, k: int, v: T) -> None:
        raise NotImplementedError

    @abstractmethod
    def __str__(self):
        raise NotImplementedError

    @abstractmethod
    def __repr__(self):
        raise NotImplementedError


from typing import Generic, Iterable, TypeVar, Callable, Union, List

T = TypeVar("T")


class SegmentTree(SegmentTreeInterface, Generic[T]):
    """セグ木です。非再帰です。"""

    def __init__(self, n_or_a: Union[int, Iterable[T]], op: Callable[[T, T], T], e: T) -> None:
        """``SegmentTree`` を構築します。
        :math:`O(n)` です。

        Args:
          n_or_a (Union[int, Iterable[T]]): ``n: int`` のとき、 ``e`` を初期値として長さ ``n`` の ``SegmentTree`` を構築します。
                                            ``a: Iterable[T]`` のとき、 ``a`` から ``SegmentTree`` を構築します。
          op (Callable[[T, T], T]): 2項演算の関数です。
          e (T): 単位元です。
        """
        self._op = op
        self._e = e
        if isinstance(n_or_a, int):
            self._n = n_or_a
            self._log = (self._n - 1).bit_length()
            self._size = 1 << self._log
            self._data = [e] * (self._size << 1)
        else:
            n_or_a = list(n_or_a)
            self._n = len(n_or_a)
            self._log = (self._n - 1).bit_length()
            self._size = 1 << self._log
            _data = [e] * (self._size << 1)
            _data[self._size : self._size + self._n] = n_or_a
            for i in range(self._size - 1, 0, -1):
                _data[i] = op(_data[i << 1], _data[i << 1 | 1])
            self._data = _data

    def set(self, k: int, v: T) -> None:
        """一点更新です。
        :math:`O(\\log{n})` です。

        Args:
          k (int): 更新するインデックスです。
          v (T): 更新する値です。

        制約:
          :math:`-n \\leq n \\leq k < n`
        """
        assert (
            -self._n <= k < self._n
        ), f"IndexError: {self.__class__.__name__}.set({k}, {v}), n={self._n}"
        if k < 0:
            k += self._n
        k += self._size
        self._data[k] = v
        for _ in range(self._log):
            k >>= 1
            self._data[k] = self._op(self._data[k << 1], self._data[k << 1 | 1])

    def get(self, k: int) -> T:
        """一点取得です。
        :math:`O(1)` です。

        Args:
          k (int): インデックスです。

        制約:
          :math:`-n \\leq n \\leq k < n`
        """
        assert (
            -self._n <= k < self._n
        ), f"IndexError: {self.__class__.__name__}.get({k}), n={self._n}"
        if k < 0:
            k += self._n
        return self._data[k + self._size]

    def prod(self, l: int, r: int) -> T:
        """区間 ``[l, r)`` の総積を返します。
        :math:`O(\\log{n})` です。

        Args:
          l (int): インデックスです。
          r (int): インデックスです。

        制約:
          :math:`0 \\leq l \\leq r \\leq n`
        """
        assert 0 <= l <= r <= self._n, f"IndexError: {self.__class__.__name__}.prod({l}, {r})"
        l += self._size
        r += self._size
        lres = self._e
        rres = self._e
        while l < r:
            if l & 1:
                lres = self._op(lres, self._data[l])
                l += 1
            if r & 1:
                rres = self._op(self._data[r ^ 1], rres)
            l >>= 1
            r >>= 1
        return self._op(lres, rres)

    def all_prod(self) -> T:
        """区間 ``[0, n)`` の総積を返します。
        :math:`O(1)` です。
        """
        return self._data[1]

    def max_right(self, l: int, f: Callable[[T], bool]) -> int:
        """Find the largest index R s.t. f([l, R)) == True. / O(\\log{n})"""
        assert (
            0 <= l <= self._n
        ), f"IndexError: {self.__class__.__name__}.max_right({l}, f) index out of range"
        # assert f(self._e), \
        #     f'{self.__class__.__name__}.max_right({l}, f), f({self._e}) must be true.'
        if l == self._n:
            return self._n
        l += self._size
        s = self._e
        while True:
            while l & 1 == 0:
                l >>= 1
            if not f(self._op(s, self._data[l])):
                while l < self._size:
                    l <<= 1
                    if f(self._op(s, self._data[l])):
                        s = self._op(s, self._data[l])
                        l |= 1
                return l - self._size
            s = self._op(s, self._data[l])
            l += 1
            if l & -l == l:
                break
        return self._n

    def min_left(self, r: int, f: Callable[[T], bool]) -> int:
        """Find the smallest index L s.t. f([L, r)) == True. / O(\\log{n})"""
        assert (
            0 <= r <= self._n
        ), f"IndexError: {self.__class__.__name__}.min_left({r}, f) index out of range"
        # assert f(self._e), \
        #     f'{self.__class__.__name__}.min_left({r}, f), f({self._e}) must be true.'
        if r == 0:
            return 0
        r += self._size
        s = self._e
        while True:
            r -= 1
            while r > 1 and r & 1:
                r >>= 1
            if not f(self._op(self._data[r], s)):
                while r < self._size:
                    r = r << 1 | 1
                    if f(self._op(self._data[r], s)):
                        s = self._op(self._data[r], s)
                        r ^= 1
                return r + 1 - self._size
            s = self._op(self._data[r], s)
            if r & -r == r:
                break
        return 0

    def tolist(self) -> List[T]:
        """リストにして返します。
        :math:`O(n)` です。
        """
        return [self.get(i) for i in range(self._n)]

    def show(self) -> None:
        """デバッグ用のメソッドです。"""
        print(
            f"<{self.__class__.__name__}> [\n"
            + "\n".join(
                [
                    "  " + " ".join(map(str, [self._data[(1 << i) + j] for j in range(1 << i)]))
                    for i in range(self._log + 1)
                ]
            )
            + "\n]"
        )

    def __getitem__(self, k: int) -> T:
        assert (
            -self._n <= k < self._n
        ), f"IndexError: {self.__class__.__name__}.__getitem__({k}), n={self._n}"
        return self.get(k)

    def __setitem__(self, k: int, v: T):
        assert (
            -self._n <= k < self._n
        ), f"IndexError: {self.__class__.__name__}.__setitem__{k}, {v}), n={self._n}"
        self.set(k, v)

    def __len__(self) -> int:
        return self._n

    def __str__(self) -> str:
        return str(self.tolist())

    def __repr__(self) -> str:
        return f"{self.__class__.__name__}({self})"


from typing import Optional, List, Dict, Final
import random
import string

_titan_pylib_HashString_MOD: Final[int] = (1 << 61) - 1
_titan_pylib_HashString_DIC: Final[Dict[str, int]] = {
    c: i for i, c in enumerate(string.ascii_lowercase, 1)
}
_titan_pylib_HashString_MASK30: Final[int] = (1 << 30) - 1
_titan_pylib_HashString_MASK31: Final[int] = (1 << 31) - 1
_titan_pylib_HashString_MASK61: Final[int] = _titan_pylib_HashString_MOD


class HashStringBase:
    """HashStringのベースクラスです。"""

    def __init__(self, n: int, base: int = -1, seed: Optional[int] = None) -> None:
        """
        :math:`O(n)` です。

        Args:
          n (int): 文字列の長さの上限です。
          base (int, optional): Defaults to -1.
          seed (Optional[int], optional): Defaults to None.
        """
        random.seed(seed)
        base = random.randint(37, 10**9) if base < 0 else base
        powb = [1] * (n + 1)
        invb = [1] * (n + 1)
        invbpow = pow(base, -1, _titan_pylib_HashString_MOD)
        for i in range(1, n + 1):
            powb[i] = HashStringBase.get_mul(powb[i - 1], base)
            invb[i] = HashStringBase.get_mul(invb[i - 1], invbpow)
        print(powb)
        print(invb)
        print(invbpow)
        self.n = n
        self.powb = powb
        self.invb = invb

    @staticmethod
    def get_mul(a: int, b: int) -> int:
        au = a >> 31
        ad = a & _titan_pylib_HashString_MASK31
        bu = b >> 31
        bd = b & _titan_pylib_HashString_MASK31
        mid = ad * bu + au * bd
        midu = mid >> 30
        midd = mid & _titan_pylib_HashString_MASK30
        return HashStringBase.get_mod(au * bu * 2 + midu + (midd << 31) + ad * bd)

    @staticmethod
    def get_mod(x: int) -> int:
        xu = x >> 61
        xd = x & _titan_pylib_HashString_MASK61
        res = xu + xd
        if res >= _titan_pylib_HashString_MOD:
            res -= _titan_pylib_HashString_MOD

        return res

    def unite(self, h1: int, h2: int, k: int) -> int:
        # len(h2) == k
        # h1 <- h2
        return self.get_mod(self.get_mul(h1, self.powb[k]) + h2)


class HashString:
    def __init__(self, hsb: HashStringBase, s: str, update: bool = False) -> None:
        """ロリハを構築します。
        :math:`O(n)` です。

        Args:
          hsb (HashStringBase): ベースクラスです。
          s (str): ロリハを構築する文字列です。
          update (bool, optional): ``update=True`` のとき、1点更新が可能になります。
        """
        n = len(s)
        data = [0] * n
        acc = [0] * (n + 1)
        powb = hsb.powb
        for i, c in enumerate(s):
            data[i] = hsb.get_mul(powb[n - i - 1], ord(c))
            acc[i + 1] = hsb.get_mod(acc[i] + data[i])
        self.hsb = hsb
        self.n = n
        self.acc = acc
        self.used_seg = False
        if update:
            self.seg = SegmentTree(data, lambda s, t: (s + t) % _titan_pylib_HashString_MOD, 0)

    def get(self, l: int, r: int) -> int:
        """``s[l, r)`` のハッシュ値を返します。
        1点更新処理後は :math:`O(\\log{n})` 、そうでなければ :math:`O(1)` です。

        Args:
          l (int): インデックスです。
          r (int): インデックスです。

        Returns:
          int: ハッシュ値です。
        """
        if self.used_seg:
            return self.hsb.get_mul(self.seg.prod(l, r), self.hsb.invb[self.n - r])
        return self.hsb.get_mul(
            self.hsb.get_mod(self.acc[r] - self.acc[l]), self.hsb.invb[self.n - r]
        )

    def __getitem__(self, k: int) -> int:
        """``s[k]`` のハッシュ値を返します。
        1点更新処理後は :math:`O(\\log{n})` 、そうでなければ :math:`O(1)` です。

        Args:
          k (int): インデックスです。

        Returns:
          int: ハッシュ値です。
        """
        return self.get(k, k + 1)

    def set(self, k: int, c: str) -> None:
        """`k` 番目の文字を `c` に更新します。
        :math:`O(\\log{n})` です。また、今後の ``get()`` が :math:`O(\\log{n})` になります。

        Args:
          k (int): インデックスです。
          c (str): 更新する文字です。
        """
        self.used_seg = True
        self.seg[k] = self.hsb.get_mul(self.hsb.powb[self.n - k - 1], ord(c))

    def __setitem__(self, k: int, c: str) -> None:
        return self.set(k, c)

    def __len__(self):
        return self.n

    def get_lcp(self) -> List[int]:
        """lcp配列を返します。
        :math:`O(n\\log{n})` です。
        """
        a = [0] * self.n
        memo = [-1] * (self.n + 1)
        for i in range(self.n):
            ok, ng = 0, self.n - i + 1
            while ng - ok > 1:
                mid = (ok + ng) >> 1
                if memo[mid] == -1:
                    memo[mid] = self.get(0, mid)
                if memo[mid] == self.get(i, i + mid):
                    ok = mid
                else:
                    ng = mid
            a[i] = ok
        return a


if __name__ == "__main__":
    s = "asezfvgbadpihoamgkcmco"
    base = HashStringBase(len(s), 37)
    hs = HashString(base, s)
    for i in range(4):
        print(hs.get(i, i + 1))

    class Solution:
        def sumScores(self, s: str) -> int:
            def countPre(curLen: int, start: int) -> int:
                left, right = 1, curLen
                while left <= right:
                    mid = (left + right) // 2
                    if hasher.get(start, start + mid) == hasher.get(0, mid):
                        left = mid + 1
                    else:
                        right = mid - 1
                return right

            n = len(s)
            base = HashStringBase(n, -1)
            hasher = HashString(base, s)
            res = 0
            for i in range(1, n + 1):
                if s[-i] != s[0]:
                    continue
                count = countPre(i, n - i)
                res += count
            return res
