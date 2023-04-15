from itertools import accumulate
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)
# 定义一个数组 arr 的 转换数组 conver 为：

# conver[i] = arr[i] + max(arr[0..i])，其中 max(arr[0..i]) 是满足 0 <= j <= i 的所有 arr[j] 中的最大值。
# 定义一个数组 arr 的 分数 为 arr 转换数组中所有元素的和。

# 给你一个下标从 0 开始长度为 n 的整数数组 nums ，请你返回一个长度为 n 的数组 ans ，其中 ans[i]是前缀 nums[0..i] 的分数。


class Solution:
    def findPrefixScore(self, nums: List[int]) -> List[int]:
        preMax = [-INF] + nums[:]
        for i in range(1, len(preMax)):
            if preMax[i] < preMax[i - 1]:
                preMax[i] = preMax[i - 1]
        preMin = [INF] + nums[:]
        for i in range(1, len(preMin)):
            if preMin[i] > preMin[i - 1]:
                preMin[i] = preMin[i - 1]
        sufMax = nums[:] + [-INF]
        for i in range(len(sufMax) - 2, -1, -1):
            if sufMax[i] < sufMax[i + 1]:
                sufMax[i] = sufMax[i + 1]
        sufMin = nums[:] + [INF]
        for i in range(len(sufMin) - 2, -1, -1):
            if sufMin[i] > sufMin[i + 1]:
                sufMin[i] = sufMin[i + 1]

        sufMax2: List[int] = ([-INF] + list(accumulate(nums[::-1], max)))[::-1]
        sufMin2: List[int] = ([INF] + list(accumulate(nums[::-1], min)))[::-1]
        print(sufMax, sufMax2)
        assert sufMax == sufMax2
        assert sufMin == sufMin2


# [2,3,7,5,10]
# [1,1,2,4,8,16]
# [17,17,12,14,13,8,8,9,10,18,10,18,10]
print(Solution().findPrefixScore([17, 17, 12, 14, 13, 8, 8, 9, 10, 18, 10, 18, 10]))
print(Solution().findPrefixScore([2, 3, 7, 5, 10]))
