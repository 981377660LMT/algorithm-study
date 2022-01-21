from typing import List
from collections import defaultdict

# 0 <= rods.length <= 20


# https://leetcode-cn.com/problems/tallest-billboard/solution/yi-quan-ji-ben-mei-shuo-ming-bai-de-zhe-pian-kan-l/
# 找到两个和相等且最大的子序列


# 等价转换：对任何一个数，`可以用三种方式对待它，乘以1，-1或0`，目标是求和为0时的最大正数和
# 用字典来存储每一步的结果，键和值分别是(k:v) 总和以及正数和，
# 初始化时dp={0:0},表示和为0时的最大正数和为0
# 那么最后只需要求dp[0]的最大值就ok
class Solution:
    def tallestBillboard(self, rods: List[int]) -> int:
        dp = defaultdict(int, {0: 0})
        for num in rods:
            for preSum, maxSum in list(dp.items()):
                dp[preSum + num] = max(dp[preSum + num], maxSum + num)
                dp[preSum - num] = max(dp[preSum - num], maxSum)
        return dp[0]


print(Solution().tallestBillboard([1, 2, 3, 4, 5, 6]))
# 输出：10
# 解释：我们有两个不相交的子集 {2,3,5} 和 {4,6}，它们具有相同的和 sum = 10。
