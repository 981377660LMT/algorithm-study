# 数组三个相邻的数可以循环移位
# 问 A 能否变成 B
# n<=5000

# !循环移位:偶排列
# !看两个排列的奇偶性是否一样 (统计逆序数对)
# 对一个数列，如果总的逆序数为奇数，则此排列为奇排列，否则为偶排列。

# 判断所有数字出现个数相不相同，如果不相同一定不可以
# 判断是否有两个或者两个以上相同的数，有的话一定可以
# 判断逆序对数奇偶性是不是相同，相同可以，不相同不可以


from bisect import bisect_right
import sys
import os

from sortedcontainers import SortedList

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)

from typing import List
from collections import Counter, defaultdict


def countInv(nums: List[int]) -> int:
    """求逆序对数"""
    n = len(nums)
    res = 0
    visited = SortedList()
    for i in range(n - 1, -1, -1):
        smaller = visited.bisect_left(nums[i])
        res += smaller
        visited.add(nums[i])

    return res


def countSmaller2(nums: List[int]) -> List[int]:
    """sortedList求每个位置处的逆序对数量"""
    n = len(nums)
    res = [0] * n
    visited = SortedList()
    for i in range(n - 1, -1, -1):
        smaller = visited.bisect_left(nums[i])
        res[i] = smaller
        visited.add(nums[i])

    return res


def countSmaller(nums: List[int]) -> List[int]:
    """求每个位置处的逆序对数量  注意值域很大时需要离散化"""
    n = len(nums)
    arr = sorted(nums)
    res = [0] * n
    bit = BIT1(n + 10)
    for i in range(len(nums) - 1, -1, -1):
        pos1 = bisect_right(arr, nums[i] - 1) + 1
        cur = bit.query(pos1)
        res[i] = cur
        pos2 = bisect_right(arr, nums[i]) + 1
        bit.add(pos2, 1)
    return res


class BIT1:
    def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(int)

    def add(self, index: int, delta: int) -> None:
        if index <= 0:
            raise ValueError("index 必须是正整数")
        while index <= self.size:
            self.tree[index] += delta
            index += index & -index

    def query(self, index: int) -> int:
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res += self.tree[index]
            index -= index & -index
        return res

    def queryRange(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)


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
