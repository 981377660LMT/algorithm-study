# 1 <= nums[i] <= 100
# 2 <= nums.length <= 105
# 1 <= queries.length <= 2 * 104

# 看到 100 这个数据,就知道可以暴力枚举了，直接开数组
# 我们需要知道不同区间之间有哪些数，所以前缀统计每个区间不同的数的个数(1-100)
from typing import List

INF = int(1e20)


class Solution:
    def minDifference(self, nums: List[int], queries: List[List[int]]) -> List[int]:
        max_ = max(nums)
        # n = len(nums)
        # preSum = [[0] * (max_ + 1) for _ in range(n + 1)]
        # for i in range(1, n + 1):
        #     preSum[i][nums[i - 1]] += 1
        #     for j in range((max_ + 1)):
        #         preSum[i][j] += preSum[i - 1][j]
        preSum = [[0] * (max_ + 1)]
        for num in nums:
            cur = preSum[-1][:]
            cur[num] += 1
            preSum.append(cur)

        res = []
        for left, right in queries:
            minDiff = INF
            pre = -INF
            for cur in range(1, (max_ + 1)):
                # 我们通过`个数前缀和数组求得l到r之间有哪些数存在`
                if preSum[right + 1][cur] - preSum[left][cur] > 0:
                    minDiff = min(minDiff, cur - pre)
                    pre = cur
            res.append(minDiff if minDiff != INF else -1)
        return res


print(Solution().minDifference(nums=[1, 3, 4, 8], queries=[[0, 1], [1, 2], [2, 3], [0, 3]]))
# 输出：[2,1,4,1]
# 解释：查询结果如下：
# - queries[0] = [0,1]：子数组是 [1,3] ，差绝对值的最小值为 |1-3| = 2 。
# - queries[1] = [1,2]：子数组是 [3,4] ，差绝对值的最小值为 |3-4| = 1 。
# - queries[2] = [2,3]：子数组是 [4,8] ，差绝对值的最小值为 |4-8| = 4 。
# - queries[3] = [0,3]：子数组是 [1,3,4,8] ，差的绝对值的最小值为 |3-4| = 1 。
