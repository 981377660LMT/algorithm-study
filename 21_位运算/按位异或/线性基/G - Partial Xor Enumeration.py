# G - Partial Xor Enumeration
# 给你一个长度为 n 的数组，
# 将它所有子序列（可能为空）进行 XOR 操作后得到的值放进一个有序的数组 S 。
# 求数组 S 的 L 到 R 项分别为多少。
# R-L<=2e5 n<=2e5 nums[i]<=60

from typing import List
from LinearBase import LinearBase


def partialXorEnumeration(nums: List[int], left: int, right: int) -> List[int]:
    lb = LinearBase.fromlist(nums)
    res = []
    for i in range(left, right + 1):
        res.append(lb.kthXor(i))
    return res


n, left, right = map(int, input().split())
nums = list(map(int, input().split()))
res = partialXorEnumeration(nums, left, right)
print(*res)
