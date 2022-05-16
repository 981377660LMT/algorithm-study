# https://codeforces.com/blog/entry/91768

# 平均数不小于k的子数组的个数  => dp 公式 变形
# n<=1e5
# nums[i]<=1e4

# from sortedcontainers import SortedList
from itertools import accumulate
from bisect import bisect_left, bisect_right, insort_left
from typing import Any, Generic, Iterable, Optional, Protocol, TypeVar, Union


n, k = map(int, input().split())
nums = list(map(int, input().split()))


# 1 3 4
# preSum[i] - preSum[j] >= k * (i - j)

# 即 preSum[i]-k*i >= preSum[j] - k*j


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
        return self._list.pop(index)

    def bisect_left(self, item: S) -> int:
        return bisect_left(self._list, item)

    def bisect_right(self, item: S) -> int:
        return bisect_right(self._list, item)

    def __getitem__(self, index: int) -> S:
        return self._list[index]

    def __len__(self) -> int:
        return len(self._list)


preSum = [0] + list(accumulate(nums))
sortedList = SortedList()
res = 0
for i in range(n + 1):
    cur = preSum[i] - k * i
    pos = sortedList.bisect_right(cur)
    res += pos
    sortedList.add(cur)
print(res)

######################################

