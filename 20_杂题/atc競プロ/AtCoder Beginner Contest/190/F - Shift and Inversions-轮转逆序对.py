# shift 逆序对
# 轮转逆序对
# n <= 3e5
from typing import List

from sortedcontainers import SortedList


def shiftAndInversions(nums: List[int]) -> List[int]:
    sl = SortedList()
    inv = 0
    for num in nums[::-1]:
        inv += sl.bisect_left(num)
        sl.add(num)

    res = []
    for num in nums:
        res.append(inv)
        inv -= sl.bisect_left(num)
        sl.remove(num)
        inv += len(sl) - sl.bisect_right(num)
        sl.add(num)
    return res


if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    print(*shiftAndInversions(nums), sep="\n")
