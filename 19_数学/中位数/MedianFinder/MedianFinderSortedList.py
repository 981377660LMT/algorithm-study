from typing import Optional

from sortedcontainers import SortedList


class MedianFinderSortedList:
    __slots__ = "sl"

    def __init__(self, sl: Optional[SortedList] = None):
        self.sl = SortedList() if sl is None else sl

    def add(self, num: int) -> None:
        self.sl.add(num)

    def discard(self, num: int) -> None:
        self.sl.discard(num)

    def median(self) -> int:
        """返回向下取整的中位数."""
        len_ = len(self.sl)
        if not len_:
            raise ValueError("No median for empty list")
        if len_ & 1:
            return self.sl[len_ >> 1]
        else:
            mid = len_ >> 1
            return (self.sl[mid - 1] + self.sl[mid]) >> 1

    def __len__(self):
        return len(self.sl)
