from typing import List


class Solution:
    def arrayRankTransform(self, arr: List[int]) -> List[int]:
        nums = sorted(set(arr))
        rank_by_value = {}
        for index, value in enumerate(nums, 1):
            rank_by_value[value] = index
        return [rank_by_value[value] for value in arr]


print(Solution().arrayRankTransform([40, 10, 20, 30]))
# 输出：[4,1,2,3]
# 解释：40 是最大的元素。 10 是最小的元素。 20 是第二小的数字。 30 是第三小的数字。
print(Solution().arrayRankTransform([100, 100, 100]))
# 输出：[1,1,1]
# 解释：所有元素有相同的序号。
