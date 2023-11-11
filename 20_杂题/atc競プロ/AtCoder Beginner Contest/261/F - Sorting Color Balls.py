# 给定一个数组,有两个属性:颜色c和大小x
# 我们可以进行以下操作来使得数组的大小非递减:
# 对于颜色不同的相邻两项,交换,并产生1代价
# 问:将数组变成非递减时最小的代价是多少?
# n<=3e5

# !变为不减数组的邻位交换次数为冒泡排序交换次数，也为逆序对个数
# 考虑到有颜色的影响，要求的是就是颜色不相同的逆序对的数量, 但是该如何维护呢?
# !颜色不同的逆序对的数量 = 所有逆序对的数量 - 颜色相同的逆序对的数量.
from collections import defaultdict
import sys
import os
from typing import List

from sortedcontainers import SortedList

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def countInv(nums: List[int]) -> int:
    """求逆序对数量之和"""
    n = len(nums)
    res = 0
    visited = SortedList()
    for i in range(n - 1, -1, -1):
        smaller = visited.bisect_left(nums[i])
        res += smaller
        visited.add(nums[i])

    return res


def main() -> None:
    n = int(input())
    colors = list(map(int, input().split()))
    nums = list(map(int, input().split()))
    adjMap = defaultdict(list)
    for i, c in enumerate(colors):
        adjMap[c].append(nums[i])

    inv1 = countInv(nums)
    inv2 = sum(countInv(arr) for arr in adjMap.values())
    print(inv1 - inv2)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
