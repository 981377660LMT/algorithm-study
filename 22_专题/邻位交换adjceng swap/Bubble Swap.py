# 求邻位交换的最小次数，使得nums1变为nums2


# 2193. 得到回文串的最少操作次数-贪心找最近
# 贪心：每个字符 找到他的最近的对应位置 交换
from typing import List

from sortedcontainers import SortedList

# n<=500


class Solution:
    def solve(self, lst0: List[int], lst1: List[int]) -> int:
        res = 0

        for num in lst0:
            index = lst1.index(num)
            res += index
            lst1.pop(index)  # 已经被换到最左边了，所以减1

        return res

        # res = 0
        # sl=SortedList(lst0)
        # for num in lst0:
        #     index = lst1.index(num)
        #     res += index
        #     lst1.pop(index)

        # return res


print(Solution().solve([0, 1, 2], [2, 0, 1]))

