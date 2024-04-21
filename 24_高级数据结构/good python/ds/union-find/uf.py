# from titan_pylib.data_structures.union_find.union_find import UnionFind
from typing import List
from collections import defaultdict


class UnionFind:
    def __init__(self, n: int) -> None:
        """``n`` 個の要素からなる ``UnionFind`` を構築します。

        :math:`O(n)` です。
        """
        self._n: int = n
        self._group_numbers: int = n
        self._parents: List[int] = [-1] * n

    def root(self, x: int) -> int:
        """要素 ``x`` を含む集合の代表元を返します。

        :math:`O(\\alpha(n))` です。
        """
        a = x
        while self._parents[a] >= 0:
            a = self._parents[a]
        while self._parents[x] >= 0:
            y = x
            x = self._parents[x]
            self._parents[y] = a
        return a

    def unite(self, x: int, y: int) -> bool:
        """要素 ``x`` を含む集合と要素 ``y`` を含む集合を併合します。
        :math:`O(\\alpha(n))` です。

        Returns:
          bool: もともと同じ集合であれば ``False``、そうでなければ ``True`` を返します。
        """
        x = self.root(x)
        y = self.root(y)
        if x == y:
            return False
        self._group_numbers -= 1
        if self._parents[x] > self._parents[y]:
            x, y = y, x
        self._parents[x] += self._parents[y]
        self._parents[y] = x
        return True

    def unite_right(self, x: int, y: int) -> int:
        # x -> y
        x = self.root(x)
        y = self.root(y)
        if x == y:
            return x
        self._group_numbers -= 1
        self._parents[y] += self._parents[x]
        self._parents[x] = y
        return y

    def unite_left(self, x: int, y: int) -> int:
        # x <- y
        x = self.root(x)
        y = self.root(y)
        if x == y:
            return x
        self._group_numbers -= 1
        self._parents[x] += self._parents[y]
        self._parents[y] = x
        return x

    def size(self, x: int) -> int:
        """要素 ``x`` を含む集合の要素数を返します。
        :math:`O(\\alpha(n))` です。
        """
        return -self._parents[self.root(x)]

    def same(self, x: int, y: int) -> bool:
        """
        要素 ``x`` と ``y`` が同じ集合に属するなら ``True`` を、
        そうでないなら ``False`` を返します。
        :math:`O(\\alpha(n))` です。
        """
        return self.root(x) == self.root(y)

    def members(self, x: int) -> List[int]:
        """要素 ``x`` を含む集合を返します。"""
        x = self.root(x)
        return [i for i in range(self._n) if self.root(i) == x]

    def all_roots(self) -> List[int]:
        """全ての集合の代表元からなるリストを返します。
        :math:`O(n)` です。

        Returns:
            List[int]: 昇順であることが保証されます。
        """
        return [i for i, x in enumerate(self._parents) if x < 0]

    def group_count(self) -> int:
        """集合の総数を返します。
        :math:`O(1)` です。
        """
        return self._group_numbers

    def all_group_members(self) -> defaultdict:
        """
        `key` に代表元、 `value` に `key` を代表元とする集合のリストをもつ `defaultdict` を返します。

        :math:`O(n\\alpha(n))` です。
        """
        group_members = defaultdict(list)
        for member in range(self._n):
            group_members[self.root(member)].append(member)
        return group_members

    def clear(self) -> None:
        """集合の連結状態をなくします(初期状態に戻します)。

        :math:`O(n)` です。
        """
        self._group_numbers = self._n
        for i in range(self._n):
            self._parents[i] = -1

    def __str__(self) -> str:
        """よしなにします。

        :math:`O(n\\alpha(n))` です。
        """
        return (
            f"<{self.__class__.__name__}> [\n"
            + "\n".join(f"  {k}: {v}" for k, v in self.all_group_members().items())
            + "\n]"
        )
