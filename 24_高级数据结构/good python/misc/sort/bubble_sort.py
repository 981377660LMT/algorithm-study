import random
from typing import Iterable, Protocol, TypeVar, Callable, List


class SupportsLessThan(Protocol):
    def __lt__(self, other) -> bool:
        ...


T = TypeVar("T", bound=SupportsLessThan)


def bubble_sort(
    a: List[T], key: Callable[[T, T], bool] = lambda s, t: s < t, inplace: bool = True
) -> List[T]:
    """バブルソートです。

    非破壊的です。
    最悪 :math:`O(n^2)` 時間です。

    Args:
      a (Iterable[T]): ソートする列です。
      key (Callable[[T, T], bool], optional): 比較関数 `key` にしたがって比較演算をします。
                                              (第1引数)<(第2引数) のとき、 ``True`` を返すようにしてください。
    """
    a = a[:] if inplace else a
    n = len(a)
    for i in range(n):
        flag = True
        for j in range(n - 1, i - 1, -1):
            if not key(a[j - 1], a[j]):
                a[j], a[j - 1] = a[j - 1], a[j]
                flag = False
        if flag:
            break
    return a


def merge_sort(a: Iterable[T], key: Callable[[T, T], bool] = lambda s, t: s < t) -> List[T]:
    """マージソートです。

    非破壊的です。
    最悪 :math:`O(n\\log{n})` 時間です。

    Args:
      a (Iterable[T]): ソートする列です。
      key (Callable[[T, T], bool], optional): 比較関数 `key` にしたがって比較演算をします。
                                              (第1引数)<(第2引数) のとき、 ``True`` を返すようにしてください。
    """

    def _sort(a: List[T]) -> List[T]:
        n = len(a)
        if n <= 1:
            return a
        if n == 2:
            if not key(a[0], a[1]):
                a[0], a[1] = a[1], a[0]
            return a
        left = _sort(a[: n // 2])
        right = _sort(a[n // 2 :])
        res = newlist_hint(n)
        i, j, l, r = 0, 0, len(left), len(right)
        while i < l and j < r:
            if key(left[i], right[j]):
                res.append(left[i])
                i += 1
            else:
                res.append(right[j])
                j += 1
        for i in range(i, l):
            res.append(left[i])
        for j in range(j, r):
            res.append(right[j])
        return res

    return _sort(list(a))


def quick_sort(a: Iterable[T], key: Callable[[T, T], bool] = lambda s, t: s < t) -> List[T]:
    """クイックソートです。

    非破壊的です。
    期待 :math:`O(n\\log{n})` 時間です。

    Args:
      a (Iterable[T]): ソートする列です。
      key (Callable[[T, T], bool], optional): 比較関数 `key` にしたがって比較演算をします。
                                              (第1引数)<(第2引数) のとき、 ``True`` を返すようにしてください。
    """
    a = list(a)

    def sort(i: int, j: int):
        if i >= j:
            return
        pivot = a[random.randint(i, j)]
        l, r = i, j
        while True:
            while key(a[l], pivot):
                l += 1
            while key(pivot, a[r]):
                r -= 1
            if l >= r:
                break
            a[l], a[r] = a[r], a[l]
            l += 1
            r -= 1
        sort(i, l - 1)
        sort(r + 1, j)

    sort(0, len(a) - 1)
    return a
