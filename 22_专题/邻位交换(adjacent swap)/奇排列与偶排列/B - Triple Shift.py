# https://blog.csdn.net/qq_51152918/article/details/123265469
# https://atcoder.jp/contests/arc136/editorial/3467

# 数组三个相邻的数可以循环移位
# 问 A 能否变成 B
# n<=5000
# !循环移位:偶排列
# !看两个排列的奇偶性是否一样 (统计逆序数对)
# 对一个数列，如果总的逆序数为奇数，则此排列为奇排列，否则为偶排列。

# 1.判断所有数字出现个数相不相同，如果不相同一定不可以
# 2.判断是否有两个或者两个以上相同的数，有的话一定可以(可以利用该种数字来改变数组逆序对个数的奇偶性)
# 3.判断逆序对数奇偶性是不是相同，相同可以，不相同不可以


from sortedcontainers import SortedList


from typing import List
from collections import Counter


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


def solve(A: List[int], B: List[int]) -> bool:
    """A能否rotate成B"""
    C1 = Counter(A)
    C2 = Counter(B)
    if C1 != C2:
        return False
    if max(C1.values(), default=0) > 1:
        return True
    inv1, inv2 = sum(countSmaller2(A)), sum(countSmaller2(B))
    return inv1 % 2 == inv2 % 2


if __name__ == "__main__":
    # !每次选三个相连的数字进行逆时针旋转,问能否使得两个数组相等
    n = int(input())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))
    if solve(nums1, nums2):
        print("Yes")
    else:
        print("No")
