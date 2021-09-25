# https://leetcode-cn.com/problems/split-array-into-consecutive-subsequences/solution/zui-hao-li-jie-de-pythonban-ben-by-user2198v/
from collections import defaultdict
from heapq import heappop, heappush
from typing import List


# 直观解法:O(n^2)
def isPossible1(self, nums: List[int]) -> bool:
    res = []
    for n in nums:
        for v in res:
            if n == v[-1] + 1:
                v.append(n)
                break
        else:
            res.insert(0, [n])

    return all([len(v) >= 3 for v in res])


# 时间复杂度nlog(n)
class Solution:
    def isPossible2(self, nums: List[int]) -> bool:
        # 只要知道子序列的最后一个数字和子序列的长度，就能确定子序列。
        # 字典的键为序列结尾数值，值为结尾为该数值的所有序列长度（以堆存储）
        # 每遍历一个数，将该数加入能加入的长度最短(此时需要最小堆)的序列中，不能加入序列则新建一个序列；然后更新字典。
        tails = defaultdict(list)
        for num in nums:
            if tails[num - 1]:
                pre_min_len = heappop(tails[num - 1])
                heappush(tails[num], pre_min_len + 1)
            else:
                heappush(tails[num], 1)
        res = all([val >= 3 for tail in tails.values() for val in tail])
        return res


# print(Solution().isPossible2([1, 2, 3, 3, 4, 4, 5, 5]))


