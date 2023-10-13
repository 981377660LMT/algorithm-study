# erase, prev, next 的时间复杂度均为 O(1)
# erase: 删除某个元素
# prev: 找到前一个元素
# next: 找到后一个元素


class FinderLinkedList:
    """
    使用双向链表维护前驱后继.
    初始时, 0~n-1 个元素都是未访问过的.
    """

    __slots__ = "_n", "_prev", "_next", "_erased"

    def __init__(self, n: int) -> None:
        """初始化元素0~n-1."""
        self._n = n
        self._prev = [i - 1 for i in range(n + 2)]
        self._next = [i + 1 for i in range(n + 2)]
        self._erased = [False for _ in range(n)]

    def erase(self, i: int) -> bool:
        """删除元素i."""
        if self._erased[i]:
            return False
        self._erased[i] = True
        i += 1
        self._prev[self._next[i]] = self._prev[i]
        self._next[self._prev[i]] = self._next[i]
        return True

    def has(self, i: int) -> bool:
        """判断元素i是否存在."""
        return not self._erased[i]

    def prev(self, i: int) -> int:
        """
        找到i左侧第一个未被访问过的位置(包含i).
        如果不存在, 返回-1.
        """
        if self.has(i):
            return i
        res = self._prev[i + 1] - 1
        return res if res >= 0 else -1

    def next(self, i: int) -> int:
        """
        找到i右侧第一个未被访问过的位置(包含i).
        如果不存在, 返回-1.
        """
        if self.has(i):
            return i
        res = self._next[i + 1] - 1
        return res if res < self._n else -1


if __name__ == "__main__":
    finder = FinderLinkedList(5)

    for i in range(5):
        print(finder.prev(i))
        print(finder.next(i))

    finder.erase(2)
    for i in range(5):
        print(finder.prev(i))
        print(finder.next(i))
