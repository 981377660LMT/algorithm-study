# 输入：nums = [0,1,2,4,5,7]
# 输出：["0->2","4->5","7"]
# 解释：区间范围是：
# [0,2] --> "0->2"
# [4,5] --> "4->5"
# [7,7] --> "7"
# 列表中的每个区间范围 [a,b] 应该按如下格式输出：

# "a->b" ，如果 a != b
# "a" ，如果 a == b

from typing import List


class Solution:
    def summaryRanges(self, nums: List[int]) -> List[str]:
        nums.append(int(1e20))
        n, res = len(nums), []
        pre = nums[0]

        for i in range(1, n):
            if nums[i] == nums[i - 1] + 1:
                continue
            if nums[i - 1] == pre:
                res.append(str(pre))
            else:
                res.append(f'{pre}->{nums[i - 1]}')
            pre = nums[i]
        return res


print(Solution().summaryRanges([0, 1, 2, 4, 5, 7]))
