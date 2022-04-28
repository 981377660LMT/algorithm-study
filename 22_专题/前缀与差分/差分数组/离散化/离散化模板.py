from typing import Iterable
from bisect import bisect_left, bisect_right


class Discretizer:
    """离散化"""

    def __init__(self, nums: Iterable[int]) -> None:
        self.allNums = sorted(set(nums))
        self.mapping = {self.allNums[i]: i + 1 for i in range(len(self.allNums))}

    def get(self, num: int) -> int:
        if num not in self.mapping:
            raise ValueError(f'{num} not in discretizer')
        return self.mapping[num]

    def bisectLeft(self, left: int) -> int:
        """离散化后的左端点

        >>> d = Discretizer([1, 3, 5, 8, 9])
        >>> d.bisectLeft(2)
        2
        >>> d.bisectLeft(10)
        Traceback (most recent call last):
          ...
        ValueError: 10 is bigger than max value in discretizer
        """
        pos = bisect_left(self.allNums, left)
        if pos == len(self.allNums):
            raise ValueError(f'{left} is bigger than max value in discretizer')
        return self.mapping[self.allNums[pos]]

    def bisectRight(self, right: int) -> int:
        """离散化后的右端点

        >>> d = Discretizer([1, 3, 5, 8, 9])
        >>> d.bisectRight(4)
        2
        >>> d.bisectRight(0)
        Traceback (most recent call last):
          ...
        ValueError: 0 is smaller than min value in discretizer
        """
        pos = bisect_right(self.allNums, right) - 1
        if pos < 0:
            raise ValueError(f'{right} is smaller than min value in discretizer')
        return self.mapping[self.allNums[pos]]

    def __len__(self) -> int:
        return len(self.allNums)


if __name__ == '__main__':
    D = Discretizer([1, 3, 5, 8, 9])
    assert D.get(5) == 3
    assert D.bisectLeft(0) == 1
    assert D.bisectLeft(3) == 2
    assert D.bisectLeft(4) == 3
    assert D.bisectRight(4) == 2
    assert D.bisectRight(5) == 3
    assert D.bisectRight(10) == 5

