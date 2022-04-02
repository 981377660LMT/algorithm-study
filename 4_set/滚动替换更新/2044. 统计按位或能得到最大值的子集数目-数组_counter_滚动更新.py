from operator import or_
from functools import reduce, lru_cache
from typing import Counter, List


class Solution:
    def countMaxOrSubsets(self, nums: List[int]) -> int:
        # 数组滚动更新
        # res = [0]
        # for num in nums:
        #     res += [sub | num for sub in res]
        # counter = Counter(res)
        # return counter[max(counter)]

        # counter滚动更新
        counter = Counter({0: 1})
        for num in nums:
            for key, count in list(counter.items()):
                counter[key | num] += count
        return counter[max(counter)]


print(Solution().countMaxOrSubsets(nums=[3, 2, 1, 5]))
# 输出：6
# 解释：子集按位或可能的最大值是 7 。有 6 个子集按位或可以得到 7 ：
# - [3,5]
# - [3,1,5]
# - [3,2,5]
# - [3,2,1,5]
# - [2,5]
# - [2,1,5]
