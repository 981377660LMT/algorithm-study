# from titan_pylib.data_structures.deque.deque import Deque
from typing import Iterable, List, Generic, TypeVar

T = TypeVar("T")


class Deque(Generic[T]):
    """Deque です。
    ランダムアクセスが :math:`O(1)` で可能です。
    """

    def __init__(self, a: Iterable[T] = []):
        """
        :math:`O(n)` です。

        Args:
          a (Iterable[T], optional): ``Deque`` を構築する配列です。
        """
        self.front: List[T] = []
        self.back: List[T] = list(a)

    def _rebuild(self) -> None:
        new = self.front[::-1] + self.back
        self.front = new[: len(new) // 2][::-1]
        self.back = new[len(new) // 2 :]

    def append(self, v: T) -> None:
        """要素 ``v`` を末尾に追加します。
        :math:`O(1)` です。

        Args:
          v (T): 追加する要素です。
        """
        self.back.append(v)

    def appendleft(self, v: T) -> None:
        """要素 ``v`` を先頭に追加します。
        :math:`O(1)` です。

        Args:
          v (T): 追加する要素です。
        """
        self.front.append(v)

    def pop(self) -> T:
        """末尾の要素を削除し、その値を返します。
        :math:`O(1)` です。
        """
        if not self.back:
            self._rebuild()
        return self.back.pop() if self.back else self.front.pop()

    def popleft(self) -> T:
        """先頭の要素を削除し、その値を返します。
        :math:`O(1)` です。
        """
        if not self.front:
            self._rebuild()
        return self.front.pop() if self.front else self.back.pop()

    def tolist(self) -> List[T]:
        """``list`` に変換します。
        :math:`O(n)` です。
        """
        return self.front[::-1] + self.back

    def __getitem__(self, k: int) -> T:
        """``k`` 番目の要素を取得します。
        :math:`O(1)` です。
        """
        assert (
            -len(self) <= k < len(self)
        ), f"IndexError: {self.__class__.__name__}[{k}], len={len(self)}"
        if k < 0:
            k += len(self)
        return (
            self.front[len(self.front) - k - 1]
            if k < len(self.front)
            else self.back[k - len(self.front)]
        )

    def __setitem__(self, k: int, v: T):
        """``k`` 番目の要素を ``v`` に更新します。
        :math:`O(1)` です。
        """
        assert (
            -len(self) <= k < len(self)
        ), f"IndexError: {self.__class__.__name__}[{k} = {v}, len={len(self)}"
        if k < 0:
            k += len(self)
        if k < len(self.front):
            self.front[len(self.front) - k - 1] = v
        else:
            self.back[k - len(self.front)] = v

    def __bool__(self):
        return bool(self.front or self.back)

    def __len__(self):
        """要素数を取得します。
        :math:`O(1)` です。
        """
        return len(self.front) + len(self.back)

    def __contains__(self, v: T):
        return (v in self.front) or (v in self.back)

    def __str__(self):
        return str(self.tolist())

    def __repr__(self):
        return f"{self.__class__.__name__}({self})"
