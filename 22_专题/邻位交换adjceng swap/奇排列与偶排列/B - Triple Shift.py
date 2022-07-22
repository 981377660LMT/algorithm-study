# 数组三个相邻的数可以循环移位
# 问 A 能否变成 B
# n<=5000

# !循环移位:偶排列
# !看两个排列的奇偶性是否一样 (统计逆序数对)
# 对一个数列，如果总的逆序数为奇数，则此排列为奇排列，否则为偶排列。

# 判断所有数字出现个数相不相同，如果不相同一定不可以
# 判断是否有两个或者两个以上相同的数，有的话一定可以
# 判断逆序对数奇偶性是不是相同，相同可以，不相同不可以


import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)

from typing import List
from collections import Counter, defaultdict


# 1 <= nums.length <= 105
# -104 <= nums[i] <= 104


def countSmaller(nums: List[int]) -> List[int]:
    """求逆序对数量"""
    OFFSET = int(1e4) + 10
    res = []
    bit = BIT(3 * OFFSET)
    for i in range(len(nums) - 1, -1, -1):
        cur = bit.query(0, nums[i] - 1 + OFFSET)
        res.append(cur)
        bit.add(nums[i] + OFFSET, nums[i] + OFFSET, 1)
    return list(reversed(res))


class BIT:
    def __init__(self, n: int):
        self.size = n
        self._tree1 = defaultdict(int)
        self._tree2 = defaultdict(int)

    def add(self, left: int, right: int, delta: int) -> None:
        """闭区间[left, right]加delta"""
        self._add(left, delta)
        self._add(right + 1, -delta)

    def query(self, left: int, right: int) -> int:
        """闭区间[left, right]的和"""
        return self._query(right) - self._query(left - 1)

    def _add(self, index: int, delta: int) -> None:
        if index <= 0:
            raise ValueError("index 必须是正整数")

        rawIndex = index
        while index <= self.size:
            self._tree1[index] += delta
            self._tree2[index] += (rawIndex - 1) * delta
            index += index & -index

    def _query(self, index: int) -> int:
        if index > self.size:
            index = self.size

        rawIndex = index
        res = 0
        while index > 0:
            res += rawIndex * self._tree1[index] - self._tree2[index]
            index -= index & -index
        return res


def solve(nums1: List[int], nums2: List[int]) -> bool:
    """nums1能否rotate成nums2"""
    if Counter(nums1) != Counter(nums2):
        return False
    if len(set(nums1)) != len(nums1):
        return True
    inv1, inv2 = sum(countSmaller(nums1)), sum(countSmaller(nums2))
    return inv1 % 2 == inv2 % 2


def main() -> None:
    n = int(input())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))
    if solve(nums1, nums2):
        print("Yes")
    else:
        print("No")


if __name__ == "__main__":

    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
