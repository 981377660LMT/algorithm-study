# from titan_pylib.algorithm.mo import Mo
from typing import Callable
from itertools import chain
from math import sqrt, ceil


class Mo:
    """長さ `n` の列、クエリ数 `q` に対する `Mo's algorithm` です。
    :math:`O(\\frac{n}{\\sqrt{q}})` です。

    Args:
      n (int): 列の長さです。
      q (int): クエリの数です。

    制約:
      :math:`0 \\leq n, 0 \\leq q`
    """

    def __init__(self, n: int, q: int):
        assert 0 <= n and 0 <= q, f"ValueError: {n=} {q=}"
        self.n = n
        self.q = q
        self.bucket_size = ceil(sqrt(3) * n / sqrt(2 * q)) if q > 0 else n
        if self.bucket_size == 0:
            self.bucket_size = 1
        self.bit = max(n, q).bit_length()
        self.msk = (1 << self.bit) - 1
        self.bucket = [[] for _ in range(n // self.bucket_size + 1)]
        self.cnt = 0

    def add_query(self, l: int, r: int) -> None:
        """区間 ``[l, r)`` に対するクエリを追加します。
        :math:`O(1)` です。

        Args:
          l (int):
          r (int):

        制約:
          :math:`0 \\leq l \\leq r \\leq n`
        """
        assert (
            0 <= l <= r <= self.n
        ), f"IndexError: {self.__class__.__name__}.add_query({l}, {r}), self.n={self.n}"
        self.bucket[l // self.bucket_size].append((((r << self.bit) | l) << self.bit) | self.cnt)
        self.cnt += 1

    def run(
        self,
        add: Callable[[int], None],
        delete: Callable[[int], None],
        out: Callable[[int], None],
    ) -> None:
        """クエリを実行します。
        :math:`O(q\\sqrt{n})` です。

        Args:
          add (Callable[[int], None]): 引数のインデックスに対応する要素を追加します。
          delete (Callable[[int], None]): 引数のインデックスに対応する要素を削除します。
          out (Callable[[int], None]): クエリ番号に対する答えを処理します。

        制約:
          ``q`` 回のクエリを ``add_query`` メソッドで追加する必要があります。
        """
        assert self.cnt == self.q, f"Not Enough Queries, now:{self.cnt}, expected:{self.q}"
        bucket, bit, msk = self.bucket, self.bit, self.msk
        for i, b in enumerate(bucket):
            b.sort(reverse=i & 1)
        nl, nr = 0, 0
        for rli in chain(*bucket):
            r, l = rli >> bit >> bit, rli >> bit & msk
            while nl > l:
                nl -= 1
                add(nl)
            while nr < r:
                add(nr)
                nr += 1
            while nl < l:
                delete(nl)
                nl += 1
            while nr > r:
                nr -= 1
                delete(nr)
            out(rli & msk)

    def runrun(
        self,
        add_left: Callable[[int], None],
        add_right: Callable[[int], None],
        delete_left: Callable[[int], None],
        delete_right: Callable[[int], None],
        out: Callable[[int], None],
    ) -> None:
        """クエリを実行します。

        :math:`O(q\\sqrt{n})` です。

        Args:
          add_left (Callable[[int], None]): 引数のインデックスに対応する要素を左から追加します。
          add_right (Callable[[int], None]): 引数のインデックスに対応する要素を右から追加します。
          delete_left (Callable[[int], None]): 引数のインデックスに対応する要素を左から削除します。
          delete_right (Callable[[int], None]): 引数のインデックスに対応する要素を右から削除します。
          out (Callable[[int], None]): クエリ番号に対する答えを処理します。

        制約:
          ``q`` 回のクエリを ``add_query`` メソッドで追加する必要があります。
        """
        assert self.cnt == self.q, f"Not Enough Queries, now:{self.cnt}, expected:{self.q}"
        bucket, bit, msk = self.bucket, self.bit, self.msk
        for i, b in enumerate(bucket):
            b.sort(reverse=i & 1)
        nl, nr = 0, 0
        for rli in chain(*bucket):
            r, l = rli >> bit >> bit, rli >> bit & msk
            while nl > l:
                nl -= 1
                add_left(nl)
            while nr < r:
                add_right(nr)
                nr += 1
            while nl < l:
                delete_left(nl)
                nl += 1
            while nr > r:
                nr -= 1
                delete_right(nr)
            out(rli & msk)
