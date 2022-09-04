from typing import List

# 你需要通过交换位置，将数组中 任何位置 上的 1 组合到一起，并返回所有可能中所需 最少的交换次数。

# 固定区间，假设最后的全1就在当前的窗口
# 则窗口内0的个数==需要交换的次数

INF = int(1e20)


class Solution:
    def minSwaps(self, nums: List[int]) -> int:
        k = nums.count(1)  # 定长滑窗长度
        n = len(nums)
        res, curSum = INF, 0
        for right in range(n):
            curSum += int(nums[right] == 1)
            if right >= k:
                curSum -= nums[right - k]
            if right >= k - 1:
                res = min(res, k - curSum)
        return res


print(Solution().minSwaps([1, 0, 1, 0, 1]))
# 输出：1
# 解释：
# 有三种可能的方法可以把所有的 1 组合在一起：
# [1,1,1,0,0]，交换 1 次；
# [0,1,1,1,0]，交换 2 次；
# [0,0,1,1,1]，交换 1 次。
# 所以最少的交换次数为 1。
