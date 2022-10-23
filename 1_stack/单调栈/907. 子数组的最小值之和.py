from typing import List
from 每个元素作为最值的影响范围 import getRange

MOD = int(1e9 + 7)

# 找到 min(b) 的总和，其中 b 的范围为 arr 的每个（连续）子数组。
# 需要维护区间极值:单调栈
# https://leetcode-cn.com/problems/sum-of-subarray-minimums/solution/python3-tong-84ti-zui-da-zhi-fang-tu-by-5ersw/

# 思路：考虑每个极小值在多少个子数组里产生贡献
# !注意这里不能计算重复元素的影响范围，因此一边开一边闭


class Solution:
    def sumSubarrayMins(self, arr: List[int]) -> int:
        """求所有子数组的最小值之和"""
        minRange = getRange(arr, isLeftStrict=True, isRightStrict=False)
        res = 0
        for i, num in enumerate(arr):
            left, right = minRange[i]
            count = (right - i + 1) * (i - left + 1)  # 出现在了多少个子数组里
            res = (res + num * count) % MOD
        return res


print(Solution().sumSubarrayMins(arr=[3, 1, 2, 4]))
# 解释：
# 子数组为 [3]，[1]，[2]，[4]，[3,1]，[1,2]，[2,4]，[3,1,2]，[1,2,4]，[3,1,2,4]。
# 最小值为 3，1，2，4，1，1，2，1，1，1，和为 17。
