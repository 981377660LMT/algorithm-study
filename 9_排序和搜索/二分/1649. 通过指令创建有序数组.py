from typing import List
from bisect import bisect_left, bisect_right, insort
from sortedcontainers import SortedList

# 返回将 instructions 中所有元素依次插入 nums 后的 总最小代价
# 每一次插入操作的 代价 是以下两者的 较小值 ：
# nums 中 严格小于  instructions[i] 的数字数目。
# nums 中 严格大于  instructions[i] 的数字数目。
class Solution:
    def createSortedArray(self, instructions: List[int]) -> int:
        res = 0
        sortedList = list()
        for num in instructions:
            # 严格小于
            smaller = bisect_left(sortedList, num)
            # 严格大于
            bigger = len(sortedList) - bisect_right(sortedList, num)
            res += min(smaller, bigger)
            sortedList[smaller:smaller] = [num]

        return res % (10 ** 9 + 7)


# 插入不要用bisect_insort而是直接切片
# 使用insort_left直接超时

print(Solution().createSortedArray([1, 5, 6, 2]))
# 输出：1
# 解释：一开始 nums = [] 。
# 插入 1 ，代价为 min(0, 0) = 0 ，现在 nums = [1] 。
# 插入 5 ，代价为 min(1, 0) = 0 ，现在 nums = [1,5] 。
# 插入 6 ，代价为 min(2, 0) = 0 ，现在 nums = [1,5,6] 。
# 插入 2 ，代价为 min(1, 2) = 1 ，现在 nums = [1,2,5,6] 。
# 总代价为 0 + 0 + 0 + 1 = 1 。
class Solution2:
    def createSortedArray(self, instructions: List[int]) -> int:
        res = 0
        list = SortedList()
        for num in instructions:
            # 严格小于
            smaller = list.bisect_left(num)
            # 严格大于
            bigger = len(list) - list.bisect_right(num)
            res += min(smaller, bigger)
            list.add(num)

        return res % (10 ** 9 + 7)
