# 2386. 找出数组的第 K 大和
# https://leetcode.cn/problems/find-the-k-sum-of-an-array/solutions/1764389/zhuan-huan-dui-by-endlesscheng-8yiq/
# 子集和的第 k 小
# n,k<=1e5
# k从1开始
# 初始时插入 (a[0],0)，然后执行 k-1 次操作：
# 取出堆顶，插入 (top.v+a[top.i+1],top.i+1) 以及 (top.v+a[top.i+1]-a[top.i],top.i+1)

from heapq import heappop, heappush
from typing import List


def kthSubsetSum(nums: List[int], k: int) -> int:
    """第k小的子序列和."""
    nums = nums[:]
    sum_ = 0
    for i, x in enumerate(nums):
        if x >= 0:
            sum_ += x
        else:
            nums[i] = -x
    nums.sort()
    pq = [(-sum_, 0)]
    for _ in range(k - 1):
        s, i = heappop(pq)
        if i < len(nums):
            heappush(pq, (s + nums[i], i + 1))  # type: ignore
            if i:
                heappush(pq, (s + nums[i] - nums[i - 1], i + 1))
    return -pq[0][0]
