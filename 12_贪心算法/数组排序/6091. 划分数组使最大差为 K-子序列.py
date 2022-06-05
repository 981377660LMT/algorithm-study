from typing import List, Tuple


MOD = int(1e9 + 7)
INF = int(1e20)

# 在满足每个子序列中最大值和最小值之间的差值最多为 k 的前提下，返回需要划分的 最少 子序列数目。
# !子序列可以排序
class Solution:
    def partitionArray(self, nums: List[int], k: int) -> int:
        nums.sort()
        res = 1
        pre = nums[0]
        for num in nums:
            if num - pre > k:
                res += 1
                pre = num
        return res

