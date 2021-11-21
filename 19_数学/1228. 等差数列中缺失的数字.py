# 我们会从该数组中删除一个 既不是第一个 也 不是最后一个的值，得到一个新的数组  arr。

# 给你这个缺值的数组 arr，请你帮忙找出被删除的那个数。

# 等差数列和=（首项+尾项）* 项数 ÷ 2
# 项数为当前arr中元素个数+1
# “数列和”与“当前数组的和” 的差值就是要求的元素
from typing import List


class Solution:
    def missingNumber(self, arr: List[int]) -> int:
        return (arr[0] + arr[-1]) * (len(arr) + 1) // 2 - sum(arr)

