# 现给出一个正整数 K，要求经过 K 次数组相邻位置元素交换（必须完成 K 次交换），使得这个数组代表的数字最大。


from bisect import bisect_left, bisect_right, insort_left
from collections import defaultdict, deque
from typing import Any, Generic, Iterable, List, Optional, Protocol, TypeVar, Union


class SupportsDunderLT(Protocol):
    def __lt__(self, __other: Any) -> bool:
        ...


class SupportsDunderGT(Protocol):
    def __gt__(self, __other: Any) -> bool:
        ...


S = TypeVar('S', bound=Union[SupportsDunderLT, SupportsDunderGT])


class SortedList(Generic[S]):
    """用bisect模拟 插入和删除的时候用切片"""

    def __init__(self, iterable: Optional[Iterable[S]] = None) -> None:
        self._list = []
        if iterable is not None:
            for item in iterable:
                self.add(item)

    def add(self, item: S) -> None:
        insort_left(self._list, item)

    def pop(self, index: int) -> S:
        res = self._list[index]
        self._list = self._list[:index] + self._list[index + 1 :]
        return res

    def bisect_left(self, item: S) -> int:
        return bisect_left(self._list, item)

    def bisect_right(self, item: S) -> int:
        return bisect_right(self._list, item)

    def __getitem__(self, index: int) -> S:
        return self._list[index]

    def __len__(self) -> int:
        return len(self._list)


def maxInteger(nums: List[int], k: int) -> List[int]:
    """最多 K 次交换相邻数位后得到的最大整数"""
    n = len(nums)
    indexMap = defaultdict(deque)
    for index, char in enumerate(nums):
        indexMap[char].append(index)

    res = []
    sl = SortedList()

    for i in range(n):
        for digit in range(9, -1, -1):
            indexes = indexMap[digit]
            if not indexes:
                continue

            pos = indexes[0]
            dist = pos - i  # 距离
            preCount = len(sl) - sl.bisect_right(pos)  # 原来在pos右边的数,换到前面了,看右边多少个数比他大
            cost = dist + preCount
            if cost <= k:
                k -= cost
                res.append(digit)
                sl.add(indexes.popleft())
                break

    return res


T = int(input())
for _ in range(T):
    k = int(input())
    n = int(input())
    nums = list(map(int, input().split()))
    res = maxInteger(nums, k)
    print(' '.join(map(str, res)))

