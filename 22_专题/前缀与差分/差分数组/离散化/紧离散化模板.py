from bisect import bisect_right


class Discretizer:
    """离散化模板"""

    __slots__ = ("_s", "_mp", "_allNums")

    def __init__(self) -> None:
        self._s = set()
        self._mp = dict()

    def add(self, num: int) -> None:
        self._s.add(num)

    def build(self) -> None:
        self._allNums = sorted(self._s)
        for i, num in enumerate(self._allNums):
            self._mp[num] = i + 1

    def get(self, num: int) -> int:
        if num in self._mp:
            return self._mp[num]
        return bisect_right(self._allNums, num)

    def __len__(self) -> int:
        return len(self._allNums)


if __name__ == "__main__":
    D = Discretizer()
    for num in [3, 1, 8, 9, 5]:
        D.add(num)
    D.build()
    assert D.get(5) == 3
    assert D.get(100) == len(D)
    assert D.get(-1) == 0
