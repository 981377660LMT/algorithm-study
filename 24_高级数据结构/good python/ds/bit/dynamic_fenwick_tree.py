# from titan_pylib.data_structures.fenwick_tree.dynamic_fenwick_tree import DynamicFenwickTree
from typing import Optional, Final, Dict


class DynamicFenwickTree:
    """必要なところだけノードを作ります。"""

    def __init__(self, u: int):
        """Build DynamicFenwickTree [0, u)."""
        assert isinstance(u, int), f"TypeError: DynamicFenwickTree({u}), {u} must be int"
        self._u: Final[int] = u
        self._tree: Dict[int, int] = {}
        self._s: Final[int] = 1 << (u - 1).bit_length()

    def add(self, k: int, x: int) -> None:
        assert 0 <= k < self._u, f"IndexError: DynamicFenwickTree.add({k}, {x}), u={self._u}"
        k += 1
        while k <= self._u:
            if k in self._tree:
                self._tree[k] += x
            else:
                self._tree[k] = x
            k += k & -k

    def pref(self, r: int) -> int:
        assert 0 <= r <= self._u, f"IndexError: DynamicFenwickTree.pref({r}), u={self._u}"
        ret = 0
        while r > 0:
            ret += self._tree.get(r, 0)
            r -= r & -r
        return ret

    def sum(self, l: int, r: int) -> int:
        assert 0 <= l <= r <= self._u, f"IndexError: DynamicFenwickTree.sum({l}, {r}), u={self._u}"
        # return self.pref(r) - self.pref(l)
        _tree = self._tree
        res = 0
        while r > l:
            res += _tree.get(r, 0)
            r -= r & -r
        while l > r:
            res -= _tree.get(l, 0)
            l -= l & -l
        return res

    def bisect_left(self, w: int) -> Optional[int]:
        i, s = 0, self._s
        while s:
            if i + s <= self._u:
                if i + s in self._tree and self._tree[i + s] < w:
                    w -= self._tree[i + s]
                    i += s
                elif i + s not in self._tree and 0 < w:
                    i += s
            s >>= 1
        return i if w else None

    def bisect_right(self, w: int) -> int:
        i, s = 0, self._s
        while s:
            if i + s <= self._u:
                if i + s in self._tree and self._tree[i + s] <= w:
                    w -= self._tree[i + s]
                    i += s
                elif i + s not in self._tree and 0 <= w:
                    i += s
            s >>= 1
        return i

    def __str__(self):
        return str(self._tree)
