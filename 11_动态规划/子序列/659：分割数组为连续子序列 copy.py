# https://leetcode-cn.com/problems/split-array-into-consecutive-subsequences/solution/zui-hao-li-jie-de-pythonban-ben-by-user2198v/
from collections import defaultdict
from heapq import heappop, heappush
from typing import List


class Solution:
    def isPossible(self, nums: List[int]) -> bool:
        #  tails 键为数字 值为一个优先队列 记录以键值结尾的序列长度
        tails = defaultdict(list)
        for num in nums:
            if not tails[num - 1]:
                heappush(tails[num], 1)
            else:
                heappush(tails[num], heappop(tails[num - 1]) + 1)
        res = all([val >= 3 for tail in tails.values() for val in tail])
        return res


print(Solution().isPossible([1, 2, 3, 3, 4, 4, 5, 5]))

# 时间复杂度nlog(n)
