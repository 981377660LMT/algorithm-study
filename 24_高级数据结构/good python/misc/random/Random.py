# from titan_pylib.algorithm.random.random import Random
from typing import List, Any


class Random:
    """Random
    乱数系のライブラリです。
    標準ライブラリよりも高速なつもりでいます。
    """

    def __init__(self) -> None:
        self._x = 123456789
        self._y = 362436069
        self._z = 521288629
        self._w = 88675123

    def _xor(self) -> int:
        t = (self._x ^ ((self._x << 11) & 0xFFFFFFFF)) & 0xFFFFFFFF
        self._x, self._y, self._z = self._y, self._z, self._w
        self._w = (self._w ^ (self._w >> 19)) ^ (t ^ ((t >> 8)) & 0xFFFFFFFF) & 0xFFFFFFFF
        return self._w

    def random(self) -> float:
        """0以上1以下の一様ランダムな値を1つ生成して返すはずです。
        :math:`O(1)` です。
        """
        return self._xor() / 0xFFFFFFFF

    def randint(self, a: int, b: int) -> int:
        """``a`` 以上 ``b`` **以下** のランダムな整数を返します。
        :math:`O(1)` です。

        制約:
          :math:`a \\leq b`
        """
        assert a <= b
        return a + self._xor() % (b - a + 1)

    def randrange(self, begin: int, end: int) -> int:
        """``begin`` 以上 ``end`` **未満** のランダムな整数を返します。
        :math:`O(1)` です。

        制約:
          :math:`begin < end`
        """
        assert begin < end
        return begin + self._xor() % (end - begin)

    def shuffle(self, a: List[Any]) -> None:
        """``a`` をインプレースにシャッフルします。
        :math:`O(n)` です。

        Args:
          a (List[Any]): ``a`` をシャッフルします。
        """
        n = len(a)
        for i in range(n - 1):
            j = self.randrange(i, n)
            a[i], a[j] = a[j], a[i]
