from typing import List
from collections import defaultdict


class Solution:
    def distinctNumbers(self, nums: List[int], k: int) -> List[int]:
        # defaultdict 比Counter快很多
        counter = defaultdict(int)
        res = []
        for right, cur in enumerate(nums):
            counter[cur] += 1
            if right >= k:
                counter[nums[right - k]] -= 1
                if not counter[nums[right - k]]:
                    del counter[nums[right - k]]
            if right >= k - 1:
                res.append(len(counter))
        return res


print(Solution().distinctNumbers(nums=[1, 2, 3, 2, 2, 1, 3], k=3))

# 输出: [3,2,2,2,3]
# 解释：每个子数组的数字种类计算方法如下：
# - nums[0:2] = [1,2,3] 有'1','2','3'三种数字所以      ans[0] = 3
# - nums[1:3] = [2,3,2] 有'2','3'两种数字所以          ans[1] = 2
# - nums[2:4] = [3,2,2] 有'2','3'两种数字所以          ans[2] = 2
# - nums[3:5] = [2,2,1] 有'1','2'两种数字所以          ans[3] = 2
# - nums[4:6] = [2,1,3] 有'1','2','3'三种数字所以      ans[4] = 3
