from typing import List, Tuple


class LinearBaseRestore:
    __slots__ = "bases"

    def __init__(self, nums: List[int]):
        self.bases = []
        for i, x in enumerate(nums):
            self.add(x, i)

    def add(self, num: int, id: int) -> bool:
        """插入一个向量,如果插入成功返回True,否则返回False."""
        v = [num, {id}]
        for b in self.bases:
            if v[0] > (v[0] ^ b[0]):
                self._apply(v, b)
        if v[0] != 0:
            self.bases.append(v)
            return True
        return False

    def restore(self, x: int) -> Tuple[List[int], bool]:
        """向量表出.返回基底的下标列表和是否成功."""
        v = [x, set()]
        for b in self.bases:
            if v[0] > (v[0] ^ b[0]):
                self._apply(v, b)
        if v[0] != 0:
            return [], False
        ids = sorted(v[1])  # type: ignore
        return ids, True

    def _apply(self, p, o):
        p[0] ^= o[0]
        for x in o[1]:
            self._toggle(p[1], x)

    def _toggle(self, set, x):
        if x in set:
            set.remove(x)
        else:
            set.add(x)


if __name__ == "__main__":
    s = LinearBaseRestore([1, 2, 3, 4, 5])
    print(s.restore(6))
