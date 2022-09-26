from itertools import pairwise
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 力扣城计划在两地设立「力扣嘉年华」的分会场，气象小组正在分析两地区的气温变化趋势，对于第 i ~ (i+1) 天的气温变化趋势，将根据以下规则判断：

# 若第 i+1 天的气温 高于 第 i 天，为 上升 趋势
# 若第 i+1 天的气温 等于 第 i 天，为 平稳 趋势
# 若第 i+1 天的气温 低于 第 i 天，为 下降 趋势
# 已知 temperatureA[i] 和 temperatureB[i] 分别表示第 i 天两地区的气温。
# 组委会希望找到一段天数尽可能多，且两地气温变化趋势相同的时间举办嘉年华活动。请分析并返回两地气温变化趋势相同的最大连续天数。

# 即最大的 n，使得第 i~i+n 天之间，两地气温变化趋势相同


class Solution:
    def temperatureTrend(self, temperatureA: List[int], temperatureB: List[int]) -> int:
        nums1, nums2 = [], []
        for a, b in pairwise(temperatureA):
            if a > b:
                nums1.append(1)
            elif a == b:
                nums1.append(0)
            else:
                nums1.append(-1)

        for a, b in pairwise(temperatureB):
            if a > b:
                nums2.append(1)
            elif a == b:
                nums2.append(0)
            else:
                nums2.append(-1)

        dp = 0
        res = 0
        for a, b in zip(nums1, nums2):
            if a == b:
                dp += 1
                res = max(res, dp)
            else:
                dp = 0
        return res


print(
    Solution().temperatureTrend(
        temperatureA=[21, 18, 18, 18, 31], temperatureB=[34, 32, 16, 16, 17]
    )
)
