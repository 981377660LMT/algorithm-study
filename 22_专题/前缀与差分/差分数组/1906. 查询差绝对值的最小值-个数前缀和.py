# 1 <= nums[i] <= 100
# 2 <= nums.length <= 105
# 1 <= queries.length <= 2 * 104

# 看到 100 这个数据,就知道可以暴力枚举了，直接开数组
# 我们需要知道不同区间之间有哪些数，所以前缀统计每个区间不同的数的个数(1-100)
from typing import List


class Solution:
    def minDifference(self, nums: List[int], queries: List[List[int]]) -> List[int]:
        counters = [[0] * 101]
        for num in nums:
            counters.append(list(counters[-1]))
            counters[-1][num] += 1

        res = []
        for left, right in queries:
            minDiff = 0x3FFFFFFF
            pre = -0x3FFFFFFF
            for cand in range(1, 101):
                # 我们通过个数前缀和数组求得l到r之间有哪些数存在
                if counters[right + 1][cand] - counters[left][cand] > 0:
                    minDiff = min(minDiff, cand - pre)
                    pre = cand
            res.append(minDiff if minDiff != 0x3FFFFFFF else -1)
        return res


print(Solution().minDifference(nums=[1, 3, 4, 8], queries=[[0, 1], [1, 2], [2, 3], [0, 3]]))
# 输出：[2,1,4,1]
# 解释：查询结果如下：
# - queries[0] = [0,1]：子数组是 [1,3] ，差绝对值的最小值为 |1-3| = 2 。
# - queries[1] = [1,2]：子数组是 [3,4] ，差绝对值的最小值为 |3-4| = 1 。
# - queries[2] = [2,3]：子数组是 [4,8] ，差绝对值的最小值为 |4-8| = 4 。
# - queries[3] = [0,3]：子数组是 [1,3,4,8] ，差的绝对值的最小值为 |3-4| = 1 。
