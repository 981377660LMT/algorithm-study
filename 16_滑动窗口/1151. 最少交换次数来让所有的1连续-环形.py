from typing import List

# 你需要通过交换位置，将数组中的1连续

# 固定区间，假设最后的全1就在当前的窗口
# 则窗口内0的个数==需要交换的次数


class Solution:
    def minSwaps(self, nums: List[int]) -> int:
        winLen = nums.count(1)
        nums = nums * 2

        res = int(1e20)
        curOne = 0
        left = 0
        for right, cur in enumerate(nums):
            if cur == 1:
                curOne += 1
            if right >= winLen:
                if nums[left] == 1:
                    curOne -= 1
                left += 1

            if right >= winLen - 1:
                res = min(res, winLen - curOne)
        return res


print(Solution().minSwaps([1, 0, 1, 0, 1]))
# 输出：1
# 解释：
# 有三种可能的方法可以把所有的 1 组合在一起：
# [1,1,1,0,0]，交换 1 次；
# [0,1,1,1,0]，交换 2 次；
# [0,0,1,1,1]，交换 1 次。
# 所以最少的交换次数为 1。
