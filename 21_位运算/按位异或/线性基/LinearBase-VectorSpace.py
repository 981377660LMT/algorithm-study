# 可合并的线性基/线性基合并


from typing import List, Optional


class VectorSpace:
    @staticmethod
    def merge(v1: "VectorSpace", v2: "VectorSpace") -> "VectorSpace":
        """merge."""
        if len(v1) < len(v2):
            v1, v2 = v2, v1
        res = v1.copy()
        for v in v2.bases:
            res.add(v)
        return res

    @staticmethod
    def mergeDestructively(v1: "VectorSpace", v2: "VectorSpace") -> "VectorSpace":
        """merge."""
        if len(v1) < len(v2):
            v1, v2 = v2, v1
        res = v1
        for v in v2.bases:
            res.add(v)
        return res

    @staticmethod
    def f2Intersection(nums1: List[int], nums2: List[int], maxLog=32) -> List[int]:
        tmp = VectorSpace()
        for a in nums1:
            tmp.add((a << maxLog) + a)
        upper = 1 << maxLog
        res = []
        for b in nums2:
            v = b << maxLog
            u = tmp.add2(v)
            if u < upper:
                res.append(u)
        return res

    __slots__ = "bases"

    def __init__(self, nums: Optional[List[int]] = None) -> None:
        self.bases = []
        if nums is not None:
            for v in nums:
                self.add(v)

    def add(self, v: int) -> bool:
        """插入一个向量,如果插入成功(不能被表出)返回True,否则返回False."""
        for e in self.bases:
            if e == 0 or v == 0:
                break
            v = min2(v, v ^ e)
        if v:
            self.bases.append(v)
            return True
        return False

    def add2(self, v: int) -> int:
        """插入一个向量,返回表出的向量."""
        for e in self.bases:
            if e == 0 or v == 0:
                break
            v = min2(v, v ^ e)
        if v:
            self.bases.append(v)
        return v

    def getMax(self, xorVal=0) -> int:
        """求xorVal与所有向量异或的最大值."""
        res = xorVal
        for e in self.bases:
            res = max2(res, res ^ e)
        return res

    def getMin(self, xorVal=0) -> int:
        """求xorVal与所有向量异或的最小值."""
        res = xorVal
        for e in self.bases:
            res = min2(res, res ^ e)
        return res

    def copy(self) -> "VectorSpace":
        res = VectorSpace()
        res.bases = self.bases[:]
        return res

    def __len__(self) -> int:
        return len(self.bases)

    def __iter__(self):
        return iter(self.bases)

    def __repr__(self) -> str:
        return repr(self.bases)

    def __contains__(self, v: int) -> bool:
        for e in self.bases:
            if v == 0:
                break
            v = min2(v, v ^ e)
        return v == 0

    def __or__(self, other: "VectorSpace") -> "VectorSpace":
        """merge."""
        v1, v2 = self, other
        if len(v1) < len(v2):
            v1, v2 = v2, v1
        res = v1.copy()
        for v in v2.bases:
            res.add(v)
        return res

    def __ior__(self, other: "VectorSpace") -> "VectorSpace":
        """merge."""
        for v in other.bases:
            self.add(v)
        return self


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


if __name__ == "__main__":
    # a = VectorSpace([1, 2, 3])
    # b = VectorSpace([6])

    # # # # https://atcoder.jp/contests/abc141/tasks/abc141_f
    # # # # !把一个数组分成两个非空子集, 使得两个集合的异或和之和最大

    # n = int(input())
    # nums = list(map(int, input().split()))

    # xor_ = 0
    # V1 = VectorSpace()
    # for v in nums:
    #     xor_ ^= v
    #     V1.add(v)

    # mask = ~xor_
    # V2 = VectorSpace()
    # for v in V1.bases:
    #     V2.add(v & mask)

    # res = V2.getMax()
    # print(res + (xor_ ^ res))

    # https://judge.yosupo.jp/problem/intersection_of_f2_vector_spaces

    def intersection_of_f2_vector_spaces():
        nums1 = list(map(int, input().split()))[1:]
        nums2 = list(map(int, input().split()))[1:]

        res = VectorSpace.f2Intersection(nums1, nums2)
        print(len(res), end=" ")
        print(*res)

    T = int(input())
    for _ in range(T):
        intersection_of_f2_vector_spaces()
