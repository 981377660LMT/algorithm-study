from bisect import bisect_left, bisect_right
from typing import List


def maxSumOfMountainSubsequence(
    nums: List[int], leftStrict=True, rightStrict=True, allowEmptySide=False
) -> int:
    """最大山脉子序列和."""
    n = len(nums)
    rev = nums[::-1]
    preMax = [0] + LISMaxSum(nums, leftStrict)
    sufMax = LISMaxSum(rev, rightStrict)[::-1] + [0]
    res = 0
    if allowEmptySide:
        for i in range(n):
            res = max(res, preMax[i + 1] + sufMax[i] - nums[i])
    else:
        leftLen = caldp(nums, leftStrict)
        rightLen = caldp(rev, rightStrict)[::-1]
        for i in range(n):
            if leftLen[i] < 2 or rightLen[i] < 2:
                continue
            res = max(res, preMax[i + 1] + sufMax[i] - nums[i])
    return res


def LISMaxSum(nums: List[int], strict=True) -> List[int]:
    """求以每个位置为结尾的LIS最大和(包括自身)."""

    def max(a: int, b: int) -> int:
        return a if a > b else b

    class BITPrefixMax:
        __slots__ = ("_max", "_tree")

        def __init__(self, max: int):
            self._max = max
            self._tree = dict()

        def set(self, index: int, value: int) -> None:
            index += 1
            while index <= self._max:
                self._tree[index] = max(self._tree.get(index, 0), value)
                index += index & -index

        def query(self, end: int) -> int:
            """Query max of [0, end)."""
            if end > self._max:
                end = self._max
            res = 0
            while end > 0:
                res = max(res, self._tree.get(end, 0))
                end -= end & -end
            return res

    n = len(nums)
    if n <= 1:
        return nums[:]
    max_ = 0
    for v in nums:
        max_ = max(max_, v)
    dp = BITPrefixMax(max_ + 5)
    res = [0] * n
    for i, v in enumerate(nums):
        preMax = dp.query(v) if strict else dp.query(v + 1)
        cur = preMax + v
        res[i] = cur
        dp.set(v, cur)
    return res


def caldp(nums: List[int], strict=True) -> List[int]:
    """求以每个位置为结尾的LIS长度(包括自身)"""
    if not nums:
        return []
    n = len(nums)
    res = [0] * n
    lis = []
    f = bisect_left if strict else bisect_right
    for i in range(n):
        pos = f(lis, nums[i])
        if pos == len(lis):
            lis.append(nums[i])
            res[i] = len(lis)
        else:
            lis[pos] = nums[i]
            res[i] = pos + 1
    return res


if __name__ == "__main__":
    nums = [2, 1, 4, 7, 3, 2, 5, 100]
    print(maxSumOfMountainSubsequence(nums, allowEmptySide=True))
    print(maxSumOfMountainSubsequence(nums, allowEmptySide=False))

    def check(nums: List[int], leftStrict=True, rightStrict=True, allowEmptySide=False):
        def validMountainArray(
            arr: List[int], leftStrict=True, rightStrict=True, allowEmptySide=False
        ) -> bool:
            """有效的山脉数组."""
            n = len(arr)
            ptr = 0

            if leftStrict:
                while ptr + 1 < n and arr[ptr] < arr[ptr + 1]:
                    ptr += 1
            else:
                while ptr + 1 < n and arr[ptr] <= arr[ptr + 1]:
                    ptr += 1

            if not allowEmptySide and (ptr == 0 or ptr == n - 1):
                return False

            if rightStrict:
                while ptr + 1 < n and arr[ptr] > arr[ptr + 1]:
                    ptr += 1
            else:
                while ptr + 1 < n and arr[ptr] >= arr[ptr + 1]:
                    ptr += 1

            return ptr == n - 1

        res = 0
        for state in range(1 << len(nums)):
            if state == 0:
                continue
            selected = [v for i, v in enumerate(nums) if state & (1 << i)]
            if validMountainArray(selected, leftStrict, rightStrict, allowEmptySide):
                res = max(res, sum(selected))
        return res

    from random import randint
    from itertools import product

    for n, ls, rs, ae in product(range(1, 15), [True, False], [True, False], [True, False]):
        nums = [randint(0, 100) for _ in range(n)]
        # nums = [22, 17, 66, 13, 56]
        # nums = [1, 2, 3, 4, 5]
        # nums = [2, 2, 2, 2, 2]
        if maxSumOfMountainSubsequence(nums, ls, rs, ae) != check(nums, ls, rs, ae):
            print(nums, ls, rs, ae)
            print(maxSumOfMountainSubsequence(nums, ls, rs, ae))
            print(check(nums, ls, rs, ae))
            break
    else:
        print("pass")
