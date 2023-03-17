from collections import defaultdict
from typing import List


class LinearBase:
    __slots__ = ("bases", "_rows", "_bit")

    @staticmethod
    def fromlist(nums: List[int]) -> "LinearBase":
        res = LinearBase(bit=max(nums, default=0).bit_length())
        for x in nums:
            res.add(x)
        res.build()
        return res

    def __init__(self, bit=62):
        self.bases = []  # 基底
        self._rows = defaultdict(int)  # 高斯消元的行
        self._bit = bit  # 最大数的位数

    def add(self, x: int) -> bool:
        """插入一个向量,如果插入成功返回True,否则返回False"""
        x = self._normalize(x)
        if x == 0:
            return False
        i = x.bit_length() - 1
        for j in range(self._bit):
            if (self._rows[j] >> i) & 1:
                self._rows[j] ^= x
        self._rows[i] = x
        return True

    def build(self) -> None:
        res = []
        for _, v in sorted(self._rows.items()):
            if v > 0:
                res.append(v)
        self.bases = res

    def kthXor(self, k: int) -> int:
        """子序列(子集,包含空集)第k小的异或 1<=k<=2**len(self.bases)"""
        assert 1 <= k <= 2 ** len(self.bases)
        k -= 1
        res = 0
        for i in range(k.bit_length()):
            if (k >> i) & 1:
                res ^= self.bases[i]
        return res

    def maxXor(self) -> int:
        return self.kthXor(2 ** len(self.bases))

    def _normalize(self, x: int) -> int:
        for i in range(x.bit_length() - 1, -1, -1):
            if (x >> i) & 1:
                x ^= self._rows[i]
        return x

    def __len__(self) -> int:
        return len(self.bases)

    def __contains__(self, x: int) -> bool:
        """x是否能由线性基表出"""
        return self._normalize(x) == 0


if __name__ == "__main__":
    nums = [1, 2, 3, 4, 5, 6, 7, 8, 9, 999]
    lb = LinearBase.fromlist(nums)
    print(lb.bases, len(lb))
    print(lb.maxXor())
    print(lb.kthXor(2))
    print(lb.kthXor(17))

    # test __contains__
    res = set()
    for i in range(1 << len(lb)):
        bases = []
        for j in range(len(lb)):
            if (i >> j) & 1:
                bases.append(lb.bases[j])
        cur = 0
        for b in bases:
            cur ^= b
        res.add(cur)
    res = sorted(res)
    ok = [i for i in range(lb.maxXor() + 1) if i in lb]
    assert res == ok
