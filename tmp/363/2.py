from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始、长度为 n 的整数数组 nums ，其中 n 是班级中学生的总数。班主任希望能够在让所有学生保持开心的情况下选出一组学生：

# 如果能够满足下述两个条件之一，则认为第 i 位学生将会保持开心：


# 这位学生被选中，并且被选中的学生人数 严格大于 nums[i] 。
# 这位学生没有被选中，并且被选中的学生人数 严格小于 nums[i] 。
# 返回能够满足让所有学生保持开心的分组方法的数目。
class Solution:
    def countWays(self, nums: List[int]) -> int:
        nums.sort()
        res = 0
        # 选择前缀并检查
        for select in range(len(nums) + 1):
            if select == 0:
                res += 1 if nums[0] > 0 else 0
            else:
                ptr = select - 1
                if select > nums[ptr] and ((ptr + 1 == len(nums) or (select < nums[ptr + 1]))):
                    res += 1
        return res


nums = [1, 1]

print(Solution().countWays(nums))
