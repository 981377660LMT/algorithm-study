from typing import List
from collections import defaultdict
from bisect import bisect_left

# 2 <= nums.length <= 105
# 1 <= queries.length <= 2 * 104
# 如果 a 中所有元素都 相同 ，那么差绝对值的最小值为 -1 。

# 差分数组
class Solution:
    def minDifference(self, nums: List[int], queries: List[List[int]]) -> List[int]:
        indexes = defaultdict(list)
        res = []

        for left, right in queries:
            ...

        return res


print(Solution().minDifference(nums=[1, 3, 4, 8], queries=[[0, 1], [1, 2], [2, 3], [0, 3]]))
# 输出：[2,1,4,1]
# 解释：查询结果如下：
# - queries[0] = [0,1]：子数组是 [1,3] ，差绝对值的最小值为 |1-3| = 2 。
# - queries[1] = [1,2]：子数组是 [3,4] ，差绝对值的最小值为 |3-4| = 1 。
# - queries[2] = [2,3]：子数组是 [4,8] ，差绝对值的最小值为 |4-8| = 4 。
# - queries[3] = [0,3]：子数组是 [1,3,4,8] ，差的绝对值的最小值为 |3-4| = 1 。

